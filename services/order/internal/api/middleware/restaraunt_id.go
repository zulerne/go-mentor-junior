package middleware

import (
	"log/slog"
	"net/http"
)

func RestaurantID(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("Checking for X-Restaurant-ID header")

			h := r.Header.Get("X-Restaurant-ID")

			if h == "" {
				log.Warn("X-Restaurant-ID header is missing")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("X-Restaurant-ID header is required"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
