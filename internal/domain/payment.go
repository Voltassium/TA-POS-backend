package domain

import (
	"backend-ta/internal/constants"
	"time"

	"github.com/uptrace/bun"
)

type Payment struct {
	bun.BaseModel `bun:"table:payments"`

	ID            int64                   `bun:"id,pk,autoincrement"`
	OrderID       int64                   `bun:"order_id,notnull,unique"`
	PaymentMethod constants.PaymentMethod `bun:"payment_method,notnull"`
	AmountPaid    float64                 `bun:"amount_paid,notnull"`
	Timestamp     time.Time               `bun:"timestamp,nullzero,notnull,default:current_timestamp"`
}
