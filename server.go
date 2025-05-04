package muxo

import "net/http"

type Server interface {
	Init() error
	Mux() http.ServeMux
	Shutdown() []error
}
