package domain

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type OrderItem struct {
	bun.BaseModel `bun:"table:order_items"`

	ID             int64    `bun:"id,pk,autoincrement"`
	OrderID        int64    `bun:"order_id,notnull"`
	ProductID      int64    `bun:"product_id,notnull"`
	Quantity       int      `bun:"quantity,notnull"`
	UnitPrice      float64  `bun:"unit_price,notnull"`
	DiscountType   *string  `bun:"discount_type"`
	DiscountValue  float64  `bun:"discount_value,notnull,default:0"`
	DiscountAmount float64  `bun:"discount_amount,notnull,default:0"`
	Subtotal       float64  `bun:"subtotal,notnull"`
	ServedQty      int      `bun:"served_qty,notnull,default:0"`
	Product        *Product `bun:"rel:belongs-to,join:product_id=id"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func (m *OrderItem) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
