package helper

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

// FindProjectRoot returns the closest parent directory containing a go.mod.
func FindProjectRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("resolve project path %q: %w", start, err)
	}

	for {
		modPath := filepath.Join(dir, "go.mod")
		info, statErr := os.Stat(modPath)
		if statErr == nil {
			if info.IsDir() {
				return "", fmt.Errorf("%s is a directory, not a file", modPath)
			}
			return dir, nil
		}
		if !os.IsNotExist(statErr) {
			return "", fmt.Errorf("stat %s: %w", modPath, statErr)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s to filesystem root", start)
		}
		dir = parent
	}
}

// GetProjectName reads the module path from the closest go.mod at or above dir.
func GetProjectName(dir string) (string, error) {
	root, err := FindProjectRoot(dir)
	if err != nil {
		return "", err
	}
	return ReadModulePath(root)
}

// ReadModulePath reads the module path from go.mod in dir without searching
// parent directories. It is useful when validating a newly cloned project.
func ReadModulePath(dir string) (string, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("resolve module directory %q: %w", dir, err)
	}
	modPath := filepath.Join(root, "go.mod")
	data, err := os.ReadFile(modPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", modPath, err)
	}
	parsed, err := modfile.Parse(modPath, data, nil)
	if err != nil {
		return "", fmt.Errorf("parse %s: %w", modPath, err)
	}
	if parsed.Module == nil || parsed.Module.Mod.Path == "" {
		return "", fmt.Errorf("module directive not found in %s", modPath)
	}
	return parsed.Module.Mod.Path, nil
}

func SplitArgs(cmd *cobra.Command, args []string) (cmdArgs, programArgs []string) {
	dashAt := cmd.ArgsLenAtDash()
	if dashAt >= 0 {
		return args[:dashAt], args[dashAt:]
	}
	return args, []string{}
}
func FindMain(base, excludeDir string) (map[string]string, error) {
	base, err := filepath.Abs(base)
	if err != nil {
		return nil, fmt.Errorf("resolve search path: %w", err)
	}
	excludeDirArr := normalizeExcludedDirs(excludeDir)
	cmdPath := make(map[string]string)
	seenDirs := make(map[string]struct{})
	err = filepath.WalkDir(base, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}
		if isExcludedRelativePath(relPath, excludeDirArr) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		isMain, err := isMainFile(path)
		if err != nil {
			return err
		}
		if !isMain {
			return nil
		}

		dir := filepath.Dir(path)
		if _, exists := seenDirs[dir]; exists {
			return nil
		}
		seenDirs[dir] = struct{}{}
		cmdPath[filepath.ToSlash(relPath)] = dir
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cmdPath, nil
}

func normalizeExcludedDirs(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		result = append(result, filepath.Clean(filepath.FromSlash(part)))
	}
	return result
}

func isExcludedRelativePath(path string, excluded []string) bool {
	path = filepath.Clean(path)
	for _, item := range excluded {
		if path == item || strings.HasPrefix(path, item+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

func isMainFile(path string) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("read %s: %w", path, err)
	}
	file, _ := parser.ParseFile(token.NewFileSet(), path, content, parser.SkipObjectResolution)
	if file == nil {
		return false, nil
	}
	if file.Name.Name != "main" {
		return false, nil
	}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && function.Recv == nil && function.Name.Name == "main" {
			return true, nil
		}
	}
	return false, nil
}
