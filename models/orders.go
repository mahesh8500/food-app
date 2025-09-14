package models

import "time"

type Order struct {
	ID        string    `db:"id"`
	Items     []byte    `db:"items"`
	Products  []byte    `db:"products"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
