package customer_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestAddItemToCart_NewCart(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{}, domain.ErrNotFound)
	rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
		ID:          testMenuItemID,
		Name:        "Test Item",
		Description: "Test Description",
		PriceMinor:  100,
		Currency:    "USD",
		Available:   true,
	}, nil)
	cStore.EXPECT().Update(mock.Anything, testCustomerID, mock.Anything).Return(nil)

	cart, err := service.AddItemToCart(context.Background(), testCustomerID, testRestaurantID, domain.CartItem{
		MenuItemID:   testMenuItemID,
		Quantity:     1,
		Instructions: "",
	})

	require.NoError(t, err)
	assert.NotNil(t, cart)
	assert.Equal(t, "Test Item", cart.Items[0].Name)
	assert.Equal(t, int64(100), cart.Items[0].UnitPriceMinor)
}

func TestAddItemToCart_ExistingCartSameRestaurant(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID:  testRestaurantID,
		Items:         []domain.CartItem{},
		SubtotalMinor: 0,
		Currency:      "USD",
	}, nil)
	rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
		ID:          testMenuItemID,
		Name:        "Test Item",
		Description: "Test Description",
		PriceMinor:  100,
		Currency:    "USD",
		Available:   true,
	}, nil)
	cStore.EXPECT().Update(mock.Anything, testCustomerID, mock.Anything).Return(nil)

	cart, err := service.AddItemToCart(context.Background(), testCustomerID, testRestaurantID, domain.CartItem{
		MenuItemID:   testMenuItemID,
		Quantity:     1,
		Instructions: "",
	})

	require.NoError(t, err)
	assert.NotNil(t, cart)
	assert.Equal(t, "Test Item", cart.Items[0].Name)
	assert.Equal(t, int64(100), cart.Items[0].UnitPriceMinor)
}

func TestAddItemToCart_ExistingCartDifferentRestaurant(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID:  testRestaurantID,
		Items:         []domain.CartItem{},
		SubtotalMinor: 0,
		Currency:      "USD",
	}, nil)

	_, err := service.AddItemToCart(context.Background(), testCustomerID, uuid.New(), domain.CartItem{
		MenuItemID:   testMenuItemID,
		Quantity:     1,
		Instructions: "",
	})

	require.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestAddItemToCart_UpdatesQuantityIfItemExists(t *testing.T) {
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
				Name:           "Test Item",
				UnitPriceMinor: 15,
				Currency:       "USD",
				Quantity:       10,
				Instructions:   "",
			},
		},
		SubtotalMinor: 0,
		Currency:      "USD",
	}, nil)
	rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
		ID:          testMenuItemID,
		Name:        "Test Item",
		Description: "Test Description",
		PriceMinor:  20,
		Currency:    "USD",
		Available:   true,
	}, nil)
	cStore.EXPECT().Update(mock.Anything, testCustomerID, mock.Anything).Return(nil)

	cart, err := service.AddItemToCart(context.Background(), testCustomerID, testRestaurantID, domain.CartItem{
		MenuItemID:   testMenuItemID,
		Quantity:     1,
		Instructions: "",
	})

	require.NoError(t, err)
	assert.Equal(t, "Test Item", cart.Items[0].Name)
	assert.Len(t, cart.Items, 1)
	assert.Equal(t, int32(11), cart.Items[0].Quantity)
}

func TestAddItemToCart_CartUnitLimit(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID:  testRestaurantID,
		Items:         []domain.CartItem{},
		SubtotalMinor: 0,
		Currency:      "USD",
	}, nil)

	_, err := service.AddItemToCart(context.Background(), testCustomerID, testRestaurantID, domain.CartItem{
		MenuItemID:   testMenuItemID,
		Quantity:     51,
		Instructions: "",
	})

	require.Error(t, err)
	assert.Equal(t, domain.ErrCartLimit, err)
}

func TestAddItemToCart_CartLineLimit(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID:  testRestaurantID,
		Items:         make([]domain.CartItem, 20),
		SubtotalMinor: 1500,
		Currency:      "USD",
	}, nil)

	_, err := service.AddItemToCart(context.Background(), testCustomerID, testRestaurantID, domain.CartItem{
		MenuItemID:   testMenuItemID,
		Quantity:     6,
		Instructions: "",
	})

	require.Error(t, err)
	assert.Equal(t, domain.ErrCartLimit, err)
}

func TestAddItemToCart_ItemNotAvailable(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID:  testRestaurantID,
		Items:         []domain.CartItem{},
		SubtotalMinor: 1500,
		Currency:      "USD",
	}, nil)
	rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
		ID:          testMenuItemID,
		Name:        "Test Item",
		Description: "Test Description",
		PriceMinor:  100,
		Currency:    "USD",
		Available:   false,
	}, nil)

	_, err := service.AddItemToCart(context.Background(), testCustomerID, testRestaurantID, domain.CartItem{
		MenuItemID:   testMenuItemID,
		Quantity:     6,
		Instructions: "",
	})

	require.Error(t, err)
	assert.Equal(t, domain.ErrMenuItemNotAvailable, err)
}

func TestAddItemToCart_SubtotalCalculation(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID:  testRestaurantID,
		Items:         make([]domain.CartItem, 5),
		SubtotalMinor: 1500,
		Currency:      "USD",
	}, nil)
	rest.EXPECT().FindItem(mock.Anything, testRestaurantID, testMenuItemID).Return(domain.MenuItem{
		ID:          testMenuItemID,
		Name:        "Test Item",
		Description: "Test Description",
		PriceMinor:  100,
		Currency:    "USD",
		Available:   true,
	}, nil)
	cStore.EXPECT().Update(mock.Anything, testCustomerID, mock.Anything).Return(nil)

	cart, err := service.AddItemToCart(context.Background(), testCustomerID, testRestaurantID, domain.CartItem{
		MenuItemID:   testMenuItemID,
		Quantity:     6,
		Instructions: "",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1500+6*100), cart.SubtotalMinor)
	assert.Len(t, cart.Items, 6)
}
