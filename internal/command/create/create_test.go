package create

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCreateRejectsInvalidComponentName(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/project\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)

	err := runCreate(CmdCreateHandler, []string{"-"})
	if err == nil || !strings.Contains(err.Error(), "valid identifier") {
		t.Fatalf("runCreate() error = %v, want invalid identifier error", err)
	}
}

func TestRunCreateFromNestedDirectoryWritesToProjectRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/project\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "cmd", "server")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(nested)

	if err := runCreate(CmdCreateModel, []string{"user"}); err != nil {
		t.Fatalf("runCreate() error = %v", err)
	}
	want := filepath.Join(root, "internal", "model", "user.go")
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("generated model %s: %v", want, err)
	}
	unexpected := filepath.Join(nested, "internal", "model", "user.go")
	if _, err := os.Stat(unexpected); !os.IsNotExist(err) {
		t.Fatalf("generated model relative to nested working directory: %s", unexpected)
	}
}

func TestCreateCommandRejectsUnknownPositionalArguments(t *testing.T) {
	if err := CmdCreate.Args(CmdCreate, []string{"unknown", "name"}); err == nil {
		t.Fatal("create command accepted positional arguments instead of a subcommand")
	}
}
