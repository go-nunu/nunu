package new

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle   = lipgloss.NewStyle().Margin(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
)

type listModel struct {
	list     list.Model
	choice   string
	selected bool
	err      error
}

type listItem struct {
	title, desc string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }

func newListModel(items []list.Item, title string) *listModel {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = title
	l.Styles.Title = titleStyle
	return &listModel{
		list:     l,
		selected: false,
	}
}

func (m *listModel) Init() tea.Cmd {
	return nil
}

func (m *listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.err = errors.New("quit nunu project-layout")
			return m, tea.Quit
		case tea.KeyEnter:
			if i, ok := m.list.SelectedItem().(listItem); ok {
				m.choice = i.Title()
				m.selected = true
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *listModel) View() string {
	if m.selected {
		return fmt.Sprintf("You chose %s\n", m.choice)
	}
	return docStyle.Render(m.list.View())
}
