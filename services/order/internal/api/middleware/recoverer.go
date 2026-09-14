package middleware

import (
	"log/slog"
	"net/http"
)

func Recoverer(log *slog.Logger) Middleware {
	log.Debug("Recoverer middleware initialized")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("recovered from panic", "error", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
