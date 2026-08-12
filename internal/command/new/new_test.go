package new

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/AlecAivazis/survey/v2"
)

func TestCloneTemplateDoesNotDeleteExistingProjectWhenLayoutSelectionFails(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "existing")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(project, "keep.txt")
	if err := os.WriteFile(marker, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}

	oldAskOne, oldRepoURL := askOne, repoURL
	repoURL = ""
	calls := 0
	askOne = func(_ survey.Prompt, response interface{}, _ ...survey.AskOpt) error {
		calls++
		if calls == 1 {
			*response.(*bool) = true
			return nil
		}
		return errors.New("selection cancelled")
	}
	t.Cleanup(func() {
		askOne = oldAskOne
		repoURL = oldRepoURL
	})

	p := &Project{ProjectName: project}
	if _, err := p.cloneTemplate(); err == nil {
		t.Fatal("cloneTemplate() returned nil error after layout selection failed")
	}
	data, err := os.ReadFile(marker)
	if err != nil {
		t.Fatalf("existing project was removed: %v", err)
	}
	if string(data) != "keep" {
		t.Fatalf("existing project marker = %q, want keep", data)
	}
}
