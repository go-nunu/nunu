package new

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
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
	RunE:    run,
}
var (
	repoURL string
)

type layoutOption struct {
	Name        string
	Repo        string
	Description string
}

var layoutOptions = []layoutOption{
	{
		Name:        "Advanced",
		Repo:        config.RepoAdvanced,
		Description: "Full-featured API service with examples for database, Redis, JWT, cron, migration, and tests.",
	},
	{
		Name:        "Basic",
		Repo:        config.RepoBase,
		Description: "Minimal layered API service for projects that only need the core Nunu structure.",
	},
	{
		Name:        "Admin",
		Repo:        config.RepoAdmin,
		Description: "Admin system template with Gin APIs, Vue 3 UI, JWT authentication, and Casbin RBAC.",
	},
	{
		Name:        "MCP Server",
		Repo:        config.RepoMCP,
		Description: "Model Context Protocol server template with STDIO, SSE, and Streamable HTTP transports.",
	},
	{
		Name:        "Monorepo",
		Repo:        config.RepoMonorepo,
		Description: "Multi-application workspace with shared packages, admin backend, embedded UI, and home app.",
	},
	{
		Name:        "Chat",
		Repo:        config.RepoChat,
		Description: "Real-time service template with WebSocket and TCP chat server examples.",
	},
}

func layoutNames() []string {
	names := make([]string, 0, len(layoutOptions))
	for _, option := range layoutOptions {
		names = append(names, option.Name)
	}
	return names
}

func findLayoutOption(name string) (layoutOption, bool) {
	for _, option := range layoutOptions {
		if option.Name == name {
			return option, true
		}
	}
	return layoutOption{}, false
}

func init() {
	CmdNew.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")

}
func NewProject() *Project {
	return &Project{}
}

func run(cmd *cobra.Command, args []string) error {
	p := NewProject()
	if len(args) == 0 {
		err := survey.AskOne(&survey.Input{
			Message: "What is your project name?",
			Help:    "project name.",
			Suggest: nil,
		}, &p.ProjectName, survey.WithValidator(survey.Required))
		if err != nil {
			return err
		}
	} else {
		p.ProjectName = args[0]
	}

	// clone repo
	cloned, err := p.cloneTemplate()
	if err != nil || !cloned {
		return err
	}

	err = p.replacePackageName()
	if err != nil {
		return err
	}
	err = p.modTidy()
	if err != nil {
		return err
	}
	p.rmGit()
	if err := p.installWire(); err != nil {
		return err
	}
	fmt.Printf("\n _   _                   \n| \\ | |_   _ _ __  _   _ \n|  \\| | | | | '_ \\| | | |\n| |\\  | |_| | | | | |_| |\n|_| \\_|\\__,_|_| |_|\\__,_| \n \n" + "\x1B[38;2;66;211;146mA\x1B[39m \x1B[38;2;67;209;149mC\x1B[39m\x1B[38;2;68;206;152mL\x1B[39m\x1B[38;2;69;204;155mI\x1B[39m \x1B[38;2;70;201;158mt\x1B[39m\x1B[38;2;71;199;162mo\x1B[39m\x1B[38;2;72;196;165mo\x1B[39m\x1B[38;2;73;194;168ml\x1B[39m \x1B[38;2;74;192;171mf\x1B[39m\x1B[38;2;75;189;174mo\x1B[39m\x1B[38;2;76;187;177mr\x1B[39m \x1B[38;2;77;184;180mb\x1B[39m\x1B[38;2;78;182;183mu\x1B[39m\x1B[38;2;79;179;186mi\x1B[39m\x1B[38;2;80;177;190ml\x1B[39m\x1B[38;2;81;175;193md\x1B[39m\x1B[38;2;82;172;196mi\x1B[39m\x1B[38;2;83;170;199mn\x1B[39m\x1B[38;2;83;167;202mg\x1B[39m \x1B[38;2;84;165;205mg\x1B[39m\x1B[38;2;85;162;208mo\x1B[39m \x1B[38;2;86;160;211ma\x1B[39m\x1B[38;2;87;158;215mp\x1B[39m\x1B[38;2;88;155;218ml\x1B[39m\x1B[38;2;89;153;221mi\x1B[39m\x1B[38;2;90;150;224mc\x1B[39m\x1B[38;2;91;148;227ma\x1B[39m\x1B[38;2;92;145;230mt\x1B[39m\x1B[38;2;93;143;233mi\x1B[39m\x1B[38;2;94;141;236mo\x1B[39m\x1B[38;2;95;138;239mn\x1B[39m\x1B[38;2;96;136;243m.\x1B[39m\n\n")
	fmt.Printf("🎉 Project \u001B[36m%s\u001B[0m created successfully!\n\n", p.ProjectName)
	fmt.Printf("Done. Now run:\n\n")
	fmt.Printf("› \033[36mcd %s \033[0m\n", p.ProjectName)
	fmt.Printf("› \033[36mnunu run \033[0m\n\n")
	return nil
}

func (p *Project) cloneTemplate() (bool, error) {
	_, err := os.Stat(p.ProjectName)
	if err == nil {
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
			return false, fmt.Errorf("remove old project: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("stat project directory: %w", err)
	}
	repo := config.RepoBase

	if repoURL == "" {
		layout := ""
		prompt := &survey.Select{
			Message: "Please select a layout:",
			Options: layoutNames(),
			Description: func(value string, index int) string {
				option, ok := findLayoutOption(value)
				if !ok {
					return ""
				}
				return option.Description
			},
		}
		err := survey.AskOne(prompt, &layout)
		if err != nil {
			return false, err
		}
		if option, ok := findLayoutOption(layout); ok {
			repo = option.Repo
		}
		err = os.RemoveAll(p.ProjectName)
		if err != nil {
			return false, fmt.Errorf("remove old project: %w", err)
		}
	} else {
		repo = repoURL
	}

	fmt.Printf("git clone %s\n", repo)
	cmd := exec.Command("git", "clone", repo, p.ProjectName)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("git clone %s: %w", repo, err)
	}
	return true, nil
}

func (p *Project) replacePackageName() error {
	packageName := helper.GetProjectName(p.ProjectName)
	if packageName == "" {
		return fmt.Errorf("read module name from %s/go.mod", p.ProjectName)
	}

	err := p.replaceFiles(packageName)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "edit", "-module", p.ProjectName)
	cmd.Dir = p.ProjectName
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod edit: %w", err)
	}
	return nil
}
func (p *Project) modTidy() error {
	fmt.Println("go mod tidy")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = p.ProjectName
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}
	return nil
}
func (p *Project) rmGit() {
	os.RemoveAll(p.ProjectName + "/.git")
}
func (p *Project) installWire() error {
	fmt.Printf("go install %s\n", config.WireCmd)
	cmd := exec.Command("go", "install", config.WireCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go install %s: %w", config.WireCmd, err)
	}
	return nil
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
		return fmt.Errorf("walk file: %w", err)
	}
	return nil
}
