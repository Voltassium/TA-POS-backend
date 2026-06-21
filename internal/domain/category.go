package domain

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:categories"`

	ID        string    `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StoreID   int64     `bun:"store_id,notnull"`
	Name      string    `bun:"name,notnull"`
	Products  []Product `bun:"rel:has-many,join:id=category_id"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func (m *Category) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
