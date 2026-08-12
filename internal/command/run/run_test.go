package run

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/fsnotify/fsnotify"
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
	if err := addWatchDirs(watcher, root, root, []string{"vendor"}); err != nil {
		t.Fatal(err)
	}

	newDir := filepath.Join(root, "internal", "service")
	if err := os.MkdirAll(newDir, 0o755); err != nil {
		t.Fatal(err)
	}
	relevant, err := handleWatchEvent(watcher, root, fsnotify.Event{Name: filepath.Join(root, "internal"), Op: fsnotify.Create}, []string{"vendor"}, includeExtSet("go"))
	if err != nil {
		t.Fatalf("handleWatchEvent() error = %v", err)
	}
	if !relevant {
		t.Fatal("handleWatchEvent() did not mark a new directory as relevant")
	}

	goFile := filepath.Join(newDir, "user.go")
	relevant, err = handleWatchEvent(watcher, root, fsnotify.Event{Name: goFile, Op: fsnotify.Write}, []string{"vendor"}, includeExtSet("go"))
	if err != nil || !relevant {
		t.Fatalf("handleWatchEvent(go file) = %v, %v; want true, nil", relevant, err)
	}
	markdown := filepath.Join(newDir, "README.md")
	relevant, err = handleWatchEvent(watcher, root, fsnotify.Event{Name: markdown, Op: fsnotify.Write}, []string{"vendor"}, includeExtSet("go"))
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
		fsnotify.Event{Name: filepath.Join(root, "vendor", "module.go"), Op: fsnotify.Write},
		[]string{"vendor"},
		includeExtSet("go"),
	)
	if err != nil || relevant {
		t.Fatalf("handleWatchEvent(excluded) = %v, %v; want false, nil", relevant, err)
	}
}
