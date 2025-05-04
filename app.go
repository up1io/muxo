package muxo

import (
	"context"
	"github.com/up1io/muxo/runtime"
	"log"
	"os"
	"os/signal"
)

type App struct {
	srv     Server
	runtime runtime.Runtime
}

type AppOption func(app *App)

func NewApp(opts ...AppOption) (*App, error) {
	app := &App{
		runtime: runtime.NewDefaultRuntime(":8080"),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app, nil
}

func WithRuntime(runtime runtime.Runtime) AppOption {
	return func(app *App) {
		app.runtime = runtime
	}
}

func WithServer(srv Server) AppOption {
	return func(app *App) {
		app.srv = srv
	}
}

func (app *App) Serve() error {
	ctx := context.Background()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go func() {
		defer func() {
			errs := app.srv.Shutdown()
			if len(errs) != 0 {
				log.Printf("errors occured during server shutdown: %v", errs)
			}
		}()

		if err := app.srv.Init(); err != nil {
			log.Printf("fail to configure server %s", err.Error())
			stop <- os.Interrupt
		}

		if err := app.runtime.Serve(ctx, app.srv.Mux()); err != nil {
			log.Printf("unable to run server %s", err.Error())
			stop <- os.Interrupt
		}
	}()

	<-stop
	return nil
}
