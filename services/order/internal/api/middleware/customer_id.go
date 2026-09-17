package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
)

// TODO (review): Do I need these middleware(customer/restaurant ids) or should I extract them manually from the request?

const (
	customerIDKey contextKey = "customer_id"

	customerIDHeader = "X-Customer-ID"
)

func CustomerID(log *slog.Logger) Middleware {
	log.Debug("CustomerID middleware initialized")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.DebugContext(r.Context(), "Checking for X-Customer-ID header")

			h := r.Header.Get(customerIDHeader)

			if h == "" {
				log.ErrorContext(r.Context(), "X-Customer-ID header is missing")

				common.RespondJSON(
					log,
					w,
					http.StatusBadRequest,
					common.NewError(common.BadRequestErrorCode, "X-Customer-ID header is required", nil),
				)
				return
			}

			ctx := WithCustomerID(r.Context(), h)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithCustomerID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, customerIDKey, id)
}

func CustomerIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(customerIDKey).(string)
	return id, ok
}
