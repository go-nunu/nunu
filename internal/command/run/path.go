package run

import (
	"os"
	"path/filepath"
	"strings"
)

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			items = append(items, part)
		}
	}
	return items
}

func includeExtSet(value string) map[string]struct{} {
	exts := splitCSV(value)
	includeExtMap := make(map[string]struct{}, len(exts))
	for _, ext := range exts {
		includeExtMap[strings.TrimPrefix(ext, ".")] = struct{}{}
	}
	return includeExtMap
}

func isExcludedPath(path string, excludeDirs []string) bool {
	cleanPath := filepath.Clean(path)
	for _, excludeDir := range excludeDirs {
		excludeDir = filepath.Clean(excludeDir)
		if cleanPath == excludeDir || strings.HasPrefix(cleanPath, excludeDir+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}
