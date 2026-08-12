package helper

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetProjectNameFromNestedDirectory(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "go.mod"), "// generated file\n\nmodule example.com/acme/service\n\ngo 1.25\n")
	nested := filepath.Join(root, "internal", "service")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	gotRoot, err := FindProjectRoot(nested)
	if err != nil {
		t.Fatalf("FindProjectRoot() error = %v", err)
	}
	if gotRoot != root {
		t.Fatalf("FindProjectRoot() = %q, want %q", gotRoot, root)
	}

	module, err := GetProjectName(nested)
	if err != nil {
		t.Fatalf("GetProjectName() error = %v", err)
	}
	if module != "example.com/acme/service" {
		t.Fatalf("GetProjectName() = %q", module)
	}
}

func TestGetProjectNameRejectsInvalidGoMod(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "go.mod"), "go 1.25\n")

	_, err := GetProjectName(root)
	if err == nil || !strings.Contains(err.Error(), "module directive not found") {
		t.Fatalf("GetProjectName() error = %v, want missing module error", err)
	}
}

func TestReadModulePathDoesNotFallBackToParentModule(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "go.mod"), "module example.com/parent\n\ngo 1.25\n")
	child := filepath.Join(root, "new-project")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := ReadModulePath(child)
	if err == nil || !strings.Contains(err.Error(), "go.mod") {
		t.Fatalf("ReadModulePath() error = %v, want missing child go.mod error", err)
	}
}

func TestFindProjectRootReturnsErrorOutsideModule(t *testing.T) {
	_, err := FindProjectRoot(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "go.mod not found") {
		t.Fatalf("FindProjectRoot() error = %v, want not found error", err)
	}
}

func TestFindMainUsesGoSyntaxAndExcludesDirectories(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "cmd", "api", "main.go"), "package main\nfunc main() {}\n")
	writeTestFile(t, filepath.Join(root, "cmd", "api", "other.go"), "package main\nfunc helper() {}\n")
	writeTestFile(t, filepath.Join(root, "cmd", "fake", "fake.go"), "package fake\n// package main; func main() {}\n")
	writeTestFile(t, filepath.Join(root, "cmd", "broken", "broken.go"), "package main\nfunc broken(\n")
	writeTestFile(t, filepath.Join(root, "vendor", "tool", "main.go"), "package main\nfunc main() {}\n")
	writeTestFile(t, filepath.Join(root, "cmd", "testonly", "main_test.go"), "package main\nfunc main() {}\n")

	got, err := FindMain(root, "vendor")
	if err != nil {
		t.Fatalf("FindMain() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("FindMain() returned %v, want one entry", got)
	}
	path, ok := got["cmd/api/main.go"]
	if !ok {
		t.Fatalf("FindMain() returned %v, want normalized cmd/api/main.go key", got)
	}
	wantDir := filepath.Join(root, "cmd", "api")
	if path != wantDir {
		t.Fatalf("FindMain() directory = %q, want %q", path, wantDir)
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
