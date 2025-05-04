package muxo

import (
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"net/http"
)

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func EncodeRender(w http.ResponseWriter, r *http.Request, v templ.Component) error {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	if err := v.Render(r.Context(), w); err != nil {
		return fmt.Errorf("render: %w", err)
	}
	return nil
}
