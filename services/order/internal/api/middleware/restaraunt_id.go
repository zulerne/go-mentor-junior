package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

const (
	RestaurantIDKey    = "restaurant_id"
	restaurantIDHeader = "X-Restaurant-ID"
)

func RestaurantID(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("Checking for X-Restaurant-ID header")

			h := r.Header.Get(restaurantIDHeader)

			if h == "" {
				log.Warn("X-Restaurant-ID header is missing")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("X-Restaurant-ID header is required"))
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
