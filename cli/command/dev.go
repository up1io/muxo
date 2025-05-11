package command

import (
	"github.com/spf13/cobra"
	"github.com/up1io/muxo/locales"
	"github.com/up1io/muxo/processor"
	"github.com/up1io/muxo/templater"
	"github.com/up1io/muxo/watcher"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
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
	p := processor.New()

	builder := &locales.Builder{Root: "web/locales"}
	templ := &templater.Templater{Dir: "template"}

	p.Add(builder)
	p.Add(templ)

	p.Run()

	restartCh := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1)

	var (
		debounceMu    sync.Mutex
		debounceTimer *time.Timer
	)

	scheduleRestart := func() {
		debounceMu.Lock()
		defer debounceMu.Unlock()
		if debounceTimer != nil {
			debounceTimer.Stop()
		}
		debounceTimer = time.AfterFunc(500*time.Millisecond, func() {
			select {
			case restartCh <- struct{}{}:
				log.Println("queued app restart after debounce")
			default:
				log.Println("restart already pending; skipping")
			}
		})
	}

	onChange := func(path string) {
		ext := filepath.Ext(path)
		switch ext {
		case ".po":
			if err := builder.Process(); err != nil {
				log.Fatal(err)
			}
			scheduleRestart()
		case ".templ":
			if err := templ.Process(); err != nil {
				log.Fatal(err)
			}
		case ".go":
			log.Println("Go file changed, consider restarting the app")
			scheduleRestart()
		default:
			// ignore
		}

	}

	fileWatcher, err := watcher.NewWatcher(".", []string{".po", ".templ", ".go"}, onChange)
	if err != nil {
		log.Fatalf("unable to create watcher: %s", err)
	}
	if err := fileWatcher.AddDir("."); err != nil {
		log.Fatalf("unable to watch directories: %s", err)
	}

	go fileWatcher.Run()

	go supervise(restartCh)

	go func() {
		for range sigs {
			scheduleRestart()
		}
	}()

	select {}
}

// supervise runs the app, kills its process group on restart, and loops
func supervise(restart <-chan struct{}) {
	cmd := runApp()

	for range restart {
		pgid := cmd.Process.Pid
		if err := syscall.Kill(-pgid, syscall.SIGINT); err != nil {
			log.Printf("failed to send SIGINT to group: %v", err)
		}

		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			log.Printf("failed to kill process: %v", err)
		}

		time.Sleep(200 * time.Millisecond)

		exited := make(chan error, 1)
		go func() {
			println("closed program")
			exited <- cmd.Wait()
		}()

		select {
		case err := <-exited:
			if err != nil {
				log.Printf("process exited with error: %v", err)
			}
		case <-time.After(5 * time.Second):
			log.Println("timeout waiting for graceful shutdown, killing")
			if err := cmd.Process.Kill(); err != nil {
				log.Fatalf("failed to kill process: %v", err)
			}
			<-exited
		}

		time.Sleep(200 * time.Millisecond)

		cmd = runApp()
	}
}

// runApp starts the Go application and returns the *exec.Cmd
func runApp() *exec.Cmd {
	cmd := exec.Command("go", "run", "cmd/local/main.go")

	// ensure subprocesses die with parent
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	if err := cmd.Start(); err != nil {
		log.Fatalf("failed to start app: %s", err)
	}

	log.Println("app started with PID", cmd.Process.Pid)
	return cmd
}
