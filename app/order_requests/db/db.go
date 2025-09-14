package db

import (
	"context"
	"fmt"
	"foodApp/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgOrderRequestRepo struct {
	p *pgxpool.Pool
}

func NewOrderRequestRepository(pool *pgxpool.Pool) *pgOrderRequestRepo {
	return &pgOrderRequestRepo{p: pool}
}

// Insert saves a raw order request into order_requests table
func (r *pgOrderRequestRepo) Insert(ctx context.Context, or models.OrderRequest) (models.OrderRequest, error) {
	var id string
	err := r.p.QueryRow(
		ctx,
		`INSERT INTO order_requests (items, coupon_code)
         VALUES ($1::jsonb, $2)
         RETURNING id, created_at`,
		or.Items, or.CouponCode,
	).Scan(&id, &or.CreatedAt)

	if err != nil {
		return models.OrderRequest{}, fmt.Errorf("insert order_request: %w", err)
	}

	or.ID = id
	return or, nil
}
