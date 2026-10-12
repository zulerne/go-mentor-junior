package order

import (
	"context"
	"fmt"
	"time"

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

func (s *Store) db() querier {
	if s.tx != nil {
		return s.tx
	}
	return s.pool
}

type querier interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

const selectOrders = `
 SELECT o.id, o.customer_id, o.restaurant_id, o.status, o.subtotal_minor,
        o.currency, o.delivery_address, o.rejection_reason, o.delivery_status,
        o.created_at, o.updated_at, i.menu_item_id, coalesce(i.name, ''),
        coalesce(i.unit_price_minor, 0), coalesce(i.quantity, 0), coalesce(i.instructions, '')
 FROM orders o LEFT JOIN order_items i ON i.order_id = o.id `

func readOrders(ctx context.Context, db querier, clause string, id uuid.UUID) ([]domain.Order, error) {
	rows, err := db.Query(ctx, selectOrders+clause, id)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()
	var result []domain.Order
	for rows.Next() {
		var value domain.Order
		var delivery pgtype.Text
		var menuID pgtype.UUID
		var item domain.OrderItem
		if err = rows.Scan(
			&value.ID, &value.CustomerID, &value.RestaurantID, &value.Status, &value.SubtotalMinor,
			&value.Currency, &value.DeliveryAddress, &value.RejectionReason, &delivery,
			&value.CreatedAt, &value.UpdatedAt, &menuID, &item.Name,
			&item.UnitPriceMinor, &item.Quantity, &item.Instructions,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		value.DeliveryStatus = domain.OrderDeliveryStatus(delivery.String)
		value.CreatedAt = value.CreatedAt.UTC()
		value.UpdatedAt = value.UpdatedAt.UTC()
		if len(result) == 0 || result[len(result)-1].ID != value.ID {
			result = append(result, value)
		}
		if menuID.Valid {
			item.MenuItemID = uuid.UUID(menuID.Bytes)
			last := &result[len(result)-1]
			last.Items = append(last.Items, item)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("read orders: %w", err)
	}
	return result, nil
}

func (s *Store) Find(ctx context.Context, id uuid.UUID) (domain.Order, error) {
	orders, err := readOrders(ctx, s.db(), "WHERE o.id = $1 ORDER BY i.menu_item_id", id)
	if err != nil {
		return domain.Order{}, err
	}
	if len(orders) == 0 {
		return domain.Order{}, domain.ErrNotFound
	}
	return orders[0], nil
}

func (s *Store) GetAllOrders(ctx context.Context, customerID uuid.UUID) ([]domain.Order, error) {
	return readOrders(ctx, s.db(),
		"WHERE o.customer_id = $1 ORDER BY o.created_at, o.id, i.menu_item_id", customerID)
}

func (s *Store) FindByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]domain.Order, error) {
	return readOrders(ctx, s.db(),
		"WHERE o.restaurant_id = $1 ORDER BY o.created_at, o.id, i.menu_item_id", restaurantID)
}

// Update creates a pending order with immutable item snapshots or advances its
// status under a row lock. Repeating the current status leaves timestamps intact.
func (s *Store) Update(ctx context.Context, value domain.Order) (domain.Order, error) {
	if s.tx != nil {
		return updateOrder(ctx, s.tx, value)
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.Order{}, fmt.Errorf("begin order update: %w", err)
	}
	defer func() { _ = tx.Rollback(context.Background()) }()
	result, err := updateOrder(ctx, tx, value)
	if err != nil {
		return domain.Order{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Order{}, fmt.Errorf("commit order update: %w", err)
	}
	return result, nil
}

func updateOrder(ctx context.Context, tx pgx.Tx, value domain.Order) (domain.Order, error) {
	var err error

	if value.Status == domain.Pending && value.UpdatedAt.IsZero() {
		if err = insertOrder(ctx, tx, value); err != nil {
			return domain.Order{}, err
		}
	}
	orders, err := readOrders(ctx, tx, "WHERE o.id = $1 ORDER BY i.menu_item_id FOR UPDATE OF o", value.ID)
	if err != nil {
		return domain.Order{}, err
	}
	if len(orders) == 0 {
		return domain.Order{}, domain.ErrNotFound
	}
	current := orders[0]
	if current.CustomerID != value.CustomerID || current.RestaurantID != value.RestaurantID {
		return domain.Order{}, domain.ErrOrderAccessDenied
	}
	if current.Status != value.Status {
		if !validTransition(current.Status, value.Status) {
			return domain.Order{}, domain.ErrInvalidOrderTransition
		}
		current.Status = value.Status
		current.RejectionReason = value.RejectionReason
		current.DeliveryStatus = deliveryStatus(value.Status)
		err = tx.QueryRow(ctx, `UPDATE orders
   SET status = $2, rejection_reason = $3, delivery_status = $4, updated_at = clock_timestamp()
   WHERE id = $1 RETURNING updated_at`,
			current.ID, current.Status, current.RejectionReason, nullableDelivery(current.DeliveryStatus),
		).Scan(&current.UpdatedAt)
		if err != nil {
			return domain.Order{}, fmt.Errorf("update order status: %w", err)
		}
		current.UpdatedAt = current.UpdatedAt.UTC()
	}

	return current, nil
}

func insertOrder(ctx context.Context, tx pgx.Tx, value domain.Order) error {
	if len(value.Items) == 0 {
		return domain.ErrCartEmpty
	}
	tag, err := tx.Exec(ctx, `INSERT INTO orders
  (id, customer_id, restaurant_id, status, subtotal_minor, currency, delivery_address,
   rejection_reason, delivery_status, created_at, updated_at)
  VALUES ($1, $2, $3, $4, $5, $6, $7, '', NULL, $8, $9)
  ON CONFLICT (id) DO NOTHING`,
		value.ID, value.CustomerID, value.RestaurantID, domain.Pending, value.SubtotalMinor,
		value.Currency, value.DeliveryAddress, value.CreatedAt, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil
	}
	for _, item := range value.Items {
		_, err = tx.Exec(ctx, `INSERT INTO order_items
   (id, order_id, menu_item_id, name, unit_price_minor, quantity, instructions)
   VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			uuid.New(), value.ID, item.MenuItemID, item.Name, item.UnitPriceMinor, item.Quantity, item.Instructions,
		)
		if err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}
	return nil
}

func validTransition(from, to domain.OrderStatus) bool {
	switch from {
	case domain.Pending:
		return to == domain.Accepted || to == domain.Rejected || to == domain.Cancelled
	case domain.Accepted:
		return to == domain.Preparing
	case domain.Preparing:
		return to == domain.ReadyForPickup
	case domain.ReadyForPickup, domain.Rejected, domain.Cancelled:
		return false
	default:
		return false
	}
}

func deliveryStatus(status domain.OrderStatus) domain.OrderDeliveryStatus {
	switch status {
	case domain.Accepted, domain.Preparing:
		return domain.WaitingForPreparation
	case domain.ReadyForPickup:
		return domain.InDelivery
	case domain.Pending, domain.Rejected, domain.Cancelled:
		return domain.UnspecifiedOrderDeliveryStatus
	default:
		return domain.UnspecifiedOrderDeliveryStatus
	}
}

func nullableDelivery(status domain.OrderDeliveryStatus) pgtype.Text {
	return pgtype.Text{String: string(status), Valid: status != domain.UnspecifiedOrderDeliveryStatus}
}
