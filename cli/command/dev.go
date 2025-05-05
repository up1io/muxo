package command

import (
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

type DevCommand struct {
	cmd *cobra.Command
}

func NewDevCommand(rootCmd *cobra.Command) *DevCommand {
	instance := &DevCommand{}

	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Launch a local application instance in dev-mode.",
		Run:   instance.run,
	}

	instance.cmd = cmd

	rootCmd.AddCommand(cmd)

	return instance
}

func (d *DevCommand) run(cmd *cobra.Command, args []string) {
	// Todo: Make sure the binaries are available

	// Note(john): Step 1 generate templates
	templateCmd := exec.Command("templ", "generate", ".")
	if err := templateCmd.Run(); err != nil {
		log.Fatal(err)
	}

	// Note(john): Step 2 build locales
	localCmd := exec.Command("msgfmt", "web/locales/en/default.po", "-o", "web/locales/en/default.mo")
	if err := localCmd.Run(); err != nil {
		log.Fatalf("unable to generate locales. %s", err)
	}

	// Note(john): Step 3 run local
	// Todo: add a file watcher
	runCmd := exec.Command("go", "run", "cmd/local/main.go")
	if err := runCmd.Run(); err != nil {
		log.Fatalf("unable to find local main entrypoint")
	}
}
