package projectwizard

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ModuleNameModel struct {
	input textinput.Model
	err   error
}

func NewModuleNameModel() *ModuleNameModel {
	ti := textinput.New()
	ti.Placeholder = "github.com/example/project"
	ti.Focus()
	ti.CharLimit = 255
	ti.Width = 20

	return &ModuleNameModel{
		input: ti,
		err:   nil,
	}
}

func (m *ModuleNameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ModuleNameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *ModuleNameModel) View() string {
	return fmt.Sprintf(
		"Enter your module name:\n\n%s\n\n%s",
		m.input.View(),
		"(press esc to quit)",
	)
}
