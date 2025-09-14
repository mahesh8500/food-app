package db

import (
	"context"
	"fmt"

	_interface "foodApp/app/order/interface"
	"foodApp/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgOrderRepo struct {
	p *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) _interface.Orders {
	return &pgOrderRepo{p: pool}
}

func (r *pgOrderRepo) Insert(ctx context.Context, order models.Order) (models.Order, error) {
	var id string
	err := r.p.QueryRow(
		ctx,
		`INSERT INTO orders (items, products)
     			VALUES ($1::jsonb, $2::jsonb)
     			RETURNING id`,
		order.Items, order.Products,
	).Scan(&id)

	if err != nil {
		return models.Order{}, fmt.Errorf("insert order: %w", err)
	}

	order.ID = id
	return order, nil
}
