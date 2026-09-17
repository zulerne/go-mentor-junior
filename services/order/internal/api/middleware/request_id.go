package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
)

const (
	// requestIDKey is the context key for request ID.
	requestIDKey ContextKey = "request_id"
	// RequestIDHeader is the header name for request ID.
	requestIDHeader = "X-Request-ID"
)

// TODO: Use standard library for tracing (https://www.cncf.io/projects)
type requestIDGenerator struct {
	atomic.Uint64
}

// RequestID middleware adds a unique request ID to each request context.
// If X-Request-ID header is present in the request, it uses that value.
// Otherwise, it generates a new unique ID.
func RequestID(log *slog.Logger) Middleware {
	log.Debug("RequestID middleware initialized")
	reqid := &requestIDGenerator{}

	return func(next http.Handler) http.Handler {
		log.Debug("RequestID middleware enabled")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				myid := reqid.Add(1)
				requestID = fmt.Sprintf("%016x", myid)
			}

			// Set the request ID in response header
			w.Header().Set(requestIDHeader, requestID)

			ctx = WithRequestID(ctx, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}
