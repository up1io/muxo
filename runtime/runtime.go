// Package runtime provides functionality for running HTTP servers.
package runtime

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Runtime defines the interface for server runtimes.
type Runtime interface {
	// Serve starts the HTTP server with the given handler and handles graceful shutdown.
	Serve(ctx context.Context, handler http.Handler) error
}

// DefaultRuntime is a basic implementation of the Runtime interface.
type DefaultRuntime struct {
	addr string
}

// NewDefaultRuntime creates a new DefaultRuntime with the given address.
func NewDefaultRuntime(addr string) *DefaultRuntime {
	return &DefaultRuntime{
		addr: addr,
	}
}

// Serve starts an HTTP server on the configured address with the given handler.
func (r *DefaultRuntime) Serve(ctx context.Context, handler http.Handler) error {
	addr := r.addr
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	fmt.Printf("Starting server on %s\n", addr)
	return http.ListenAndServe(addr, handler)
}
