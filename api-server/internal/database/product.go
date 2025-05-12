package database

import "time"

type Product struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Price     float64   `db:"price" json:"price"`
	Available bool      `db:"available" json:"available"`
	Created   time.Time `db:"created" json:"created"`
}
