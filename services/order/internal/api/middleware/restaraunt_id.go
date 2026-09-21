package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
)

const (
	restaurantIDKey contextKey = "restaurant_id"

	restaurantIDHeader = "X-Restaurant-ID"
)

func RestaurantID(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.DebugContext(r.Context(), "Checking for X-Restaurant-ID header")

			restaurantID, err := uuid.Parse(r.Header.Get(restaurantIDHeader))

			if err != nil {
				log.ErrorContext(r.Context(), "X-Restaurant-ID header is missing")
				common.RespondJSON(
					log,
					w,
					http.StatusBadRequest,
					common.NewError(common.BadRequestErrorCode, "X-Restaurant-ID header is required", nil),
				)
				return
			}

			ctx := WithRestaurantID(r.Context(), restaurantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithRestaurantID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, restaurantIDKey, id)
}

func RestaurantIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(restaurantIDKey).(uuid.UUID)
	return id, ok
}
