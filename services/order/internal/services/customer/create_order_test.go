package customer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestCreateOrder_EmptyCart(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID:  testRestaurantID,
		Items:         []domain.CartItem{},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	// rest.EXPECT().Find(mock.Anything, testRestaurantID).Return(domain.Restaurant{
	// 	ID:                testRestaurantID,
	// 	Name:              "restaurant",
	// 	AcceptingOrders:   false,
	// 	MinimumOrderMinor: 100,
	// 	Currency:          "USD",
	// }, nil)

	// rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
	// 	ID:         testMenuItemID,
	// 	Name:       "item",
	// 	Available:  true,
	// 	PriceMinor: 100,
	// 	Currency:   "USD",
	// }, nil)

	// oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	// cStore.EXPECT().Update(mock.Anything, mock.Anything, domain.Cart{}).Return(nil)

	_, err := service.CreateOrder(context.Background(), testCustomerID, "address")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrCartEmpty, err)
}

func TestCreateOrder_RestaurantNotFound(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID: testRestaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "item",
				UnitPriceMinor: 100,
				Currency:       "USD",
				Quantity:       1,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	rest.EXPECT().Find(mock.Anything, testRestaurantID).Return(domain.Restaurant{}, domain.ErrRestaurantNotFound)

	// rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
	// 	ID:         testMenuItemID,
	// 	Name:       "item",
	// 	Available:  true,
	// 	PriceMinor: 100,
	// 	Currency:   "USD",
	// }, nil)

	// oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	// cStore.EXPECT().Update(mock.Anything, mock.Anything, domain.Cart{}).Return(nil)

	_, err := service.CreateOrder(context.Background(), testCustomerID, "address")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrRestaurantNotFound, err)
}

func TestCreateOrder_RestaurantNotAcceptingOrders(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID: testRestaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "item",
				UnitPriceMinor: 100,
				Currency:       "USD",
				Quantity:       1,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	rest.EXPECT().Find(mock.Anything, testRestaurantID).Return(domain.Restaurant{
		ID:                testRestaurantID,
		Name:              "restaurant",
		AcceptingOrders:   false,
		MinimumOrderMinor: 100,
		Currency:          "USD",
	}, nil)

	// rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
	// 	ID:         testMenuItemID,
	// 	Name:       "item",
	// 	Available:  true,
	// 	PriceMinor: 100,
	// 	Currency:   "USD",
	// }, nil)

	// oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	// cStore.EXPECT().Update(mock.Anything, mock.Anything, domain.Cart{}).Return(nil)

	_, err := service.CreateOrder(context.Background(), testCustomerID, "address")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrRestaurantNotAcceptingOrders, err)
}

func TestCreateOrder_ItemNotAvailable(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID: testRestaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "item",
				UnitPriceMinor: 100,
				Currency:       "USD",
				Quantity:       1,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	rest.EXPECT().Find(mock.Anything, testRestaurantID).Return(domain.Restaurant{
		ID:                testRestaurantID,
		Name:              "restaurant",
		AcceptingOrders:   true,
		MinimumOrderMinor: 100,
		Currency:          "USD",
	}, nil)

	rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
		ID:         testMenuItemID,
		Name:       "item",
		Available:  false,
		PriceMinor: 100,
		Currency:   "USD",
	}, nil)

	// oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	// cStore.EXPECT().Update(mock.Anything, mock.Anything, domain.Cart{}).Return(nil)

	_, err := service.CreateOrder(context.Background(), testCustomerID, "address")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMenuItemNotAvailable, err)
}

func TestCreateOrder_MinimumOrderNotReached(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID: testRestaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "item",
				UnitPriceMinor: 100,
				Currency:       "USD",
				Quantity:       1,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	rest.EXPECT().Find(mock.Anything, testRestaurantID).Return(domain.Restaurant{
		ID:                testRestaurantID,
		Name:              "restaurant",
		AcceptingOrders:   true,
		MinimumOrderMinor: 200,
		Currency:          "USD",
	}, nil)

	rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
		ID:         testMenuItemID,
		Name:       "item",
		Available:  true,
		PriceMinor: 100,
		Currency:   "USD",
	}, nil)

	// oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	// cStore.EXPECT().Update(mock.Anything, mock.Anything, domain.Cart{}).Return(nil)

	_, err := service.CreateOrder(context.Background(), testCustomerID, "address")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMinOrderNotReached, err)
}

func TestCreateOrder_ClearsCartOnSuccess(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID: testRestaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "item",
				UnitPriceMinor: 100,
				Currency:       "USD",
				Quantity:       1,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	rest.EXPECT().Find(mock.Anything, testRestaurantID).Return(domain.Restaurant{
		ID:                testRestaurantID,
		Name:              "restaurant",
		AcceptingOrders:   true,
		MinimumOrderMinor: 0,
		Currency:          "USD",
	}, nil)

	rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
		ID:         testMenuItemID,
		Name:       "item",
		Available:  true,
		PriceMinor: 100,
		Currency:   "USD",
	}, nil)

	oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	cStore.EXPECT().Update(mock.Anything, mock.Anything, domain.Cart{}).Return(nil)

	order, err := service.CreateOrder(context.Background(), testCustomerID, "address")

	assert.NoError(t, err)
	assert.Equal(t, domain.Pending, order.Status)
}

func TestCreateOrder_SnapshotsPrices(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID: testRestaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "item",
				UnitPriceMinor: 100,
				Currency:       "USD",
				Quantity:       1,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	rest.EXPECT().Find(mock.Anything, testRestaurantID).Return(domain.Restaurant{
		ID:                testRestaurantID,
		Name:              "restaurant",
		AcceptingOrders:   true,
		MinimumOrderMinor: 0,
		Currency:          "USD",
	}, nil)

	rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
		ID:         testMenuItemID,
		Name:       "item",
		Available:  true,
		PriceMinor: 200,
		Currency:   "USD",
	}, nil)

	oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	cStore.EXPECT().Update(mock.Anything, mock.Anything, mock.Anything).Return(nil)

	order, err := service.CreateOrder(context.Background(), testCustomerID, "address")

	assert.NoError(t, err)
	assert.Equal(t, domain.Pending, order.Status)
	assert.Equal(t, testCustomerID, order.CustomerID)
	assert.Equal(t, testRestaurantID, order.RestaurantID)
	assert.Equal(t, 1, len(order.Items))
	assert.Equal(t, int64(200), order.SubtotalMinor)
	assert.Equal(t, testMenuItemID, order.Items[0].MenuItemID)
	assert.Equal(t, "item", order.Items[0].Name)
}

func TestCreateOrder_Success(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID: testRestaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "item",
				UnitPriceMinor: 100,
				Currency:       "USD",
				Quantity:       1,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	rest.EXPECT().Find(mock.Anything, testRestaurantID).Return(domain.Restaurant{
		ID:                testRestaurantID,
		Name:              "restaurant",
		AcceptingOrders:   true,
		MinimumOrderMinor: 0,
		Currency:          "USD",
	}, nil)

	rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
		ID:         testMenuItemID,
		Name:       "item",
		Available:  true,
		PriceMinor: 100,
		Currency:   "USD",
	}, nil)

	oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	cStore.EXPECT().Update(mock.Anything, mock.Anything, mock.Anything).Return(nil)

	order, err := service.CreateOrder(context.Background(), testCustomerID, "address")

	assert.NoError(t, err)
	assert.Equal(t, domain.Pending, order.Status)
	assert.Equal(t, testCustomerID, order.CustomerID)
	assert.Equal(t, testRestaurantID, order.RestaurantID)
	assert.Equal(t, 1, len(order.Items))
	assert.Equal(t, int64(100), order.SubtotalMinor)
	assert.Equal(t, testMenuItemID, order.Items[0].MenuItemID)
	assert.Equal(t, "item", order.Items[0].Name)
}
