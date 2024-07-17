package new

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type textInputModel struct {
	textInput   textinput.Model
	ProjectName string
	err         error
}

func InitTextInputModel() *textInputModel {
	ti := textinput.New()
	ti.Placeholder = "Enter your project name"
	ti.Focus()
	ti.PlaceholderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Italic(true)

	ti.TextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")). // White color
		PaddingLeft(1).
		PaddingRight(1).
		Italic(true)

	ti.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		PaddingLeft(1).
		PaddingRight(1)

	ti.CharLimit = 156
	ti.Width = 40

	return &textInputModel{
		textInput:   ti,
		ProjectName: "",
		err:         nil,
	}
}

func (m *textInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.ProjectName = m.textInput.Value()
			return m, tea.Quit
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *textInputModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32CD32")). // LimeGreen color
		Bold(true).
		Render

	quitTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#32CD32")). // LimeGreen color
		Bold(true).
		Render

	return fmt.Sprintf(
		"✨  %s\n%s\n%s😶‍🌫️",
		titleStyle("What is your project name?"),
		m.textInput.View(),
		quitTextStyle("Press ctrl+c to quit!"),
	) + "\n"
}
