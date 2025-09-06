package order

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(d *sql.DB) *Repository {
	return &Repository{db: d}
}

func (r *Repository) CreateOrder(ctx context.Context, o *Order) error {
	const query = `
	INSERT INTO orders (customer, status, created_at) VALUES ($1,$2,$3) RETURNING id`

	o.CreatedAt = time.Now().UTC()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := r.db.QueryRowContext(ctx, query, o.Customer, o.Status, o.CreatedAt).Scan(&o.ID); err != nil {
		return fmt.Errorf("create order: %w", err)
	}
	return nil
}
