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
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	orderStore := store.NewOrderStore()

	customer := customer.New(orderStore, store.NewCartStore(), restaurantProvider.NewRestaurantProvider(), log)
	restaurant := restaurant.New(orderStore, deliveryProveder.NewDeliveryProvider(), log)

	srv := &http.Server{
		Addr:         "localhost:8080",
		Handler:      handler.New(customer, restaurant, log),
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
