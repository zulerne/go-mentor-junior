package middleware

import (
	"log/slog"
	"net/http"
)

func CustomerID(log *slog.Logger) Middleware {
	log.Debug("CustomerID middleware initialized")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("Checking for X-Customer-ID header")

			h := r.Header.Get("X-Customer-ID")

			if h == "" {
				log.Warn("X-Customer-ID header is missing")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("X-Customer-ID header is required"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
