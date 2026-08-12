package wire

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindWireSearchesProjectRootAndUsesExactFilename(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "go.mod"), "module example.com/project\n\ngo 1.25\n")
	writeFile(t, filepath.Join(root, "cmd", "api", "wire.go"), "//go:build wireinject\n\npackage main\n")
	writeFile(t, filepath.Join(root, "cmd", "worker", "notwire.go"), "package main\n")
	writeFile(t, filepath.Join(root, "vendor", "example", "wire.go"), "package main\n")
	start := filepath.Join(root, "internal", "service")
	if err := os.MkdirAll(start, 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := findWire(start)
	if err != nil {
		t.Fatalf("findWire() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("findWire() = %v, want one result", got)
	}
	for label, dir := range got {
		if label != "cmd/api/wire.go" {
			t.Fatalf("findWire() label = %q", label)
		}
		if dir != filepath.Join(root, "cmd", "api") {
			t.Fatalf("findWire() dir = %q", dir)
		}
	}
}

func TestRunWireReturnsErrorWhenWireFileIsMissing(t *testing.T) {
	t.Chdir(t.TempDir())
	err := runWire(nil, nil)
	if err == nil || !strings.Contains(err.Error(), "wire.go not found") {
		t.Fatalf("runWire() error = %v, want wire.go not found", err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
