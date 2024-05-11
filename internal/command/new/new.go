package new

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-nunu/nunu/config"
	"github.com/go-nunu/nunu/internal/pkg/helper"
	"github.com/spf13/cobra"
)

type Project struct {
	ProjectName string `survey:"name"`
}

var CmdNew = &cobra.Command{
	Use:     "new",
	Example: "nunu new demo-api",
	Short:   "create a new project.",
	Long:    `create a new project with nunu layout.`,
	Run:     run,
}
var (
	repoURL string
)

func init() {
	CmdNew.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")

}
func NewProject() *Project {
	return &Project{}
}

func run(cmd *cobra.Command, args []string) {
	p := NewProject()
	if len(args) == 0 {
		err := survey.AskOne(&survey.Input{
			Message: "What is your project name?",
			Help:    "project name.",
			Suggest: nil,
		}, &p.ProjectName, survey.WithValidator(survey.Required))
		if err != nil {
			return
		}
	} else {
		p.ProjectName = args[0]
	}

	// clone repo
	yes, err := p.cloneTemplate()
	if err != nil || !yes {
		return
	}

	err = p.replacePackageName()
	if err != nil || !yes {
		return
	}

	err = p.replacePackageName()
	if err != nil || !yes {
		return
	}
	err = p.modTidy()
	if err != nil || !yes {
		return
	}
	p.rmGit()
	p.installWire()
	fmt.Printf("\n _   _                   \n| \\ | |_   _ _ __  _   _ \n|  \\| | | | | '_ \\| | | |\n| |\\  | |_| | | | | |_| |\n|_| \\_|\\__,_|_| |_|\\__,_| \n \n" + "\x1B[38;2;66;211;146mA\x1B[39m \x1B[38;2;67;209;149mC\x1B[39m\x1B[38;2;68;206;152mL\x1B[39m\x1B[38;2;69;204;155mI\x1B[39m \x1B[38;2;70;201;158mt\x1B[39m\x1B[38;2;71;199;162mo\x1B[39m\x1B[38;2;72;196;165mo\x1B[39m\x1B[38;2;73;194;168ml\x1B[39m \x1B[38;2;74;192;171mf\x1B[39m\x1B[38;2;75;189;174mo\x1B[39m\x1B[38;2;76;187;177mr\x1B[39m \x1B[38;2;77;184;180mb\x1B[39m\x1B[38;2;78;182;183mu\x1B[39m\x1B[38;2;79;179;186mi\x1B[39m\x1B[38;2;80;177;190ml\x1B[39m\x1B[38;2;81;175;193md\x1B[39m\x1B[38;2;82;172;196mi\x1B[39m\x1B[38;2;83;170;199mn\x1B[39m\x1B[38;2;83;167;202mg\x1B[39m \x1B[38;2;84;165;205mg\x1B[39m\x1B[38;2;85;162;208mo\x1B[39m \x1B[38;2;86;160;211ma\x1B[39m\x1B[38;2;87;158;215mp\x1B[39m\x1B[38;2;88;155;218ml\x1B[39m\x1B[38;2;89;153;221mi\x1B[39m\x1B[38;2;90;150;224mc\x1B[39m\x1B[38;2;91;148;227ma\x1B[39m\x1B[38;2;92;145;230mt\x1B[39m\x1B[38;2;93;143;233mi\x1B[39m\x1B[38;2;94;141;236mo\x1B[39m\x1B[38;2;95;138;239mn\x1B[39m\x1B[38;2;96;136;243m.\x1B[39m\n\n")
	fmt.Printf("🎉 Project \u001B[36m%s\u001B[0m created successfully!\n\n", p.ProjectName)
	fmt.Printf("Done. Now run:\n\n")
	fmt.Printf("› \033[36mcd %s \033[0m\n", p.ProjectName)
	fmt.Printf("› \033[36mnunu run \033[0m\n\n")
}

func (p *Project) cloneTemplate() (bool, error) {
	stat, _ := os.Stat(p.ProjectName)
	if stat != nil {
		var overwrite = false

		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Folder %s already exists, do you want to overwrite it?", p.ProjectName),
			Help:    "Remove old project and create new project.",
		}
		err := survey.AskOne(prompt, &overwrite)
		if err != nil {
			return false, err
		}
		if !overwrite {
			return false, nil
		}
		err = os.RemoveAll(p.ProjectName)
		if err != nil {
			fmt.Println("remove old project error: ", err)
			return false, err
		}
	}
	repo := config.RepoBase

	if repoURL == "" {
		layout := ""
		prompt := &survey.Select{
			Message: "Please select a layout:",
			Options: []string{
				"Advanced",
				"Basic",
				"Chat",
			},
			Description: func(value string, index int) string {
				if index == 1 {
					return "A basic project structure"
				}
				if index == 2 {
					return "A simple chat room containing websocket/tcp"
				}
				return "It has rich functions such as db, jwt, cron, migration, test, etc"
			},
		}
		err := survey.AskOne(prompt, &layout)
		if err != nil {
			return false, err
		}
		if layout == "Advanced" {
			repo = config.RepoAdvanced
		} else if layout == "Chat" {
			repo = config.RepoChat
		}
		err = os.RemoveAll(p.ProjectName)
		if err != nil {
			fmt.Println("remove old project error: ", err)
			return false, err
		}
	} else {
		repo = repoURL
	}

	fmt.Printf("git clone %s\n", repo)
	cmd := exec.Command("git", "clone", repo, p.ProjectName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("git clone %s error: %s\n", repo, err)
		return false, err
	}
	return true, nil
}

func (p *Project) replacePackageName() error {
	packageName := helper.GetProjectName(p.ProjectName)

	err := p.replaceFiles(packageName)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "edit", "-module", p.ProjectName)
	cmd.Dir = p.ProjectName
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("go mod edit error: ", err)
		return err
	}
	return nil
}
func (p *Project) modTidy() error {
	fmt.Println("go mod tidy")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = p.ProjectName
	if err := cmd.Run(); err != nil {
		fmt.Println("go mod tidy error: ", err)
		return err
	}
	return nil
}
func (p *Project) rmGit() {
	os.RemoveAll(p.ProjectName + "/.git")
}
func (p *Project) installWire() {
	fmt.Printf("go install %s\n", config.WireCmd)
	cmd := exec.Command("go", "install", config.WireCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("go install %s error\n", err)
	}
}

func (p *Project) replaceFiles(packageName string) error {
	err := filepath.Walk(p.ProjectName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newData := bytes.ReplaceAll(data, []byte(packageName), []byte(p.ProjectName))
		if err := os.WriteFile(path, newData, 0644); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("walk file error: ", err)
		return err
	}
	return nil
}
