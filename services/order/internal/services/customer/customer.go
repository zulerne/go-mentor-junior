package customer

//go:generate go run github.com/vektra/mockery/v3@v3.8.0

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store/cart"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store/order"
)

// OrderStore is a persistent store for orders.
//
//mockery:generate: true
type OrderStore interface {
	Find(ctx context.Context, id uuid.UUID) (domain.Order, error)
	Update(ctx context.Context, order domain.Order) (domain.Order, error)
	GetAllOrders(ctx context.Context, customerID uuid.UUID) ([]domain.Order, error)
}

// CartStore is a persistent store for carts.
//
//mockery:generate: true
type CartStore interface {
	Find(ctx context.Context, customerID uuid.UUID) (domain.Cart, error)
	Update(ctx context.Context, customerID uuid.UUID, cart domain.Cart) error
}

// RestaurantProvider provides access to restaurant and menu data.
//
//mockery:generate: true
type RestaurantProvider interface {
	Find(ctx context.Context, restaurantID uuid.UUID) (domain.Restaurant, error)
	FindItem(ctx context.Context, restaurantID uuid.UUID, menuItemID uuid.UUID) (domain.MenuItem, error)
}

type Service struct {
	pool               *pgxpool.Pool
	orderStore         OrderStore
	cartStore          CartStore
	restaurantProvider RestaurantProvider
	log                *slog.Logger
	now                func() time.Time
	newID              func() uuid.UUID
}

type Option func(*Service)

// withinCart owns the transaction and locks the cart for the entire operation.
func (s *Service) withinCart(ctx context.Context, customerID uuid.UUID, fn func(CartStore, OrderStore) error) error {
	// A nil pool keeps the original memory/mock stores usable without PostgreSQL.
	if s.pool == nil {
		return fn(s.cartStore, s.orderStore)
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin cart transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(context.Background()) }()
	if _, err = tx.Exec(ctx, `INSERT INTO carts (customer_id) VALUES ($1)
  ON CONFLICT (customer_id) DO NOTHING`, customerID); err != nil {
		return fmt.Errorf("ensure cart: %w", err)
	}
	var lockedID uuid.UUID
	if err = tx.QueryRow(ctx, "SELECT customer_id FROM carts WHERE customer_id = $1 FOR UPDATE", customerID).
		Scan(&lockedID); err != nil {
		return fmt.Errorf("lock cart: %w", err)
	}
	carts := cart.NewWithTx(tx)
	orders := order.NewWithTx(tx)
	if err = fn(carts, orders); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit cart transaction: %w", err)
	}
	return nil
}

func WithNow(now func() time.Time) Option {
	return func(s *Service) { s.now = now }
}

func WithNewID(newID func() uuid.UUID) Option {
	return func(s *Service) { s.newID = newID }
}

// New uses the pool for atomic cart operations. A nil pool uses the supplied
// stores directly for in-memory operation and unit tests.
func New(
	pool *pgxpool.Pool,
	orderStore OrderStore,
	cartStore CartStore,
	restaurantProvider RestaurantProvider,
	log *slog.Logger,
	opts ...Option,
) *Service {
	log = log.With("component", "customer")
	s := &Service{
		pool:               pool,
		orderStore:         orderStore,
		cartStore:          cartStore,
		restaurantProvider: restaurantProvider,
		log:                log,
		now:                func() time.Time { return time.Now().UTC() },
		newID:              uuid.New,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
