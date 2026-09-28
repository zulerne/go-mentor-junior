package customer_test

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/services/customer"
)

var (
	testCustomerID   = uuid.MustParse("ef55f77a-7738-426f-93a4-f78f2baf7970")
	testRestaurantID = uuid.MustParse("fe70cb38-d10d-452c-8860-1af5936f7037")
	testOrderID      = uuid.MustParse("2f5efdc6-c936-4fdf-a681-1e21e73e6e71")
	testMenuItemID   = uuid.MustParse("58185771-1f9f-4d71-9609-bdacea1deb2e")
	fixedNow         = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	fixedOrderID     = uuid.MustParse("aaaaaaaa-0000-0000-0000-000000000001")
)

func newCustomer(
	orderStore customer.OrderStore,
	cartStore customer.CartStore,
	rest customer.RestaurantProvider,
) *customer.Service {
	log := slog.New(slog.DiscardHandler)
	return customer.New(orderStore, cartStore, rest, log,
		customer.WithNow(func() time.Time { return fixedNow }),
		customer.WithNewID(func() uuid.UUID { return fixedOrderID }),
	)
}
