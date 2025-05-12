package muxo

import (
	"context"
	"fmt"
	"github.com/up1io/muxo/logger"
	"github.com/up1io/muxo/middleware"
	localMiddleware "github.com/up1io/muxo/module/local/middleware"
	"github.com/up1io/muxo/runtime"
	"os"
	"os/signal"
)

// App represents a web application with middleware support.
type App struct {
	srv         Server
	runtime     runtime.Runtime
	middlewares []middleware.Middleware
	log         logger.Logger
}

// AppOption is a function that configures an App.
type AppOption func(app *App)

// NewApp creates a new App with the given options.
func NewApp(opts ...AppOption) *App {
	app := &App{
		runtime: runtime.NewDefaultRuntime(":8080"),
		middlewares: []middleware.Middleware{
			localMiddleware.WithLocalization("web/locales"),
		},
		log: logger.Default,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// WithRuntime sets the runtime for the App.
func WithRuntime(runtime runtime.Runtime) AppOption {
	return func(app *App) {
		app.runtime = runtime
	}
}

// WithServer sets the server for the App.
func WithServer(srv Server) AppOption {
	return func(app *App) {
		app.srv = srv
	}
}

// WithMiddleware allows users to override the default middleware stack.
func WithMiddleware(middlewares ...middleware.Middleware) AppOption {
	return func(app *App) {
		app.middlewares = middlewares
	}
}

// WithAdditionalMiddleware allows users to add middleware to the default stack.
func WithAdditionalMiddleware(middlewares ...middleware.Middleware) AppOption {
	return func(app *App) {
		app.middlewares = append(app.middlewares, middlewares...)
	}
}

// WithLogger allows users to provide a custom logger.
func WithLogger(log logger.Logger) AppOption {
	return func(app *App) {
		app.log = log
	}
}

// Serve initializes and starts the server, applying middleware and handling graceful shutdown.
func (app *App) Serve() error {
	if app.srv == nil {
		return fmt.Errorf("server is not configured, use WithServer option")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		defer func() {
			errs := app.srv.Shutdown()
			if len(errs) != 0 {
				app.log.Error("errors occurred during server shutdown: %v", errs)
			}
		}()

		if err := app.srv.Init(); err != nil {
			app.log.Error("failed to initialize server: %s", err.Error())
			errCh <- err
			return
		}

		mux := app.srv.Mux()

		// Apply middleware to the mux
		// This uses the middleware stack configured in the App struct
		// By default, this includes core modules like localization
		// Users can override or add to this stack using WithMiddleware or WithAdditionalMiddleware
		withMiddlewares := middleware.CreateStack(app.middlewares...)
		handler := withMiddlewares(&mux)

		if err := app.runtime.Serve(ctx, handler); err != nil {
			app.log.Error("failed to run server: %s", err.Error())
			errCh <- err
			return
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}
