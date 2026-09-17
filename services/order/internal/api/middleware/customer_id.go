package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
)

// TODO (review): Do I need these middleware(customer/restaurant ids) or should I extract them manually from the request?

const (
	CustomerIDKey ContextKey = "customer_id"

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

			ctx := context.WithValue(r.Context(), CustomerIDKey, h)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetCustomerID(ctx context.Context) string {
	if id, ok := ctx.Value(CustomerIDKey).(string); ok {
		return id
	}
	return ""
}
