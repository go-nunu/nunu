//go:build !plan9 && !windows

package run

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWatchReloadsForNewFilesAndNewDirectoriesWithDebounce(t *testing.T) {
	root := t.TempDir()
	writeWatchTestFile(t, filepath.Join(root, "go.mod"), "module example.com/watchtest\n\ngo 1.25\n")
	writeWatchTestFile(t, filepath.Join(root, "main.go"), `package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	f, err := os.OpenFile("starts.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil { panic(err) }
	_, _ = fmt.Fprintln(f, "started")
	_ = f.Close()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
`)
	t.Chdir(root)

	oldExcludeDir, oldIncludeExt := excludeDir, includeExt
	excludeDir, includeExt = ".git,vendor", "go,yaml"
	t.Cleanup(func() {
		excludeDir, includeExt = oldExcludeDir, oldIncludeExt
	})

	quit := make(chan os.Signal, 1)
	done := make(chan error, 1)
	go func() {
		done <- watch(".", nil, nil, quit)
		close(done)
	}()
	t.Cleanup(func() {
		select {
		case quit <- os.Interrupt:
		default:
		}
		select {
		case <-done:
		case <-time.After(10 * time.Second):
		}
	})

	startsLog := filepath.Join(root, "starts.log")
	waitForStartCount(t, startsLog, 1)

	// Creating a new source file must be observed because directories, not only
	// files that existed at startup, are watched.
	writeWatchTestFile(t, filepath.Join(root, "helper.go"), "package main\n\nconst helperValue = 1\n")
	waitForStartCount(t, startsLog, 2)

	// A new directory is added to the watcher recursively. A later edit inside
	// that directory must trigger another reload.
	configFile := filepath.Join(root, "config", "app.yaml")
	writeWatchTestFile(t, configFile, "name: first\n")
	waitForStartCount(t, startsLog, 3)
	writeWatchTestFile(t, configFile, "name: second\n")
	waitForStartCount(t, startsLog, 4)

	// A burst of writes should collapse into one restart.
	for i := 0; i < 3; i++ {
		writeWatchTestFile(t, configFile, fmt.Sprintf("revision: %d\n", i))
		time.Sleep(25 * time.Millisecond)
	}
	waitForStartCount(t, startsLog, 5)
	time.Sleep(2 * reloadDebounce)
	if got := startCount(startsLog); got != 5 {
		t.Fatalf("debounced write burst produced %d starts, want 5", got)
	}

	quit <- os.Interrupt
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("watch() error = %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("watch() did not stop after interrupt")
	}
}

func waitForStartCount(t *testing.T, path string, want int) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if startCount(path) >= want {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("application start count = %d, want at least %d", startCount(path), want)
}

func startCount(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	return strings.Count(string(data), "started\n")
}

func writeWatchTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
