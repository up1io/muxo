// Package middleware provides utilities for HTTP middleware.
package middleware

import "net/http"

// Middleware is a function that wraps an http.Handler with additional functionality.
type Middleware func(http.Handler) http.Handler

// CreateStack creates a middleware stack from the given middleware functions.
// The middleware are applied in reverse order, so the first middleware in the list
// is the outermost middleware (the first to receive the request and the last to
// handle the response).
func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}
