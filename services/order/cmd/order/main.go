package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zulerne/go-mentor-junior/order/internal/api/v1/handler"
	"github.com/zulerne/go-mentor-junior/order/internal/config"
	"github.com/zulerne/go-mentor-junior/order/internal/logger"
	deliveryProvider "github.com/zulerne/go-mentor-junior/order/internal/provider/delivery"
	restaurantProvider "github.com/zulerne/go-mentor-junior/order/internal/provider/restaurant"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store/cart"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store/order"
	"github.com/zulerne/go-mentor-junior/order/internal/services/customer"
	"github.com/zulerne/go-mentor-junior/order/internal/services/restaurant"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		log.Printf("application failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	cfg, err := config.Load(validate)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log := logger.New(cfg.Env)

	pool, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("create database pool: %w", err)
	}
	defer pool.Close()

	orderStore := order.New(pool)
	cartStore := cart.New(pool)
	restaurantProvider := restaurantProvider.NewMemoryProvider()

	customer := customer.New(pool, orderStore, cartStore, restaurantProvider, log)
	restaurant := restaurant.New(orderStore, deliveryProvider.NewStubProvider(), log)

	h := handler.New(pool, customer, restaurant, validate, log)
	srv := &http.Server{
		Addr:         cfg.HTTP.Address,
		Handler:      h.Routes(),
		WriteTimeout: cfg.HTTP.Timeout,
		ReadTimeout:  cfg.HTTP.Timeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	if err = serve(log, srv, cfg.HTTP.ShutdownTimeout); err != nil {
		return fmt.Errorf("serve HTTP: %w", err)
	}
	return nil
}

func serve(log *slog.Logger, srv *http.Server, shutdownTimeout time.Duration) error {
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
