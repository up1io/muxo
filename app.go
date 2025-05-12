package muxo

import (
	"context"
	"github.com/up1io/muxo/middleware"
	localMiddleware "github.com/up1io/muxo/module/local/middleware"
	"github.com/up1io/muxo/runtime"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type App struct {
	srv         Server
	runtime     runtime.Runtime
	middlewares []middleware.Middleware
}

type AppOption func(app *App)

func NewApp(opts ...AppOption) (*App, error) {
	app := &App{
		runtime: runtime.NewDefaultRuntime(":8080"),
		middlewares: []middleware.Middleware{
			localMiddleware.WithLocalization("web/locales"),
		},
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

// WithMiddleware allows users to override the default middleware stack
func WithMiddleware(middlewares ...middleware.Middleware) AppOption {
	return func(app *App) {
		app.middlewares = middlewares
	}
}

// WithAdditionalMiddleware allows users to add middleware to the default stack
func WithAdditionalMiddleware(middlewares ...middleware.Middleware) AppOption {
	return func(app *App) {
		app.middlewares = append(app.middlewares, middlewares...)
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

		mux := app.srv.Mux()

		// Apply middleware to the mux
		// This uses the middleware stack configured in the App struct
		// By default, this includes core modules like localization
		// Users can override or add to this stack using WithMiddleware or WithAdditionalMiddleware
		withMiddlewares := middleware.CreateStack(app.middlewares...)

		handler := withMiddlewares(&mux)

		wrappedMux := http.NewServeMux()
		wrappedMux.Handle("/", handler)

		if err := app.runtime.Serve(ctx, *wrappedMux); err != nil {
			log.Printf("unable to run server %s", err.Error())
			stop <- os.Interrupt
		}
	}()

	<-stop
	return nil
}
