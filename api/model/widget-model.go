package model

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Widget is the type for all widgets
type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	IsRecurring    bool      `json:"is_recurring"`
	PlanID         string    `json:"plan_id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

type SqlDB struct {
	*pgxpool.Pool
}

// GetWidget gets one widget by id
func (m *SqlDB) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widget

	row := m.QueryRow(ctx, `
		select 
			id, name, description, inventory_level, price, image, is_recurring, plan_id, created_at, updated_at
		from 
			widgets
		where id = $1`, id)

	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.IsRecurring,
		&widget.PlanID,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	)

	if err != nil {
		return widget, err
	}
	return widget, nil
}
