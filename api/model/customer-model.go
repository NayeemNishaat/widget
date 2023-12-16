package model

import (
	"context"
	"time"
)

// Customer is the type for customers
type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// InsertOrder inserts a new order, and returns its id
func (m *SqlDB) InsertCustomer(c Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	stmt := `
		insert into customers
			(first_name, last_name, email, created_at, updated_at)
		values ($1, $2, $3, $4, $5)
		RETURNING id;`

	err := m.QueryRow(ctx, stmt,
		c.FirstName,
		c.LastName,
		c.Email,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
