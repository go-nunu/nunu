package wire

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type choiceModel struct {
	choices  []string
	cursor   int
	choice   string
	quitting bool
	title    string
}

func newChoiceModel(choices []string, title string) choiceModel {
	return choiceModel{
		choices: choices,
		title:   title,
	}
}

func (m choiceModel) Init() tea.Cmd {
	return nil
}

func (m choiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.choice = ""
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.choice = m.choices[m.cursor]
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m choiceModel) View() string {
	if m.quitting {
		return ""
	}
	s := fmt.Sprintf("\033[1;32m%s\033[0m %s\033[1;34m%s\033[0m\n", "?", m.title, "  [Use arrows to move, type to filter]")

	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "\033[32m>\033[0m"
			s += fmt.Sprintf("%s \033[1;33m%s\033[0m\n", cursor, choice)
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
	}

	s += "\n\033[2mPress q to quit.\033[0m"
	return s
}
