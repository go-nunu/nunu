package create

import (
	"fmt"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/go-nunu/nunu/internal/pkg/helper"
	"github.com/go-nunu/nunu/tpl"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type Create struct {
	ProjectName          string
	CreateType           string
	FilePath             string
	FileName             string
	StructName           string
	StructNameLowerFirst string
	StructNameFirstChar  string
	StructNameSnakeCase  string
	IsFull               bool
}

func NewCreate() *Create {
	return &Create{}
}

var CmdCreate = &cobra.Command{
	Use:     "create [type] [handler-name]",
	Short:   "Create a new handler/service/repository/model",
	Example: "nunu create handler user",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

	},
}
var (
	tplPath string
)

func init() {
	CmdCreateHandler.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")
	CmdCreateService.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")
	CmdCreateRepository.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")
	CmdCreateModel.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")
	CmdCreateAll.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")

}

var CmdCreateHandler = &cobra.Command{
	Use:     "handler",
	Short:   "Create a new handler",
	Example: "nunu create handler user",
	Args:    cobra.ExactArgs(1),
	RunE:    runCreate,
}
var CmdCreateService = &cobra.Command{
	Use:     "service",
	Short:   "Create a new service",
	Example: "nunu create service user",
	Args:    cobra.ExactArgs(1),
	RunE:    runCreate,
}
var CmdCreateRepository = &cobra.Command{
	Use:     "repository",
	Short:   "Create a new repository",
	Example: "nunu create repository user",
	Args:    cobra.ExactArgs(1),
	RunE:    runCreate,
}
var CmdCreateModel = &cobra.Command{
	Use:     "model",
	Short:   "Create a new model",
	Example: "nunu create model user",
	Args:    cobra.ExactArgs(1),
	RunE:    runCreate,
}
var CmdCreateAll = &cobra.Command{
	Use:     "all",
	Short:   "Create a new handler & service & repository & model",
	Example: "nunu create all user",
	Args:    cobra.ExactArgs(1),
	RunE:    runCreate,
}

func runCreate(cmd *cobra.Command, args []string) error {
	c := NewCreate()
	c.ProjectName = helper.GetProjectName(".")
	if c.ProjectName == "" {
		return fmt.Errorf("read module name from go.mod")
	}
	c.CreateType = cmd.Use
	c.FilePath, c.StructName = filepath.Split(args[0])
	c.FileName = strings.ReplaceAll(c.StructName, ".go", "")
	if c.FileName == "" {
		return fmt.Errorf("component name cannot be empty")
	}
	c.StructName = strutil.UpperFirst(strutil.CamelCase(c.FileName))
	c.StructNameLowerFirst = strutil.LowerFirst(c.StructName)
	c.StructNameFirstChar = string(c.StructNameLowerFirst[0])
	c.StructNameSnakeCase = strutil.SnakeCase(c.StructName)

	switch c.CreateType {
	case "handler", "service", "repository", "model":
		return c.genFile()
	case "all":
		c.CreateType = "handler"
		if err := c.genFile(); err != nil {
			return err
		}

		c.CreateType = "service"
		if err := c.genFile(); err != nil {
			return err
		}

		c.CreateType = "repository"
		if err := c.genFile(); err != nil {
			return err
		}

		c.CreateType = "model"
		return c.genFile()
	default:
		return fmt.Errorf("invalid handler type: %s", c.CreateType)
	}
}
func (c *Create) genFile() error {
	filePath := c.FilePath
	if filePath == "" {
		filePath = fmt.Sprintf("internal/%s/", c.CreateType)
	}
	f, err := createFile(filePath, strings.ToLower(c.FileName)+".go")
	if err != nil {
		return err
	}
	if f == nil {
		log.Printf("warn: file %s%s %s", filePath, strings.ToLower(c.FileName)+".go", "already exists.")
		return nil
	}
	defer f.Close()
	var t *template.Template
	if tplPath == "" {
		t, err = template.ParseFS(tpl.CreateTemplateFS, fmt.Sprintf("create/%s.tpl", c.CreateType))
	} else {
		t, err = template.ParseFiles(path.Join(tplPath, fmt.Sprintf("%s.tpl", c.CreateType)))
	}
	if err != nil {
		return fmt.Errorf("create %s: %w", c.CreateType, err)
	}
	err = t.Execute(f, c)
	if err != nil {
		return fmt.Errorf("create %s: %w", c.CreateType, err)
	}
	log.Printf("Created new %s: %s", c.CreateType, filePath+strings.ToLower(c.FileName)+".go")
	return nil
}
func createFile(dirPath string, filename string) (*os.File, error) {
	filePath := filepath.Join(dirPath, filename)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("create dir %s: %w", dirPath, err)
	}
	_, err = os.Stat(filePath)
	if err == nil {
		return nil, nil
	}
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("stat file %s: %w", filePath, err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("create file %s: %w", filePath, err)
	}

	return file, nil
}
