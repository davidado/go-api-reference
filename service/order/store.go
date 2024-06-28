// Package order : Order service
package order

import (
	"database/sql"

	"github.com/davidado/go-api-reference/types"
)

// Store : Order store
type Store struct {
	db *sql.DB
}

// NewStore creates a new order store
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// CreateOrder creates a new order
func (s *Store) CreateOrder(o types.Order) (int, error) {
	res, err := s.db.Exec("INSERT INTO orders (user_id, total, status, address) VALUES (?, ?, ?, ?)", o.UserID, o.Total, o.Status, o.Address)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// CreateOrderItem creates a new order item
func (s *Store) CreateOrderItem(oi types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)", oi.OrderID, oi.ProductID, oi.Quantity, oi.Price)
	return err
}
