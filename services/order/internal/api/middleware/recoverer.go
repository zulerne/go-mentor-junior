package middleware

import (
	"log/slog"
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
)

func Recoverer(log *slog.Logger) Middleware {
	log.Debug("Recoverer middleware initialized")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.ErrorContext(r.Context(), "recovered from panic", "error", err)
					common.RespondJSON(
						log,
						w,
						http.StatusInternalServerError,
						response.NewError(response.BaseErrorCode, "internal server error", nil),
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
