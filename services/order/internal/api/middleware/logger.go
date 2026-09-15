package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter wraps [http.ResponseWriter] to capture the status code.
type responseWriter struct {
	http.ResponseWriter

	statusCode int
	written    bool
}

func newWrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status code
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}

// Unwrap returns the underlying ResponseWriter for [http.ResponseController] compatibility.
func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

func Logger(log *slog.Logger) Middleware {
	log.Debug("Logger middleware initialized")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := newWrapResponseWriter(w)
			start := time.Now()
			defer func() {
				log.InfoContext(r.Context(), "HTTP",
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.statusCode,
					"duration", time.Since(start),
					"request_id", GetRequestID(r.Context()))
			}()
			next.ServeHTTP(w, r)
		})
	}
}
