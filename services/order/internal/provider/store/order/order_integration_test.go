//go:build integration

package order_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store/order"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store/testdb"
)

func sampleOrder() domain.Order {
	return domain.Order{
		ID: uuid.New(), CustomerID: uuid.New(), RestaurantID: uuid.New(), Status: domain.Pending,
		SubtotalMinor: 325, Currency: "USD", DeliveryAddress: "123 Test Street",
		CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
		Items: []domain.OrderItem{
			{MenuItemID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Soup",
				UnitPriceMinor: 125, Quantity: 2, Instructions: "No salt"},
			{MenuItemID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Bread",
				UnitPriceMinor: 75, Quantity: 1},
		},
	}
}

func TestOrderPersistence(t *testing.T) {
	pool := testdb.New(t, "../../../../migrations")
	store := order.New(pool)
	ctx := context.Background()
	_, err := store.Find(ctx, uuid.New())
	require.ErrorIs(t, err, domain.ErrNotFound)
	value := sampleOrder()
	saved, err := store.Update(ctx, value)
	require.NoError(t, err)
	require.Equal(t, value.Items, saved.Items)
	require.Equal(t, value.CreatedAt, saved.CreatedAt)
	require.False(t, saved.UpdatedAt.IsZero())
	require.Empty(t, saved.DeliveryStatus)
	got, err := store.Find(ctx, value.ID)
	require.NoError(t, err)
	require.Equal(t, saved, got)

	reopened, err := pgxpool.NewWithConfig(ctx, pool.Config())
	require.NoError(t, err)
	defer reopened.Close()
	got, err = order.New(reopened).Find(ctx, value.ID)
	require.NoError(t, err)
	require.Equal(t, saved, got)

	repeated, err := store.Update(ctx, value)
	require.NoError(t, err)
	require.Equal(t, saved, repeated)
	var noDelivery bool
	require.NoError(t, pool.QueryRow(ctx, "SELECT delivery_status IS NULL FROM orders WHERE id = $1", value.ID).Scan(&noDelivery))
	require.True(t, noDelivery)

	next := sampleOrder()
	next.CustomerID, next.RestaurantID = value.CustomerID, value.RestaurantID
	next.CreatedAt = value.CreatedAt.Add(time.Second)
	savedNext, err := store.Update(ctx, next)
	require.NoError(t, err)
	other := sampleOrder()
	_, err = store.Update(ctx, other)
	require.NoError(t, err)
	for _, tc := range []struct {
		id   uuid.UUID
		list func(context.Context, uuid.UUID) ([]domain.Order, error)
	}{
		{value.CustomerID, store.GetAllOrders},
		{value.RestaurantID, store.FindByRestaurant},
	} {
		orders, listErr := tc.list(ctx, tc.id)
		require.NoError(t, listErr)
		require.Equal(t, []domain.Order{saved, savedNext}, orders)
		orders, listErr = tc.list(ctx, uuid.New())
		require.NoError(t, listErr)
		require.Empty(t, orders)
	}

	ctx, cancel := context.WithCancel(ctx)
	cancel()
	_, err = store.Find(ctx, value.ID)
	require.Error(t, err)
	_, err = store.Update(ctx, saved)
	require.Error(t, err)
}

func TestOrderInsertRollback(t *testing.T) {
	store := order.New(testdb.New(t, "../../../../migrations"))
	ctx := context.Background()
	for _, corrupt := range []func(*domain.Order){
		func(value *domain.Order) { value.Items[1].Quantity = 11 },
		func(value *domain.Order) { value.Items[1].MenuItemID = value.Items[0].MenuItemID },
	} {
		value := sampleOrder()
		corrupt(&value)
		_, err := store.Update(ctx, value)
		require.Error(t, err)
		_, err = store.Find(ctx, value.ID)
		require.ErrorIs(t, err, domain.ErrNotFound)
	}
	value := sampleOrder()
	value.Items = nil
	_, err := store.Update(ctx, value)
	require.ErrorIs(t, err, domain.ErrCartEmpty)
}

func TestOrderTransitions(t *testing.T) {
	store := order.New(testdb.New(t, "../../../../migrations"))
	ctx := context.Background()
	value, err := store.Update(ctx, sampleOrder())
	require.NoError(t, err)
	original := value
	for _, status := range []domain.OrderStatus{domain.Accepted, domain.Preparing, domain.ReadyForPickup} {
		value.Status = status
		// Status changes must preserve the checkout snapshot even with stale input.
		value.Currency = "EUR"
		value.SubtotalMinor = 1
		value.Items = nil
		value, err = store.Update(ctx, value)
		require.NoError(t, err)
		require.Equal(t, status, value.Status)
		require.Equal(t, original.Items, value.Items)
		require.Equal(t, original.SubtotalMinor, value.SubtotalMinor)
		require.Equal(t, original.Currency, value.Currency)
		require.Equal(t, original.CreatedAt, value.CreatedAt)
		expectedDelivery := domain.WaitingForPreparation
		if status == domain.ReadyForPickup {
			expectedDelivery = domain.InDelivery
		}
		require.Equal(t, expectedDelivery, value.DeliveryStatus)
		repeated, repeatErr := store.Update(ctx, value)
		require.NoError(t, repeatErr)
		require.Equal(t, value, repeated)
	}
	value.Status = domain.Pending
	_, err = store.Update(ctx, value)
	require.ErrorIs(t, err, domain.ErrInvalidOrderTransition)

	for _, status := range []domain.OrderStatus{domain.Cancelled, domain.Rejected} {
		pending, createErr := store.Update(ctx, sampleOrder())
		require.NoError(t, createErr)
		pending.Status = status
		if status == domain.Rejected {
			pending.RejectionReason = "Kitchen closed"
		}
		saved, updateErr := store.Update(ctx, pending)
		require.NoError(t, updateErr)
		require.Equal(t, status, saved.Status)
		require.Equal(t, pending.RejectionReason, saved.RejectionReason)
		require.Empty(t, saved.DeliveryStatus)
		repeated, repeatErr := store.Update(ctx, pending)
		require.NoError(t, repeatErr)
		require.Equal(t, saved, repeated)
	}
	pending, err := store.Update(ctx, sampleOrder())
	require.NoError(t, err)
	pending.Status = domain.ReadyForPickup
	_, err = store.Update(ctx, pending)
	require.ErrorIs(t, err, domain.ErrInvalidOrderTransition)
	pending.Status = domain.Accepted
	pending.CustomerID = uuid.New()
	_, err = store.Update(ctx, pending)
	require.ErrorIs(t, err, domain.ErrOrderAccessDenied)
	missing := sampleOrder()
	missing.Status = domain.Accepted
	_, err = store.Update(ctx, missing)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestConcurrentOrderTransitions(t *testing.T) {
	store := order.New(testdb.New(t, "../../../../migrations"))
	ctx := context.Background()
	pending, err := store.Update(ctx, sampleOrder())
	require.NoError(t, err)
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	start := make(chan struct{})
	for _, status := range []domain.OrderStatus{domain.Accepted, domain.Cancelled} {
		wg.Go(func() {
			value := pending
			value.Status = status
			<-start
			_, updateErr := store.Update(ctx, value)
			errs <- updateErr
		})
	}
	close(start)
	wg.Wait()
	close(errs)
	succeeded, conflicted := 0, 0
	for updateErr := range errs {
		switch {
		case updateErr == nil:
			succeeded++
		case errors.Is(updateErr, domain.ErrInvalidOrderTransition):
			conflicted++
		default:
			require.NoError(t, updateErr)
		}
	}
	require.Equal(t, 1, succeeded)
	require.Equal(t, 1, conflicted)
	got, err := store.Find(ctx, pending.ID)
	require.NoError(t, err)
	require.True(t, got.Status == domain.Accepted || got.Status == domain.Cancelled)
	require.Equal(t, pending.Items, got.Items)
}
