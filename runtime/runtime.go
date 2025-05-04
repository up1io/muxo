package runtime

import (
	"context"
	"net/http"
)

type Runtime interface {
	Serve(ctx context.Context, mux http.ServeMux) error
}

type DefaultRuntime struct {
	addr string
}

func NewDefaultRuntime(addr string) *DefaultRuntime {
	return &DefaultRuntime{
		addr: addr,
	}
}

func (r *DefaultRuntime) Serve(ctx context.Context, mux http.ServeMux) error {
	return http.ListenAndServe(r.addr, &mux)
}
