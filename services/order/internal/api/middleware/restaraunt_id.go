package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
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
				log.ErrorContext(r.Context(), "X-Restaurant-ID header is missing")
				common.RespondJSON(
					log,
					w,
					http.StatusBadRequest,
					common.NewError(common.BadRequestErrorCode, "X-Restaurant-ID header is required", nil),
				)
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
