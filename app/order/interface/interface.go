package _interface

import (
	"context"
	"foodApp/models"
)

type Orders interface {
	Insert(ctx context.Context, order models.Order) (models.Order, error)
}
