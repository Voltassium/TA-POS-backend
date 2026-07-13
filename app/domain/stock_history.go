package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type StockHistory struct {
	bun.BaseModel `bun:"table:stock_histories"`

	ID         int64     `bun:"id,pk,autoincrement"`
	ProductID  string    `bun:"product_id,type:uuid,notnull"`
	Change     int       `bun:"change,notnull"`
	Reason     string    `bun:"reason,notnull"`
	SourceType string    `bun:"source_type,notnull,default:'manual'"`
	HargaBeli  float64   `bun:"harga_beli,notnull,default:0"`
	CreatedAt  time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`

	InitialStock int `bun:"initial_stock,notnull,default:0"`
	FinalStock   int `bun:"final_stock,notnull,default:0"`

	Product *Product `bun:"rel:belongs-to,join:product_id=id"`
}
