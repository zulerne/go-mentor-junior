package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

const (
	RestaurantIDKey ContextKey = "restaurant_id"

	restaurantIDHeader = "X-Restaurant-ID"
)

func RestaurantID(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.DebugContext(r.Context(), "Checking for X-Restaurant-ID header")

			h := r.Header.Get(restaurantIDHeader)

			if h == "" {
				log.WarnContext(r.Context(), "X-Restaurant-ID header is missing")
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte("X-Restaurant-ID header is required"))
				if err != nil {
					log.ErrorContext(r.Context(), "Failed to write response", "error", err)
				}
				return
			}

			ctx := context.WithValue(r.Context(), RestaurantIDKey, h)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetRestaurantID(ctx context.Context) string {
	if id, ok := ctx.Value(RestaurantIDKey).(string); ok {
		return id
	}
	return ""
}
