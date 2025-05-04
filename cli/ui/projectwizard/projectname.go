package projectwizard

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectNameModel struct {
	input textinput.Model
	err   error
}

func NewProjectNameModel() *ProjectNameModel {
	ti := textinput.New()
	ti.Placeholder = "Project Name"
	ti.Focus()
	ti.CharLimit = 255
	ti.Width = 20

	return &ProjectNameModel{
		input: ti,
		err:   nil,
	}
}

func (m *ProjectNameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ProjectNameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case error:
		m.err = msg
		return m, nil
	}

	_, cmd = m.input.Update(msg)

	return m, cmd
}

func (m *ProjectNameModel) View() string {
	return fmt.Sprintf(
		"Choose a project name:\n\n%s\n\n%s",
		m.input.View(),
		"(press esc to quit)",
	)
}
