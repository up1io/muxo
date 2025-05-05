package cli

import (
	"github.com/spf13/cobra"
	"github.com/up1io/muxo/cli/command"
)

type Runner struct {
	rootCmd *cobra.Command
}

func New() *Runner {
	rootCmd := &cobra.Command{
		Use:   "muxo",
		Short: "Server Side-Rendering (SSR) Go Web Framework",
		Long:  "Muxo is a primarily Server Side-Rendering (SSR) Go Web Framework.",
	}

	command.NewInitCmd(rootCmd)
	command.NewDevCommand(rootCmd)

	return &Runner{
		rootCmd: rootCmd,
	}
}

func (r *Runner) Run() error {
	return r.rootCmd.Execute()
}
