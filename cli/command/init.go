package command

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	projectwizardui "github.com/up1io/muxo/cli/ui/projectwizard"
	"github.com/up1io/muxo/cli/wizard/project"
	"log"
	"os"
)

type InitCmd struct {
	cmd *cobra.Command
}

type InitCmdOption func(command *InitCmd)

func NewInitCmd(opts ...InitCmdOption) *InitCmd {
	instance := &InitCmd{}

	rootCmd := &cobra.Command{
		Use:   "init",
		Short: "Init a new project or module",
	}

	initAppCmd := &cobra.Command{
		Use:   "app",
		Short: "Init a new project",
		Run:   instance.runProjectInit,
	}

	rootCmd.AddCommand(initAppCmd)

	instance.cmd = rootCmd

	for _, opt := range opts {
		opt(instance)
	}

	return instance
}

func WithRootCommand(root *cobra.Command) InitCmdOption {
	return func(command *InitCmd) {
		root.AddCommand(command.cmd)
	}
}

func (c *InitCmd) runProjectInit(cmd *cobra.Command, args []string) {
	wizardForm := projectwizardui.NewProjectWizard()
	p := tea.NewProgram(wizardForm)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	rootDir := "."
	if err := os.MkdirAll(fmt.Sprintf("%s/%s", rootDir, wizardForm.Input.ProjectName), 0750); err != nil {
		log.Fatal(err)
	}

	ctx := project.NewConfigContext(context.TODO(), &project.Config{
		ProjectDir:  fmt.Sprintf("%s/%s", rootDir, wizardForm.Input.ProjectName),
		ProjectName: wizardForm.Input.ProjectName,
		ModName:     wizardForm.Input.ModName,
		MuxoVersion: wizardForm.Input.MuxoVersion,
	})

	if err := project.Execute(ctx); err != nil {
		log.Fatal(err)
	}
}
