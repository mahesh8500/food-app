package db

import (
	"context"
	"fmt"
	"time"

	_interface "foodApp/app/product/interface"
	"foodApp/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = Err("not found")

type Err string

type pgRepo struct {
	p *pgxpool.Pool
}

func NewPGRepository(pool *pgxpool.Pool) _interface.Products {
	return &pgRepo{p: pool}
}

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Verify connection
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx2); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func (r *pgRepo) GetList(ctx context.Context) ([]models.Product, error) {
	rows, err := r.p.Query(ctx, `SELECT id::text, category, name, price::double precision FROM products ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	out := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Category, &p.Name, &p.Price); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		out = append(out, p)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("rows err: %w", rows.Err())
	}
	return out, nil
}

func (r *pgRepo) Get(ctx context.Context, id string) (models.Product, error) {
	var p models.Product
	row := r.p.QueryRow(ctx, `SELECT id::text, category, name, price::double precision FROM products WHERE id = $1`, id)
	if err := row.Scan(&p.ID, &p.Category, &p.Name, &p.Price); err != nil {
		// check pgx error for not found
		if err.Error() == "no rows in result set" {
			return models.Product{}, ErrNotFound
		}
		return models.Product{}, fmt.Errorf("scan product: %w", err)
	}
	return p, nil
}

func (e Err) Error() string {
	return string(e)
}
