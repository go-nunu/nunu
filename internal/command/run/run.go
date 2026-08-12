package run

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fsnotify/fsnotify"
	"github.com/go-nunu/nunu/config"
	"github.com/go-nunu/nunu/internal/pkg/helper"
	"github.com/kballard/go-shellquote"
	"github.com/spf13/cobra"
)

const (
	reloadDebounce  = 300 * time.Millisecond
	shutdownTimeout = 5 * time.Second
)

var excludeDir string
var includeExt string
var buildFlags string

func init() {
	CmdRun.Flags().StringVar(&excludeDir, "excludeDir", config.RunExcludeDir, `eg: nunu run --excludeDir="tmp,vendor,.git,.idea"`)
	CmdRun.Flags().StringVar(&includeExt, "includeExt", config.RunIncludeExt, `eg: nunu run --includeExt="go,tpl,tmpl,html,yaml,yml,toml,ini,json"`)
	CmdRun.Flags().StringVar(&buildFlags, "buildFlags", "", `eg: nunu run --buildFlags="-tags cse"`)
}

var CmdRun = &cobra.Command{
	Use:     "run [directory] [-- program arguments...]",
	Short:   "Run a Go application and reload it when project files change",
	Long:    "Run a Go application and reload it when project files change",
	Example: "nunu run cmd/server -- --config config/local.yml",
	Args: func(cmd *cobra.Command, args []string) error {
		cmdArgs, _ := helper.SplitArgs(cmd, args)
		if len(cmdArgs) > 1 {
			return fmt.Errorf("accepts at most one application directory, received %d", len(cmdArgs))
		}
		return nil
	},
	RunE: run,
}

func run(cmd *cobra.Command, args []string) error {
	cmdArgs, programArgs := helper.SplitArgs(cmd, args)
	var dir string
	if len(cmdArgs) == 1 {
		dir = cmdArgs[0]
	}
	if dir == "" {
		base, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current directory: %w", err)
		}
		if projectRoot, rootErr := helper.FindProjectRoot(base); rootErr == nil {
			base = projectRoot
		}
		cmdPaths, err := helper.FindMain(base, excludeDir)
		if err != nil {
			return fmt.Errorf("find main package: %w", err)
		}
		switch len(cmdPaths) {
		case 0:
			return errors.New("main package not found in the current project")
		case 1:
			for _, path := range cmdPaths {
				dir = path
			}
		default:
			labels := make([]string, 0, len(cmdPaths))
			for label := range cmdPaths {
				labels = append(labels, label)
			}
			sort.Strings(labels)
			var selected string
			prompt := &survey.Select{
				Message:  "Which directory do you want to run?",
				Options:  labels,
				PageSize: 10,
			}
			if err := survey.AskOne(prompt, &selected); err != nil {
				return fmt.Errorf("select main package: %w", err)
			}
			if selected == "" {
				return errors.New("no main package selected")
			}
			dir = cmdPaths[selected]
		}
	}

	buildFlagsArgs, err := splitBuildFlags(buildFlags)
	if err != nil {
		return err
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	fmt.Printf("\033[35mNunu run %s.\033[0m\n", dir)
	fmt.Printf("\033[35mWatch excludeDir %s\033[0m\n", excludeDir)
	fmt.Printf("\033[35mWatch includeExt %s\033[0m\n", includeExt)
	fmt.Printf("\033[35mWatch buildFlags %s\033[0m\n", buildFlags)
	return watch(dir, buildFlagsArgs, programArgs, quit)
}

type managedProcess struct {
	cmd  *exec.Cmd
	done chan error
}

func watch(dir string, buildFlagsArgs, programArgs []string, quit <-chan os.Signal) error {
	watchRoot, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("resolve watch root: %w", err)
	}
	if projectRoot, rootErr := helper.FindProjectRoot(watchRoot); rootErr == nil {
		watchRoot = projectRoot
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create file watcher: %w", err)
	}
	defer watcher.Close()

	excludeDirs := splitCSV(excludeDir)
	includeExts := includeExtSet(includeExt)
	if err := addWatchDirs(watcher, watchRoot, watchRoot, excludeDirs); err != nil {
		return err
	}

	process, err := startProcess(dir, buildFlagsArgs, programArgs)
	if err != nil {
		return err
	}
	processDone := (<-chan error)(process.done)

	var debounceTimer *time.Timer
	var debounceC <-chan time.Time
	defer func() {
		if debounceTimer != nil {
			debounceTimer.Stop()
		}
	}()

	scheduleReload := func() {
		if debounceTimer == nil {
			debounceTimer = time.NewTimer(reloadDebounce)
		} else {
			if !debounceTimer.Stop() {
				select {
				case <-debounceTimer.C:
				default:
				}
			}
			debounceTimer.Reset(reloadDebounce)
		}
		debounceC = debounceTimer.C
	}

	for {
		select {
		case <-quit:
			if err := stopProcess(process, shutdownTimeout); err != nil {
				return fmt.Errorf("stop application: %w", err)
			}
			fmt.Printf("\033[31mserver exiting...\033[0m\n")
			return nil

		case processErr := <-processDone:
			process = nil
			processDone = nil
			if processErr != nil {
				fmt.Fprintf(os.Stderr, "\033[31mapplication exited: %v\033[0m\n", processErr)
			} else {
				fmt.Println("application exited")
			}

		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("file watcher event channel closed")
			}
			relevant, err := handleWatchEvent(watcher, watchRoot, event, excludeDirs, includeExts)
			if err != nil {
				return err
			}
			if relevant {
				fmt.Printf("\033[36mfile modified: %s\033[0m\n", displayPath(watchRoot, event.Name))
				scheduleReload()
			}

		case watcherErr, ok := <-watcher.Errors:
			if !ok {
				return errors.New("file watcher error channel closed")
			}
			return fmt.Errorf("file watcher: %w", watcherErr)

		case <-debounceC:
			debounceC = nil
			if err := stopProcess(process, shutdownTimeout); err != nil {
				return fmt.Errorf("restart application: %w", err)
			}
			process, err = startProcess(dir, buildFlagsArgs, programArgs)
			if err != nil {
				return err
			}
			processDone = process.done
		}
	}
}

func splitBuildFlags(value string) ([]string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	args, err := shellquote.Split(value)
	if err != nil {
		return nil, fmt.Errorf("parse buildFlags: %w", err)
	}
	return args, nil
}

func addWatchDirs(watcher *fsnotify.Watcher, root, start string, excludeDirs []string) error {
	err := filepath.WalkDir(start, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() {
			return nil
		}
		if isExcludedFromRoot(root, path, excludeDirs) {
			return filepath.SkipDir
		}
		if err := watcher.Add(path); err != nil {
			return fmt.Errorf("watch directory %s: %w", path, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("discover watch directories: %w", err)
	}
	return nil
}

func handleWatchEvent(watcher *fsnotify.Watcher, root string, event fsnotify.Event, excludeDirs []string, includeExts map[string]struct{}) (bool, error) {
	if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove|fsnotify.Rename) == 0 {
		return false, nil
	}
	if isExcludedFromRoot(root, event.Name, excludeDirs) {
		return false, nil
	}

	if event.Op&fsnotify.Create != 0 {
		info, err := os.Stat(event.Name)
		if err == nil && info.IsDir() {
			if err := addWatchDirs(watcher, root, event.Name, excludeDirs); err != nil {
				return false, err
			}
			return true, nil
		}
		if err != nil && !os.IsNotExist(err) {
			return false, fmt.Errorf("stat created path %s: %w", event.Name, err)
		}
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(event.Name)), ".")
	if _, ok := includeExts[ext]; ok {
		return true, nil
	}
	// A removed or renamed extensionless path may be a source directory.
	return ext == "" && event.Op&(fsnotify.Remove|fsnotify.Rename) != 0, nil
}

func isExcludedFromRoot(root, path string, excludeDirs []string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return true
	}
	return isExcludedPath(rel, excludeDirs)
}

func displayPath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}

func startProcess(dir string, buildFlagsArgs, programArgs []string) (*managedProcess, error) {
	args := []string{"run"}
	args = append(args, buildFlagsArgs...)
	args = append(args, dir)
	args = append(args, programArgs...)
	cmd := exec.Command("go", args...)
	configureProcess(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start go %s: %w", strings.Join(args, " "), err)
	}

	process := &managedProcess{cmd: cmd, done: make(chan error, 1)}
	go func() {
		process.done <- cmd.Wait()
		close(process.done)
	}()
	fmt.Printf("\033[32;1mrunning (pid %d)...\033[0m\n", cmd.Process.Pid)
	return process, nil
}

func stopProcess(process *managedProcess, timeout time.Duration) error {
	if process == nil || process.cmd == nil || process.cmd.Process == nil {
		return nil
	}
	select {
	case <-process.done:
		return nil
	default:
	}

	interruptErr := interruptProcess(process.cmd)
	if interruptErr != nil {
		select {
		case <-process.done:
			return nil
		default:
		}
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-process.done:
		return nil
	case <-timer.C:
		if err := forceKillProcess(process.cmd); err != nil {
			select {
			case <-process.done:
				return nil
			default:
			}
			if interruptErr != nil {
				return fmt.Errorf("interrupt process: %v; force kill process: %w", interruptErr, err)
			}
			return fmt.Errorf("force kill process: %w", err)
		}
		select {
		case <-process.done:
			return nil
		case <-time.After(timeout):
			return errors.New("timed out waiting for killed process to exit")
		}
	}
}
