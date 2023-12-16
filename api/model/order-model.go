package model

import (
	"context"
	"time"
)

// Order is the type for all orders
type Order struct {
	ID            int       `json:"id"`
	WidgetID      int       `json:"widget_id"`
	TransactionID int       `json:"transaction_id"`
	CustomerID    int       `json:"customers_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

// InsertOrder inserts a new order, and returns its id
func (m *SqlDB) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	stmt := `
		insert into orders
			(widget_id, transaction_id, status_id, quantity, customers_id,
			amount, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;`

	err := m.QueryRow(ctx, stmt,
		order.WidgetID,
		order.TransactionID,
		order.StatusID,
		order.Quantity,
		order.CustomerID,
		order.Amount,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
