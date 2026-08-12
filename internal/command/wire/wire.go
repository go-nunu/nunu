package wire

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-nunu/nunu/internal/pkg/helper"
	"github.com/spf13/cobra"
)

var CmdWire = &cobra.Command{
	Use:     "wire [directory]",
	Short:   "nunu wire [wire.go path]",
	Long:    "nunu wire [wire.go path]",
	Example: "nunu wire cmd/server",
	Args:    cobra.MaximumNArgs(1),
	RunE:    runWire,
}
var CmdWireAll = &cobra.Command{
	Use:     "all",
	Short:   "nunu wire all",
	Long:    "nunu wire all",
	Example: "nunu wire all",
	Args:    cobra.NoArgs,
	RunE:    runWireAll,
}

func runWire(_ *cobra.Command, args []string) error {
	var dir string
	if len(args) == 1 {
		dir = args[0]
	}
	if dir == "" {
		base, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current directory: %w", err)
		}
		wirePaths, err := findWire(base)
		if err != nil {
			return err
		}
		switch len(wirePaths) {
		case 0:
			return fmt.Errorf("wire.go not found in the current project")
		case 1:
			for _, path := range wirePaths {
				dir = path
			}
		default:
			labels := sortedKeys(wirePaths)
			var selected string
			prompt := &survey.Select{
				Message:  "Which directory do you want to run?",
				Options:  labels,
				PageSize: 10,
			}
			if err := survey.AskOne(prompt, &selected); err != nil {
				return fmt.Errorf("select wire directory: %w", err)
			}
			if selected == "" {
				return fmt.Errorf("no wire directory selected")
			}
			dir = wirePaths[selected]
		}
	}
	return wire(dir)
}

func runWireAll(_ *cobra.Command, _ []string) error {
	base, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current directory: %w", err)
	}
	wirePaths, err := findWire(base)
	if err != nil {
		return err
	}
	if len(wirePaths) == 0 {
		return fmt.Errorf("wire.go not found in the current project")
	}
	var wireErrors []error
	for _, label := range sortedKeys(wirePaths) {
		if err := wire(wirePaths[label]); err != nil {
			wireErrors = append(wireErrors, fmt.Errorf("generate %s: %w", label, err))
		}
	}
	return errors.Join(wireErrors...)
}

func wire(wirePath string) error {
	fmt.Println("wire.go path: ", wirePath)
	cmd := exec.Command("wire")
	cmd.Dir = wirePath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("run wire in %s: %w\n%s", wirePath, err, out)
	}
	if len(out) > 0 {
		fmt.Print(string(out))
	}
	return nil
}

func findWire(base string) (map[string]string, error) {
	base, err := filepath.Abs(base)
	if err != nil {
		return nil, fmt.Errorf("resolve search path: %w", err)
	}
	searchRoot := base
	if projectRoot, rootErr := helper.FindProjectRoot(base); rootErr == nil {
		searchRoot = projectRoot
	}

	wirePaths := make(map[string]string)
	err = filepath.WalkDir(searchRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "vendor":
				if path != searchRoot {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if entry.Name() != "wire.go" {
			return nil
		}
		rel, err := filepath.Rel(searchRoot, path)
		if err != nil {
			return err
		}
		wirePaths[filepath.ToSlash(rel)] = filepath.Dir(path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("search for wire.go: %w", err)
	}
	return wirePaths, nil
}

func sortedKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
