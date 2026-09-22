package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/api/v1/handler"
	"github.com/zulerne/go-mentor-junior/order/internal/config"
	"github.com/zulerne/go-mentor-junior/order/internal/logger"
	deliveryProvider "github.com/zulerne/go-mentor-junior/order/internal/provider/delivery"
	restaurantProvider "github.com/zulerne/go-mentor-junior/order/internal/provider/restaurant"
	store "github.com/zulerne/go-mentor-junior/order/internal/provider/store"
	"github.com/zulerne/go-mentor-junior/order/internal/services/customer"
	"github.com/zulerne/go-mentor-junior/order/internal/services/restaurant"
	"golang.org/x/sync/errgroup"
)

func main() {
	validate := validator.New(validator.WithRequiredStructEnabled())

	cfg, err := config.Load(validate)
	if err != nil {
		log.Printf("failed to load config: %v", err)
		os.Exit(1)
	}

	log := logger.New(cfg.Env)

	orderStore := store.NewMemoryOrderStore()
	cartStore := store.NewMemoryCartStore()
	restaurantProvider := restaurantProvider.NewMemoryProvider()

	customer := customer.New(orderStore, cartStore, restaurantProvider, log)
	restaurant := restaurant.New(orderStore, deliveryProvider.NewStubProvider(), log)

	h := handler.New(customer, restaurant, validate, log)
	srv := &http.Server{
		Addr:         cfg.HTTPConfig.Address,
		Handler:      h.Routes(),
		WriteTimeout: cfg.HTTPConfig.Timeout,
		ReadTimeout:  cfg.HTTPConfig.Timeout,
		IdleTimeout:  cfg.HTTPConfig.IdleTimeout,
	}

	if err = run(log, srv, cfg.HTTPConfig.ShutdownTimeout); err != nil {
		log.Error("server error", "error", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger, srv *http.Server, shutdownTimeout time.Duration) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Info("starting server", "address", srv.Addr)

		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			closeErr := srv.Close()
			if closeErr != nil {
				return errors.Join(err, closeErr)
			}
			return err
		}

		log.Info("server stopped gracefully")
		return nil
	})

	return g.Wait()
}
