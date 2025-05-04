package projectwizard

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/up1io/muxo/cli/wizard/project"
)

type state int

const (
	projectName state = iota
	moduleName
)

type ProjectWizard struct {
	Input project.Config

	state state

	projectName *ProjectNameModel
	moduleName  *ModuleNameModel
}

func NewProjectWizard() *ProjectWizard {
	projectNameModel := NewProjectNameModel()
	moduleNameModel := NewModuleNameModel()

	return &ProjectWizard{
		Input: project.Config{
			MuxoVersion: "v0.0.1",
		},

		projectName: projectNameModel,
		moduleName:  moduleNameModel,
	}
}

func (m *ProjectWizard) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ProjectWizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case projectName:
		input, updateCmd := m.projectName.input.Update(msg)
		m.projectName.input = input
		m.Input.ProjectName = input.Value()
		cmd = updateCmd
	case moduleName:
		input, updateCmd := m.moduleName.input.Update(msg)
		m.moduleName.input = input
		m.Input.ModName = input.Value()
		cmd = updateCmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.state == projectName {
				m.state = moduleName
			} else {
				return m, tea.Quit
			}
		}
	}

	return m, cmd
}

func (m *ProjectWizard) View() string {
	switch m.state {
	case projectName:
		return m.projectName.View()
	case moduleName:
		return m.moduleName.View()
	default:
		return "unknown"
	}
}
