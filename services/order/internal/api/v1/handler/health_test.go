package handler_test

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zulerne/go-mentor-junior/order/internal/api/v1/handler"
)

func TestHealthEndpoints(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name       string
		path       string
		pingErr    error
		wantStatus int
		wantPings  int
	}{
		{"ready_database_available", "/readyz", nil, http.StatusOK, 1},
		{"ready_database_unavailable", "/readyz", errors.New("database unavailable"), http.StatusServiceUnavailable, 1},
		{"live_database_unavailable", "/livez", errors.New("database unavailable"), http.StatusOK, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			pings := 0
			db := databaseStub(func(context.Context) error {
				pings++
				return tc.pingErr
			})
			h := handler.New(db, nil, nil, validator.New(), slog.New(slog.DiscardHandler))
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			h.Routes().ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Equal(t, tc.wantPings, pings)
		})
	}
}

func TestReadyz_RequestCancelled(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var pingErr error
	db := databaseStub(func(ctx context.Context) error {
		pingErr = ctx.Err()
		return pingErr
	})
	h := handler.New(db, nil, nil, validator.New(), slog.New(slog.DiscardHandler))
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/readyz", nil)
	h.Routes().ServeHTTP(rec, req)

	require.Error(t, pingErr)
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}
