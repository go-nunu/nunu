package pathignore

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// Matcher applies comma-separated ignore patterns using gitignore-style path
// semantics. Patterns without a slash match a name at any directory depth;
// patterns containing a slash are relative to the project root.
type Matcher struct {
	rules []rule
}

type rule struct {
	pattern       string
	negated       bool
	directoryOnly bool
	basenameOnly  bool
	contentsOnly  bool
	contentsBase  string
}

// CompileCSV compiles the value accepted by nunu's --excludeDir flag.
func CompileCSV(value string) (*Matcher, error) {
	parts := strings.Split(value, ",")
	rules := make([]rule, 0, len(parts))
	for index, part := range parts {
		pattern := strings.TrimSpace(part)
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}

		negated := strings.HasPrefix(pattern, "!")
		if negated {
			pattern = strings.TrimSpace(strings.TrimPrefix(pattern, "!"))
		}
		if pattern == "" {
			return nil, fmt.Errorf("exclude pattern %d has no value after the negation marker", index+1)
		}

		// Flags are commonly copied between Unix and Windows shells. Internally
		// all matching uses slash-separated paths so the same value is portable.
		pattern = strings.ReplaceAll(pattern, `\`, "/")
		directoryOnly := strings.HasSuffix(pattern, "/")
		pattern = strings.TrimSuffix(pattern, "/")
		pattern = strings.TrimPrefix(pattern, "/")
		pattern = strings.TrimPrefix(pattern, "./")
		pattern = path.Clean(pattern)
		if pattern == "." || pattern == "" {
			return nil, fmt.Errorf("exclude pattern %d does not identify a path", index+1)
		}
		if !doublestar.ValidatePattern(pattern) {
			return nil, fmt.Errorf("invalid exclude pattern %q", part)
		}

		contentsBase := strings.TrimSuffix(pattern, "/**")
		rules = append(rules, rule{
			pattern:       normalizeCase(pattern),
			negated:       negated,
			directoryOnly: directoryOnly,
			basenameOnly:  !strings.Contains(pattern, "/"),
			contentsOnly:  contentsBase != pattern,
			contentsBase:  normalizeCase(contentsBase),
		})
	}
	return &Matcher{rules: rules}, nil
}

// Match reports whether a project-relative path should be ignored. isDir is
// required for patterns ending in a slash.
func (m *Matcher) Match(name string, isDir bool) bool {
	if m == nil {
		return false
	}
	name = normalizePath(name)
	if name == "" {
		return false
	}

	ignored := false
	for _, rule := range m.rules {
		if rule.matches(name, isDir) {
			ignored = !rule.negated
		}
	}
	return ignored
}

func (r rule) matches(name string, isDir bool) bool {
	parts := strings.Split(name, "/")
	if r.basenameOnly {
		for index, part := range parts {
			candidateIsDir := index < len(parts)-1 || isDir
			if r.directoryOnly && !candidateIsDir {
				continue
			}
			if doublestar.MatchUnvalidated(r.pattern, part) {
				return true
			}
		}
		return false
	}

	for index := range parts {
		candidate := strings.Join(parts[:index+1], "/")
		candidateIsDir := index < len(parts)-1 || isDir
		if r.directoryOnly && !candidateIsDir {
			continue
		}
		// In gitignore, dir/** ignores the contents of dir, not dir itself.
		// This distinction permits !dir/keep.go to re-include a direct child.
		if r.contentsOnly && candidate == r.contentsBase {
			continue
		}
		if doublestar.MatchUnvalidated(r.pattern, candidate) {
			return true
		}
	}
	return false
}

func normalizePath(name string) string {
	name = strings.ReplaceAll(name, `\`, "/")
	name = strings.TrimPrefix(name, "./")
	name = strings.TrimPrefix(name, "/")
	name = path.Clean(name)
	if name == "." {
		return ""
	}
	return normalizeCase(name)
}

func normalizeCase(value string) string {
	if runtime.GOOS == "windows" {
		return strings.ToLower(value)
	}
	return value
}
