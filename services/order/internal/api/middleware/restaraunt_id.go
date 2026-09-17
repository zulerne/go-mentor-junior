package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
)

const (
	restaurantIDKey ContextKey = "restaurant_id"

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

			ctx := WithRestaurantID(r.Context(), h)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithRestaurantID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, restaurantIDKey, id)
}

func RestaurantIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(restaurantIDKey).(string)
	return id, ok
}
