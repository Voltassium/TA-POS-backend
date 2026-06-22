package domain

import (
	"backend-ta/app/constants"
	"time"

	"github.com/uptrace/bun"
)

type Payment struct {
	bun.BaseModel `bun:"table:payments"`

	ID            string                  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	OrderID       string                  `bun:"order_id,notnull,unique,type:uuid"`
	PaymentMethod constants.PaymentMethod `bun:"payment_method,notnull"`
	AmountPaid    float64                 `bun:"amount_paid,notnull"`
	Timestamp     time.Time               `bun:"timestamp,nullzero,notnull,default:current_timestamp"`
}
