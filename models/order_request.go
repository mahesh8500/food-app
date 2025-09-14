package models

import "time"

type OrderRequest struct {
	ID         string    `db:"id"`
	Items      []byte    `db:"items"`
	CouponCode *string   `db:"coupon_code"`
	CreatedAt  time.Time `db:"created_at"`
}
