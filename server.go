package muxo

import "net/http"

// Server defines the interface for HTTP servers.
type Server interface {
	// Init initializes the server.
	Init() error

	// Mux returns the HTTP mux for the server.
	Mux() http.ServeMux

	// Shutdown gracefully shuts down the server.
	// It returns a slice of errors that occurred during shutdown.
	Shutdown() []error
}
