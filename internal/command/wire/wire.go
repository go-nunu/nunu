package wire

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-nunu/nunu/internal/pkg/helper"
	"github.com/spf13/cobra"
)

var CmdWire = &cobra.Command{
	Use:     "wire",
	Short:   "nunu wire [wire.go path]",
	Long:    "nunu wire [wire.go path]",
	Example: "nunu wire cmd/server",
	Run: func(cmd *cobra.Command, args []string) {
		cmdArgs, _ := helper.SplitArgs(cmd, args)
		var dir string
		if len(cmdArgs) > 0 {
			dir = cmdArgs[0]
		}
		base, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
			return
		}
		if dir == "" {
			// find the directory containing the cmd/*
			wirePath, err := findWire(base)

			if err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
				return
			}
			switch len(wirePath) {
			case 0:
				fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", "The wire.go cannot be found in the current directory")
				return
			case 1:
				for _, v := range wirePath {
					dir = v
				}
			default:
				var wirePaths []string
				for k := range wirePath {
					wirePaths = append(wirePaths, k)
				}
				sort.Strings(wirePaths)
				prompt := &survey.Select{
					Message:  "Which directory do you want to run?",
					Options:  wirePaths,
					PageSize: 10,
				}
				e := survey.AskOne(prompt, &dir)
				if e != nil || dir == "" {
					return
				}
				dir = wirePath[dir]
			}
		}
		wire(dir)

	},
}
var CmdWireAll = &cobra.Command{
	Use:     "all",
	Short:   "nunu wire all",
	Long:    "nunu wire all",
	Example: "nunu wire all",
	Run: func(cmd *cobra.Command, args []string) {
		cmdArgs, _ := helper.SplitArgs(cmd, args)
		var dir string
		if len(cmdArgs) > 0 {
			dir = cmdArgs[0]
		}
		base, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
			return
		}
		if dir == "" {
			// find the directory containing the cmd/*
			wirePath, err := findWire(base)

			if err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
				return
			}
			switch len(wirePath) {
			case 0:
				fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", "The wire.go cannot be found in the current directory")
				return
			default:
				for _, v := range wirePath {
					wire(v)
				}
			}
		}

	},
}

func wire(wirePath string) {
	fmt.Println("wire.go path: ", wirePath)
	cmd := exec.Command("wire")
	cmd.Dir = wirePath
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("wire fail:", err)
	}
	fmt.Println(string(out))
}
func findWire(base string) (map[string]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(wd, "/") {
		wd += "/"
	}

	var root bool
	next := func(dir string) (map[string]string, error) {
		wirePath := make(map[string]string)
		err = filepath.Walk(dir, func(walkPath string, info os.FileInfo, err error) error {
			// multi level directory is not allowed under the wirePath directory, so it is judged that the path ends with wirePath.
			if strings.HasSuffix(walkPath, "wire.go") {
				p, _ := filepath.Split(walkPath)
				wirePath[strings.TrimPrefix(walkPath, wd)] = p
				return nil
			}
			if info.Name() == "go.mod" {
				root = true
			}
			return nil
		})
		return wirePath, err
	}
	for i := 0; i < 5; i++ {
		tmp := base
		cmd, err := next(tmp)
		if err != nil {
			return nil, err
		}
		if len(cmd) > 0 {
			return cmd, nil
		}
		if root {
			break
		}
		_ = filepath.Join(base, "..")
	}
	return map[string]string{"": base}, nil
}
