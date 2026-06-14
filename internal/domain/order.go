package domain

import (
	"backend-ta/internal/constants"
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:orders"`

	ID             int64                 `bun:"id,pk,autoincrement"`
	OrderCode      string                `bun:"order_code,notnull"`
	StoreID        int64                 `bun:"store_id,notnull"`
	TableID        *int64                `bun:"table_id"`
	StaffID        int64                 `bun:"staff_id,notnull"`
	TotalAmount    float64               `bun:"total_amount,notnull,default:0"`
	DiscountType   *string               `bun:"discount_type"`
	DiscountValue  float64               `bun:"discount_value,notnull,default:0"`
	DiscountAmount float64               `bun:"discount_amount,notnull,default:0"`
	Status         constants.OrderStatus `bun:"status,notnull"`
	OrderItems     []OrderItem           `bun:"rel:has-many,join:id=order_id"`
	Payment        *Payment              `bun:"rel:has-one,join:id=order_id"`
	CreatedAt      time.Time             `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt      time.Time             `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func (m *Order) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
