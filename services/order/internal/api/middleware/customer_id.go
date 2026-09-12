package middleware

import (
	"context"
	"log/slog"
	"net/http"
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
				log.WarnContext(r.Context(), "X-Customer-ID header is missing")
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte("X-Customer-ID header is required"))
				if err != nil {
					log.ErrorContext(r.Context(), "Failed to write response", "error", err)
				}
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
