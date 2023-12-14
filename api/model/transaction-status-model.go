package model

import "time"

// TransactionStatus is the type for transaction statuses
type TransactionStatuses struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
