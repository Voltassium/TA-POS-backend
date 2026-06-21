package domain

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products"`

	ID          string    `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StoreID     int64     `bun:"store_id,notnull"`
	ProductType string    `bun:"product_type,notnull,default:'Olahan'"`
	CategoryID  string    `bun:"category_id,type:uuid,notnull"`
	SKU         *string   `bun:"sku"`
	HargaBeli   *float64  `bun:"harga_beli"`
	Name        string    `bun:"name,notnull"`
	Description string    `bun:"description"`
	Price       float64   `bun:"price,notnull"`
	IsAvailable bool      `bun:"is_available,notnull,default:true"`
	Stock       int       `bun:"stock,notnull,default:0"`
	Category    *Category `bun:"rel:belongs-to,join:category_id=id"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func (m *Product) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
