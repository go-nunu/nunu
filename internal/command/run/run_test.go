package run

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-nunu/nunu/internal/pkg/pathignore"
)

func TestSplitBuildFlagsPreservesQuotedValues(t *testing.T) {
	want := []string{"-tags", "integration", "-ldflags=-X main.version=hello world"}
	got, err := splitBuildFlags(`-tags integration '-ldflags=-X main.version=hello world'`)
	if err != nil {
		t.Fatalf("splitBuildFlags() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitBuildFlags() = %#v, want %#v", got, want)
	}
}

func TestSplitBuildFlagsRejectsUnclosedQuote(t *testing.T) {
	if _, err := splitBuildFlags(`-ldflags="unfinished`); err == nil {
		t.Fatal("splitBuildFlags() returned nil error for an unclosed quote")
	}
}

func TestIncludeExtSetIsCaseInsensitive(t *testing.T) {
	exts := includeExtSet(".GO, yaml")
	for _, ext := range []string{"go", "yaml"} {
		if _, ok := exts[ext]; !ok {
			t.Fatalf("includeExtSet() missing %q: %v", ext, exts)
		}
	}
}

func TestHandleWatchEventAddsNewDirectory(t *testing.T) {
	root := t.TempDir()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()
	excludes := mustExcludeMatcher(t, "vendor")
	if err := addWatchDirs(watcher, root, root, excludes); err != nil {
		t.Fatal(err)
	}

	newDir := filepath.Join(root, "internal", "service")
	if err := os.MkdirAll(newDir, 0o755); err != nil {
		t.Fatal(err)
	}
	relevant, err := handleWatchEvent(watcher, root, fsnotify.Event{Name: filepath.Join(root, "internal"), Op: fsnotify.Create}, excludes, includeExtSet("go"))
	if err != nil {
		t.Fatalf("handleWatchEvent() error = %v", err)
	}
	if !relevant {
		t.Fatal("handleWatchEvent() did not mark a new directory as relevant")
	}

	goFile := filepath.Join(newDir, "user.go")
	relevant, err = handleWatchEvent(watcher, root, fsnotify.Event{Name: goFile, Op: fsnotify.Write}, excludes, includeExtSet("go"))
	if err != nil || !relevant {
		t.Fatalf("handleWatchEvent(go file) = %v, %v; want true, nil", relevant, err)
	}
	markdown := filepath.Join(newDir, "README.md")
	relevant, err = handleWatchEvent(watcher, root, fsnotify.Event{Name: markdown, Op: fsnotify.Write}, excludes, includeExtSet("go"))
	if err != nil || relevant {
		t.Fatalf("handleWatchEvent(markdown) = %v, %v; want false, nil", relevant, err)
	}
}

func TestHandleWatchEventIgnoresExcludedDirectory(t *testing.T) {
	root := t.TempDir()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()

	relevant, err := handleWatchEvent(
		watcher,
		root,
		fsnotify.Event{Name: filepath.Join(root, "app", "vendor", "module.go"), Op: fsnotify.Write},
		mustExcludeMatcher(t, "vendor"),
		includeExtSet("go"),
	)
	if err != nil || relevant {
		t.Fatalf("handleWatchEvent(excluded) = %v, %v; want false, nil", relevant, err)
	}
}

func TestResolveRunTargetUsesTargetProjectRootAndAbsoluteTarget(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/project\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "app", "admin", "cmd", "server")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	other := t.TempDir()
	t.Chdir(other)

	gotTarget, gotRoot, err := resolveRunTarget(target)
	if err != nil {
		t.Fatalf("resolveRunTarget() error = %v", err)
	}
	wantTarget, err := filepath.EvalSymlinks(target)
	if err != nil {
		t.Fatal(err)
	}
	wantRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	if gotTarget != wantTarget {
		t.Fatalf("resolveRunTarget() target = %q, want %q", gotTarget, wantTarget)
	}
	if gotRoot != wantRoot {
		t.Fatalf("resolveRunTarget() root = %q, want %q", gotRoot, wantRoot)
	}
}

func TestStartProcessUsesProjectRootAsWorkingDirectory(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "cmd", "server")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/project\n\ngo 1.25.8\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	program := `package main
import "os"
func main() {
	wd, err := os.Getwd()
	if err != nil { panic(err) }
	if err := os.WriteFile("working-directory.txt", []byte(wd), 0o644); err != nil { panic(err) }
}
`
	if err := os.WriteFile(filepath.Join(target, "main.go"), []byte(program), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(t.TempDir())

	process, err := startProcess(root, target, nil, nil)
	if err != nil {
		t.Fatalf("startProcess() error = %v", err)
	}
	select {
	case err := <-process.done:
		if err != nil {
			t.Fatalf("process error = %v", err)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("process did not exit")
	}
	data, err := os.ReadFile(filepath.Join(root, "working-directory.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != root {
		t.Fatalf("process working directory = %q, want %q", data, root)
	}
}

func mustExcludeMatcher(t *testing.T, value string) *pathignore.Matcher {
	t.Helper()
	matcher, err := pathignore.CompileCSV(value)
	if err != nil {
		t.Fatal(err)
	}
	return matcher
}
