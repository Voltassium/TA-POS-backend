package domain

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products"`

	ID          int64     `bun:"id,pk,autoincrement"`
	StoreID     int64     `bun:"store_id,notnull"`
	CategoryID  int64     `bun:"category_id,notnull"`
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
