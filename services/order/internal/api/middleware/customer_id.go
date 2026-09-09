package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

const (
	CustomerIDKey    = "customer_id"
	customerIDHeader = "X-Customer-ID"
)

func CustomerID(log *slog.Logger) Middleware {
	log.Debug("CustomerID middleware initialized")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("Checking for X-Customer-ID header")

			h := r.Header.Get(customerIDHeader)

			if h == "" {
				log.Warn("X-Customer-ID header is missing")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("X-Customer-ID header is required"))
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
