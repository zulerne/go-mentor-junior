//go:build integration

package cart_test

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store/cart"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store/testdb"
)

func sampleCart() domain.Cart {
	return domain.Cart{
		RestaurantID: uuid.New(), Currency: "USD", SubtotalMinor: 250,
		Items: []domain.CartItem{{MenuItemID: uuid.New(), Name: "Soup", UnitPriceMinor: 125,
			Currency: "USD", Quantity: 2, Instructions: "No salt"}},
	}
}

func TestCartPersistence(t *testing.T) {
	pool := testdb.New(t, "../../../../migrations")
	store := cart.New(pool)
	ctx := context.Background()
	customerID := uuid.New()
	_, err := store.Find(ctx, customerID)
	require.ErrorIs(t, err, domain.ErrNotFound)

	original := sampleCart()
	require.NoError(t, store.Update(ctx, customerID, original))
	got, err := store.Find(ctx, customerID)
	require.NoError(t, err)
	require.Equal(t, original, got)

	// A new store using a separate pool reads the persisted data.
	reopened, err := pgxpool.NewWithConfig(ctx, pool.Config())
	require.NoError(t, err)
	defer reopened.Close()
	got, err = cart.New(reopened).Find(ctx, customerID)
	require.NoError(t, err)
	require.Equal(t, original, got)

	replacement := sampleCart()
	require.NoError(t, store.Update(ctx, customerID, replacement))
	got, err = store.Find(ctx, customerID)
	require.NoError(t, err)
	require.Equal(t, replacement, got)

	// The first item succeeds, the second fails: the entire replacement rolls back.
	invalid := sampleCart()
	bad := invalid.Items[0]
	bad.MenuItemID = uuid.New()
	bad.Quantity = 11
	invalid.Items = append(invalid.Items, bad)
	require.Error(t, store.Update(ctx, customerID, invalid))
	got, err = store.Find(ctx, customerID)
	require.NoError(t, err)
	require.Equal(t, replacement, got)

	absentID := uuid.New()
	require.Error(t, store.Update(ctx, absentID, invalid))
	_, err = store.Find(ctx, absentID)
	require.ErrorIs(t, err, domain.ErrNotFound)

	require.NoError(t, store.Update(ctx, customerID, domain.Cart{}))
	got, err = store.Find(ctx, customerID)
	require.NoError(t, err)
	require.Equal(t, domain.Cart{}, got)
	var empty bool
	require.NoError(t, pool.QueryRow(ctx, `SELECT restaurant_id IS NULL AND currency IS NULL
  AND subtotal_minor = 0 FROM carts WHERE customer_id = $1`, customerID).Scan(&empty))
	require.True(t, empty)

	cancelledCtx, cancel := context.WithCancel(ctx)
	cancel()
	require.Error(t, store.Update(cancelledCtx, customerID, original))
	_, err = store.Find(cancelledCtx, customerID)
	require.Error(t, err)
}

func TestConcurrentCartReplacement(t *testing.T) {
	store := cart.New(testdb.New(t, "../../../../migrations"))
	ctx := context.Background()
	customerID := uuid.New()
	first, second := sampleCart(), sampleCart()
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	start := make(chan struct{})
	for _, value := range []domain.Cart{first, second} {
		wg.Go(func() { <-start; errs <- store.Update(ctx, customerID, value) })
	}
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
	got, err := store.Find(ctx, customerID)
	require.NoError(t, err)
	require.True(t, got.RestaurantID == first.RestaurantID || got.RestaurantID == second.RestaurantID)
	if got.RestaurantID == first.RestaurantID {
		require.Equal(t, first, got)
	} else {
		require.Equal(t, second, got)
	}
}
