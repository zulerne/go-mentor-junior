//go:build integration

package customer_test

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store/cart"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store/order"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store/testdb"
	"github.com/zulerne/go-mentor-junior/order/internal/services/customer"
)

type restaurantStub struct{ onFind func(context.Context) error }

func (r restaurantStub) Find(ctx context.Context, id uuid.UUID) (domain.Restaurant, error) {
	if r.onFind != nil {
		if err := r.onFind(ctx); err != nil {
			return domain.Restaurant{}, err
		}
	}
	return domain.Restaurant{ID: id, Currency: "USD", AcceptingOrders: true}, nil
}

func (restaurantStub) FindItem(_ context.Context, _ uuid.UUID, id uuid.UUID) (domain.MenuItem, error) {
	return domain.MenuItem{ID: id, Name: "Soup", Currency: "USD", PriceMinor: 150, Available: true}, nil
}

func newService(pool *pgxpool.Pool, provider restaurantStub) *customer.Service {
	return customer.New(pool, order.New(pool), cart.New(pool), provider, slog.New(slog.DiscardHandler))
}

func seedCart(t *testing.T, pool *pgxpool.Pool, customerID uuid.UUID) domain.Cart {
	t.Helper()
	value := domain.Cart{RestaurantID: uuid.New(), Currency: "USD", SubtotalMinor: 250,
		Items: []domain.CartItem{{MenuItemID: uuid.New(), Name: "Soup", Currency: "USD",
			UnitPriceMinor: 125, Quantity: 2, Instructions: "No salt"}}}
	require.NoError(t, cart.New(pool).Update(context.Background(), customerID, value))
	return value
}

func TestCheckoutSuccess(t *testing.T) {
	pool := testdb.New(t, "../../../migrations")
	id := uuid.New()
	original := seedCart(t, pool, id)
	ctx := context.Background()
	saved, err := newService(pool, restaurantStub{}).CreateOrder(ctx, id, "123 Test Street")
	require.NoError(t, err)
	require.Equal(t, int64(300), saved.SubtotalMinor)
	require.Equal(t, original.RestaurantID, saved.RestaurantID)
	require.Equal(t, id, saved.CustomerID)
	require.Equal(t, domain.Pending, saved.Status)
	require.Equal(t, "No salt", saved.Items[0].Instructions)
	got, err := order.New(pool).Find(ctx, saved.ID)
	require.NoError(t, err)
	require.Equal(t, saved, got)
	empty, err := cart.New(pool).Find(ctx, id)
	require.NoError(t, err)
	require.Equal(t, domain.Cart{}, empty)
}

func TestCheckoutRollback(t *testing.T) {
	for _, table := range []string{"order_items", "cart_items"} {
		t.Run(table, func(t *testing.T) {
			pool := testdb.New(t, "../../../migrations")
			id := uuid.New()
			original := seedCart(t, pool, id)
			ctx := context.Background()
			event := "INSERT"
			if table == "cart_items" {
				event = "DELETE"
			}
			_, err := pool.Exec(ctx, `CREATE FUNCTION fail_checkout() RETURNS trigger LANGUAGE plpgsql AS $$
    BEGIN RAISE EXCEPTION 'injected checkout failure'; END $$;
    CREATE TRIGGER fail_checkout BEFORE `+event+` ON `+table+`
    FOR EACH ROW EXECUTE FUNCTION fail_checkout()`)
			require.NoError(t, err)
			result, err := newService(pool, restaurantStub{}).CreateOrder(ctx, id, "123 Test Street")
			require.Error(t, err)
			require.Equal(t, domain.Order{}, result)
			got, err := cart.New(pool).Find(ctx, id)
			require.NoError(t, err)
			require.Equal(t, original, got)
			orders, err := order.New(pool).GetAllOrders(ctx, id)
			require.NoError(t, err)
			require.Empty(t, orders)
		})
	}
}

func TestConcurrentCheckout(t *testing.T) {
	pool := testdb.New(t, "../../../migrations")
	id := uuid.New()
	seedCart(t, pool, id)
	service := newService(pool, restaurantStub{})
	ctx := context.Background()
	start := make(chan struct{})
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Go(func() { <-start; _, err := service.CreateOrder(ctx, id, "123 Test Street"); errs <- err })
	}
	close(start)
	wg.Wait()
	close(errs)
	successes, empty := 0, 0
	for err := range errs {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, domain.ErrCartEmpty):
			empty++
		default:
			require.NoError(t, err)
		}
	}
	require.Equal(t, 1, successes)
	require.Equal(t, 1, empty)
	orders, err := order.New(pool).GetAllOrders(ctx, id)
	require.NoError(t, err)
	require.Len(t, orders, 1)
}

func TestCheckoutConcurrentCartChange(t *testing.T) {
	pool := testdb.New(t, "../../../migrations")
	id := uuid.New()
	original := seedCart(t, pool, id)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	entered, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	unlock := func() { once.Do(func() { close(release) }) }
	defer unlock()
	service := newService(pool, restaurantStub{onFind: func(ctx context.Context) error {
		close(entered)
		select {
		case <-release:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}})
	checkoutErr := make(chan error, 1)
	go func() { _, err := service.CreateOrder(ctx, id, "123 Test Street"); checkoutErr <- err }()
	select {
	case <-entered:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	item := domain.CartItem{MenuItemID: uuid.New(), Quantity: 1}
	updateErr := make(chan error, 1)
	go func() {
		_, err := service.AddItemToCart(ctx, id, original.RestaurantID, item)
		updateErr <- err
	}()
	// Both operations have acquired connections; checkout still holds the cart lock.
	require.Eventually(t, func() bool { return pool.Stat().AcquiredConns() == 2 }, 5*time.Second, 10*time.Millisecond)
	unlock()
	require.NoError(t, <-checkoutErr)
	require.NoError(t, <-updateErr)
	got, err := cart.New(pool).Find(ctx, id)
	require.NoError(t, err)
	require.Len(t, got.Items, 1)
	require.Equal(t, item.MenuItemID, got.Items[0].MenuItemID)
	orders, err := order.New(pool).GetAllOrders(ctx, id)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	require.Len(t, orders[0].Items, 1)
	require.Equal(t, original.Items[0].MenuItemID, orders[0].Items[0].MenuItemID)
}
