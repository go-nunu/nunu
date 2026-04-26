//go:build windows
// +build windows

package run

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-nunu/nunu/config"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fsnotify/fsnotify"
	"github.com/go-nunu/nunu/internal/pkg/helper"
	"github.com/spf13/cobra"
)

type Run struct {
}

var excludeDir string
var includeExt string
var buildFlags string

func init() {
	CmdRun.Flags().StringVarP(&excludeDir, "excludeDir", "", excludeDir, `eg: nunu run --excludeDir="tmp,vendor,.git,.idea"`)
	CmdRun.Flags().StringVarP(&includeExt, "includeExt", "", includeExt, `eg: nunu run --includeExt="go,tpl,tmpl,html,yaml,yml,toml,ini,json"`)
	CmdRun.Flags().StringVarP(&buildFlags, "buildFlags", "", buildFlags, `eg: nunu run --buildFlags="-tags cse"`)
	if excludeDir == "" {
		excludeDir = config.RunExcludeDir
	}
	if includeExt == "" {
		includeExt = config.RunIncludeExt
	}
}

var CmdRun = &cobra.Command{
	Use:     "run",
	Short:   "nunu run [main.go path]",
	Long:    "nunu run [main.go path]",
	Example: "nunu run cmd/server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdArgs, programArgs := helper.SplitArgs(cmd, args)
		var dir string
		if len(cmdArgs) > 0 {
			dir = cmdArgs[0]
		}
		base, err := os.Getwd()
		if err != nil {
			return err
		}
		if dir == "" {
			cmdPath, err := helper.FindMain(base, excludeDir)

			if err != nil {
				return err
			}
			switch len(cmdPath) {
			case 0:
				return errors.New("the cmd directory cannot be found in the current directory")
			case 1:
				for _, v := range cmdPath {
					dir = v
				}
			default:
				var cmdPaths []string
				for k := range cmdPath {
					cmdPaths = append(cmdPaths, k)
				}
				sort.Strings(cmdPaths)
				prompt := &survey.Select{
					Message:  "Which directory do you want to run?",
					Options:  cmdPaths,
					PageSize: 10,
				}
				e := survey.AskOne(prompt, &dir)
				if e != nil {
					return e
				}
				if dir == "" {
					return errors.New("no directory selected")
				}
				dir = cmdPath[dir]
			}
		}
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(quit)
		fmt.Printf("\033[35mNunu run %s.\033[0m\n", dir)
		fmt.Printf("\033[35mWatch excludeDir %s\033[0m\n", excludeDir)
		fmt.Printf("\033[35mWatch includeExt %s\033[0m\n", includeExt)
		fmt.Printf("\033[35mWatch buildFlags %s\033[0m\n", buildFlags)
		return watch(dir, programArgs, quit)

	},
}

func watch(dir string, programArgs []string, quit <-chan os.Signal) error {

	// Listening file path
	watchPath := "./"

	// Create a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	excludeDirArr := splitCSV(excludeDir)
	buildFlagsArr := make([]string, 0)
	if buildFlags = strings.TrimSpace(buildFlags); buildFlags != "" {
		buildFlagsArr = strings.Fields(buildFlags)
	}
	includeExtMap := includeExtSet(includeExt)
	// Add files to watcher
	err = filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if isExcludedPath(path, excludeDirArr) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.IsDir() {
			ext := filepath.Ext(info.Name())
			if _, ok := includeExtMap[strings.TrimPrefix(ext, ".")]; ok {
				if err := watcher.Add(path); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	cmd, err := start(dir, buildFlagsArr, programArgs)
	if err != nil {
		return err
	}

	// Loop listening file modification
	for {
		select {
		case <-quit:
			if err := killProcess(cmd); err != nil {
				fmt.Printf("\033[31mserver exiting...\033[0m\n")
				return err
			}
			fmt.Printf("\033[31mserver exiting...\033[0m\n")
			return nil

		case event := <-watcher.Events:
			// The file has been modified or created
			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Printf("\033[36mfile modified: %s\033[0m\n", event.Name)
				if err := killProcess(cmd); err != nil {
					fmt.Println("Error:", err)
				}
				cmd, err = start(dir, buildFlagsArr, programArgs)
				if err != nil {
					return err
				}
			}
		case err := <-watcher.Errors:
			return err
		}
	}
}

func killProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	// 获取进程ID
	pid := cmd.Process.Pid
	// 构造taskkill命令
	taskkill := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
	err := taskkill.Run()
	if err != nil {
		return err
	}
	return nil
}
func start(dir string, buildFlagsArgs []string, programArgs []string) (*exec.Cmd, error) {
	run := []string{"run"}
	run = append(run, buildFlagsArgs...)
	run = append(run, dir)
	cmd := exec.Command("go", append(run, programArgs...)...)
	// Set a new process group to kill all child processes when the program exits

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("cmd run failed: %w", err)
	}
	time.Sleep(time.Second)
	fmt.Printf("\033[32;1mrunning...\033[0m\n")
	return cmd, nil
}
