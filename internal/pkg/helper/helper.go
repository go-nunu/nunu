package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func GetProjectName(dir string) string {
	modFile, err := os.Open(dir + "/go.mod")
	if err != nil {
		fmt.Println("go.mod does not exist", err)
		return ""
	}
	defer modFile.Close()

	var moduleName string
	_, err = fmt.Fscanf(modFile, "module %s", &moduleName)
	if err != nil {
		fmt.Println("read go mod error: ", err)
		return ""
	}
	return moduleName
}
func SplitArgs(cmd *cobra.Command, args []string) (cmdArgs, programArgs []string) {
	dashAt := cmd.ArgsLenAtDash()
	if dashAt >= 0 {
		return args[:dashAt], args[dashAt:]
	}
	return args, []string{}
}
func FindMain(base, excludeDir string) (map[string]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(wd, "/") {
		wd += "/"
	}
	excludeDirArr := strings.Split(excludeDir, ",")
	cmdPath := make(map[string]string)
	err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, s := range excludeDirArr {
			if strings.HasPrefix(strings.TrimPrefix(path, base), "/"+s) {
				return nil
			}
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if !strings.Contains(string(content), "package main") {
				return nil
			}
			re := regexp.MustCompile(`func\s+main\s*\(`)
			if re.Match(content) {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				d, _ := filepath.Split(absPath)
				cmdPath[strings.TrimPrefix(absPath, wd)] = d

			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cmdPath, nil
}
