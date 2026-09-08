package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/zulerne/go-mentor-junior/order/internal/api/handler"
	"github.com/zulerne/go-mentor-junior/order/internal/config"
	"github.com/zulerne/go-mentor-junior/order/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	log.Info("config initialized", "config", cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var storage handler.Storage = nil

	srv := &http.Server{
		Addr:         "localhost:8080",
		Handler:      handler.New(storage, log),
		WriteTimeout: cfg.HTTPConfig.Timeout,
		ReadTimeout:  cfg.HTTPConfig.Timeout,
		IdleTimeout:  cfg.HTTPConfig.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	log.Info("server started", "address", srv.Addr)
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTPConfig.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("failed to shutdown server gracefully", "error", err)
		srv.Close()
	}

	if err := <-errCh; err != nil {
		log.Error("server error", "error", err)
	}

	log.Info("server stopped gracefully")
}
