package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type StockHistory struct {
	bun.BaseModel `bun:"table:stock_histories"`

	ID        int64     `bun:"id,pk,autoincrement"`
	ProductID int64     `bun:"product_id,notnull"`
	Change    int       `bun:"change,notnull"`
	Reason    string    `bun:"reason,notnull"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`

	Product *Product `bun:"rel:belongs-to,join:product_id=id"`
}
