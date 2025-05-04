package muxo

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"net/http"
	"strings"
)

var decoder = schema.NewDecoder()

func Decode[T any](r *http.Request) (T, error) {
	var v T
	typ := r.Header.Get("Content-Type")

	if strings.HasPrefix(typ, "application/x-www-form-urlencoded") {
		if err := DecodeForm(r, &v); err != nil {
			return v, fmt.Errorf("decode form data: %w", err)
		}
		return v, nil
	}

	if strings.HasPrefix(typ, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			return v, fmt.Errorf("decode json: %w", err)
		}
		return v, nil
	}

	return v, fmt.Errorf("content type %s is not supported", typ)
}

func DecodeForm(r *http.Request, v interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	if err := decoder.Decode(v, r.PostForm); err != nil {
		return err
	}
	return nil
}

func DecodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	v, err := Decode[T](r)
	if err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}
	if problems := v.Valid(r.Context()); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return v, nil, nil
}
