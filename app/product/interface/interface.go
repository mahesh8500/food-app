package _interface

import (
	"context"
	"foodApp/models"
)

type Products interface {
	GetList(ctx context.Context) ([]models.Product, error)
	Get(ctx context.Context, id string) (models.Product, error)
}
