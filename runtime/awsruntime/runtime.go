package awsruntime

import (
	"context"
	"github.com/akrylysov/algnhsa"
	"net/http"
)

type AwsRuntime struct{}

func New() *AwsRuntime {
	return &AwsRuntime{}
}

func (r *AwsRuntime) Serve(ctx context.Context, mux http.ServeMux) error {
	algnhsa.ListenAndServe(&mux, nil)
	return nil
}
