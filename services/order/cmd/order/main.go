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

	"github.com/zulerne/go-mentor-junior/order/internal/api/handler"
	"github.com/zulerne/go-mentor-junior/order/internal/config"
	"github.com/zulerne/go-mentor-junior/order/internal/logger"
	deliveryProveder "github.com/zulerne/go-mentor-junior/order/internal/provider/delivery/memory"
	restaurantProvider "github.com/zulerne/go-mentor-junior/order/internal/provider/restaurant/memory"
	store "github.com/zulerne/go-mentor-junior/order/internal/provider/store/memory"
	"github.com/zulerne/go-mentor-junior/order/internal/services/customer"
	"github.com/zulerne/go-mentor-junior/order/internal/services/restaurant"
)

// TODO (review): Global questions:
// 1. Should we use global logger? Or pass it everywhere(like now)?
// 2. Is project structure correct?
// 3. Is it okay that I splitted service into two separate services (customer and restaurant)?
// 4. What to use: internal/services/customer and internal/services/restaurant or just internal/customer and internal/restaurant?

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("failed to load config: %v", err)
		os.Exit(1)
	}

	log := logger.New(cfg.Env)

	orderStore := store.NewOrderStore()
	cartStore := store.NewCartStore()
	restaurantProvider := restaurantProvider.NewRestaurantProvider()

	customer := customer.New(orderStore, cartStore, restaurantProvider, log)
	restaurant := restaurant.New(orderStore, deliveryProveder.NewDeliveryProvider(), log)

	h := handler.New(customer, restaurant, log)
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
	errCh := make(chan error, 1)
	go func() {
		log.Info("starting server", "address", srv.Addr)

		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			closeErr := srv.Close()
			if closeErr != nil {
				return errors.Join(err, closeErr)
			}
			return err
		} else {
			log.Info("server stopped gracefully")
		}

		if err := <-errCh; err != nil {
			return err
		}
	}

	return nil
}
