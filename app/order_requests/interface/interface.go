package _interface

import (
	"context"
	"foodApp/models"
)

type OrderRequests interface {
	Insert(ctx context.Context, req models.OrderRequest) (models.OrderRequest, error)
}
