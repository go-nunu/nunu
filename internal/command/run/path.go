package run

import (
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
		includeExtMap[strings.TrimPrefix(strings.ToLower(ext), ".")] = struct{}{}
	}
	return includeExtMap
}
