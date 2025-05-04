package projectwizard

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/up1io/muxo/cli/wizard/project"
)

type state int

const (
	projectName state = iota
)

type ProjectWizard struct {
	Input project.Config

	state state

	projectName *ProjectNameModel
}

func NewProjectWizard() *ProjectWizard {
	pn := NewProjectNameModel()

	return &ProjectWizard{
		Input: project.Config{
			MuxoVersion: "v0.0.1",
			ModName:     "example.com/project",
		},

		projectName: pn,
	}
}

func (m *ProjectWizard) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ProjectWizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			return m, tea.Quit
		}
	}

	switch m.state {
	case projectName:
		input, updateCmd := m.projectName.input.Update(msg)
		m.projectName.input = input
		m.Input.ProjectName = input.Value()
		cmd = updateCmd
	}

	return m, cmd
}

func (m *ProjectWizard) View() string {
	switch m.state {
	case projectName:
		return m.projectName.View()
	default:
		return "unknown"
	}
}
