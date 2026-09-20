package middleware

import (
	"net/http"
	"slices"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range slices.Backward(middlewares) {
		h = middleware(h)
	}
	return h
}
