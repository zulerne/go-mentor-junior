package cart

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type Store struct {
	pool *pgxpool.Pool
	tx   pgx.Tx
}

func New(pool *pgxpool.Pool) *Store { return &Store{pool: pool} }

// NewWithTx binds the store to a transaction owned by the caller.
func NewWithTx(tx pgx.Tx) *Store { return &Store{tx: tx} }

type querier interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

// Find reads the cart and its items from a single database snapshot.
func (s *Store) Find(ctx context.Context, customerID uuid.UUID) (domain.Cart, error) {
	var db querier = s.pool
	if s.tx != nil {
		db = s.tx
	}
	rows, err := db.Query(ctx, `
 SELECT c.restaurant_id, c.subtotal_minor, c.currency,
        i.menu_item_id, coalesce(i.name, ''), coalesce(i.unit_price_minor, 0),
        coalesce(i.currency, ''), coalesce(i.quantity, 0), coalesce(i.instructions, '')
 FROM carts c LEFT JOIN cart_items i ON i.customer_id = c.customer_id
 WHERE c.customer_id = $1 ORDER BY i.menu_item_id`, customerID)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("query cart: %w", err)
	}
	defer rows.Close()
	var result domain.Cart
	found := false
	for rows.Next() {
		var restaurantID, menuItemID pgtype.UUID
		var currency pgtype.Text
		var item domain.CartItem
		if err = rows.Scan(
			&restaurantID, &result.SubtotalMinor, &currency,
			&menuItemID, &item.Name, &item.UnitPriceMinor,
			&item.Currency, &item.Quantity, &item.Instructions,
		); err != nil {
			return domain.Cart{}, fmt.Errorf("scan cart: %w", err)
		}
		found = true
		result.RestaurantID = uuid.UUID(restaurantID.Bytes)
		result.Currency = currency.String
		if menuItemID.Valid {
			item.MenuItemID = uuid.UUID(menuItemID.Bytes)
			result.Items = append(result.Items, item)
		}
	}
	if err = rows.Err(); err != nil {
		return domain.Cart{}, fmt.Errorf("read cart: %w", err)
	}
	if !found {
		return domain.Cart{}, domain.ErrNotFound
	}
	return result, nil
}

// Update atomically replaces a cart and its items. The parent upsert serializes
// replacements for the same customer, including creation of a missing cart.
func (s *Store) Update(ctx context.Context, customerID uuid.UUID, cart domain.Cart) error {
	if s.tx != nil {
		return updateCart(ctx, s.tx, customerID, cart)
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin cart update: %w", err)
	}
	defer func() { _ = tx.Rollback(context.Background()) }()
	if err = updateCart(ctx, tx, customerID, cart); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit cart update: %w", err)
	}
	return nil
}

func updateCart(ctx context.Context, tx pgx.Tx, customerID uuid.UUID, cart domain.Cart) error {
	var restaurantID pgtype.UUID
	var currency pgtype.Text
	if len(cart.Items) == 0 {
		cart.SubtotalMinor = 0
	} else {
		restaurantID = pgtype.UUID{Bytes: cart.RestaurantID, Valid: true}
		currency = pgtype.Text{String: cart.Currency, Valid: true}
	}
	_, err := tx.Exec(ctx, `
 INSERT INTO carts (customer_id, restaurant_id, currency, subtotal_minor)
 VALUES ($1, $2, $3, $4)
 ON CONFLICT (customer_id) DO UPDATE SET
 restaurant_id = EXCLUDED.restaurant_id, currency = EXCLUDED.currency,
 subtotal_minor = EXCLUDED.subtotal_minor`, customerID, restaurantID, currency, cart.SubtotalMinor)
	if err != nil {
		return fmt.Errorf("save cart: %w", err)
	}
	if _, err = tx.Exec(ctx, "DELETE FROM cart_items WHERE customer_id = $1", customerID); err != nil {
		return fmt.Errorf("clear cart items: %w", err)
	}
	for _, item := range cart.Items {
		_, err = tx.Exec(
			ctx,
			`
  INSERT INTO cart_items (id, customer_id, menu_item_id, name, unit_price_minor, currency, quantity, instructions)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			uuid.New(), customerID, item.MenuItemID, item.Name,
			item.UnitPriceMinor, item.Currency, item.Quantity, item.Instructions,
		)
		if err != nil {
			return fmt.Errorf("save cart item: %w", err)
		}
	}

	return nil
}
