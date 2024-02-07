package model

import (
	"context"
	"time"
)

// Order is the type for all orders
type Order struct {
	ID            int          `json:"id"`
	WidgetID      int          `json:"widget_id"`
	TransactionID int          `json:"transaction_id"`
	CustomerID    int          `json:"customer_id"`
	StatusID      int          `json:"status_id"`
	Quantity      int          `json:"quantity"`
	Amount        int          `json:"amount"`
	CreatedAt     time.Time    `json:"-"`
	UpdatedAt     time.Time    `json:"-"`
	Widget        Widget       `json:"widget"`
	Transactions  Transactions `json:"transactions"`
	Customer      Customer     `json:"customer"`
}

// InsertOrder inserts a new order, and returns its id
func (m *SqlDB) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	stmt := `
		insert into orders
			(widget_id, transaction_id, status_id, quantity, customer_id,
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

func (m *SqlDB) GetAllOrders() ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var orders []*Order

	query := `
	select
		o.id, o.widget_id, o.transaction_id, o.customer_id, 
		o.status_id, o.quantity, o.amount, o.created_at,
		o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
		t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
		t.bank_return_code, c.id, c.first_name, c.last_name, c.email
		
	from
		orders o
		left join widgets w on (o.widget_id = w.id)
		left join transactions t on (o.transaction_id = t.id)
		left join customers c on (o.customer_id = c.id)
	where
		w.is_recurring = false
	order by
		o.created_at desc
	`

	rows, err := m.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.WidgetID,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.Widget.ID,
			&o.Widget.Name,
			&o.Transactions.ID,
			&o.Transactions.Amount,
			&o.Transactions.Currency,
			&o.Transactions.LastFour,
			&o.Transactions.ExpiryMonth,
			&o.Transactions.ExpiryYear,
			&o.Transactions.PaymentIntent,
			&o.Transactions.BankReturnCode,
			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}

func (m *SqlDB) GetAllOrdersPaginated(limit, page int) ([]*Order, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	offset := (page - 1) * limit

	var orders []*Order

	query := `
	select
		o.id, o.widget_id, o.transaction_id, o.customer_id, 
		o.status_id, o.quantity, o.amount, o.created_at,
		o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
		t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
		t.bank_return_code, c.id, c.first_name, c.last_name, c.email
		
	from
		orders o
		left join widgets w on (o.widget_id = w.id)
		left join transactions t on (o.transaction_id = t.id)
		left join customers c on (o.customer_id = c.id)
	where
		w.is_recurring = false
	order by
		o.created_at desc
	limit $1 offset $2
	`

	rows, err := m.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.WidgetID,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.Widget.ID,
			&o.Widget.Name,
			&o.Transactions.ID,
			&o.Transactions.Amount,
			&o.Transactions.Currency,
			&o.Transactions.LastFour,
			&o.Transactions.ExpiryMonth,
			&o.Transactions.ExpiryYear,
			&o.Transactions.PaymentIntent,
			&o.Transactions.BankReturnCode,
			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		orders = append(orders, &o)
	}

	query = `
		SELECT COUNT(*)
		FROM orders o
		LEFT JOIN widgets w ON (o.widget_id = w.id)
		WHERE w.is_recurring = 0
	`

	var totalRecords int
	err = m.QueryRow(ctx, query).Scan(&totalRecords)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := totalRecords / limit

	return orders, lastPage, totalRecords, nil
}

func (m *SqlDB) GetAllSubscriptions() ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var orders []*Order

	query := `
	select
		o.id, o.widget_id, o.transaction_id, o.customer_id, 
		o.status_id, o.quantity, o.amount, o.created_at,
		o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
		t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
		t.bank_return_code, c.id, c.first_name, c.last_name, c.email
		
	from
		orders o
		left join widgets w on (o.widget_id = w.id)
		left join transactions t on (o.transaction_id = t.id)
		left join customers c on (o.customer_id = c.id)
	where
		w.is_recurring = true
	order by
		o.created_at desc
	`

	rows, err := m.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.WidgetID,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.Widget.ID,
			&o.Widget.Name,
			&o.Transactions.ID,
			&o.Transactions.Amount,
			&o.Transactions.Currency,
			&o.Transactions.LastFour,
			&o.Transactions.ExpiryMonth,
			&o.Transactions.ExpiryYear,
			&o.Transactions.PaymentIntent,
			&o.Transactions.BankReturnCode,
			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}

// GetAllSubscriptionsPaginated returns a slice of a subset of subscriptions
func (m *SqlDB) GetAllSubscriptionsPaginated(pageSize, page int) ([]*Order, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	offset := (page - 1) * pageSize

	var orders []*Order

	query := `
	select
		o.id, o.widget_id, o.transaction_id, o.customer_id, 
		o.status_id, o.quantity, o.amount, o.created_at,
		o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
		t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
		t.bank_return_code, c.id, c.first_name, c.last_name, c.email
		
	from
		orders o
		left join widgets w on (o.widget_id = w.id)
		left join transactions t on (o.transaction_id = t.id)
		left join customers c on (o.customer_id = c.id)
	where
		w.is_recurring = 1
	order by
		o.created_at desc
	limit $1 offset $2
	`

	rows, err := m.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(
			&o.ID,
			&o.WidgetID,
			&o.TransactionID,
			&o.CustomerID,
			&o.StatusID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.Widget.ID,
			&o.Widget.Name,
			&o.Transactions.ID,
			&o.Transactions.Amount,
			&o.Transactions.Currency,
			&o.Transactions.LastFour,
			&o.Transactions.ExpiryMonth,
			&o.Transactions.ExpiryYear,
			&o.Transactions.PaymentIntent,
			&o.Transactions.BankReturnCode,
			&o.Customer.ID,
			&o.Customer.FirstName,
			&o.Customer.LastName,
			&o.Customer.Email,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		orders = append(orders, &o)
	}

	query = `
		select 
			count(o.id)
		from 
			orders o
			left join widgets w on (o.widget_id = w.id)
		where
			w.is_recurring = 1
	`
	var totalRecords int
	countRow := m.QueryRow(ctx, query)
	err = countRow.Scan(&totalRecords)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := totalRecords / pageSize

	return orders, lastPage, totalRecords, nil
}

func (m *SqlDB) GetOrderByID(id int) (Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var o Order

	query := `
	select
		o.id, o.widget_id, o.transaction_id, o.customer_id, 
		o.status_id, o.quantity, o.amount, o.created_at,
		o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
		t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
		t.bank_return_code, c.id, c.first_name, c.last_name, c.email
		
	from
		orders o
		left join widgets w on (o.widget_id = w.id)
		left join transactions t on (o.transaction_id = t.id)
		left join customers c on (o.customer_id = c.id)
	where
		o.id = $1
	`

	row := m.QueryRow(ctx, query, id)

	err := row.Scan(
		&o.ID,
		&o.WidgetID,
		&o.TransactionID,
		&o.CustomerID,
		&o.StatusID,
		&o.Quantity,
		&o.Amount,
		&o.CreatedAt,
		&o.UpdatedAt,
		&o.Widget.ID,
		&o.Widget.Name,
		&o.Transactions.ID,
		&o.Transactions.Amount,
		&o.Transactions.Currency,
		&o.Transactions.LastFour,
		&o.Transactions.ExpiryMonth,
		&o.Transactions.ExpiryYear,
		&o.Transactions.PaymentIntent,
		&o.Transactions.BankReturnCode,
		&o.Customer.ID,
		&o.Customer.FirstName,
		&o.Customer.LastName,
		&o.Customer.Email,
	)
	if err != nil {
		return o, err
	}

	return o, nil
}

func (m *SqlDB) UpdateOrderStatus(id, statusID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE orders SET status_id = $1 WHERE id = $2`

	_, err := m.Exec(ctx, stmt, statusID, id)
	if err != nil {
		return err
	}

	return nil
}
