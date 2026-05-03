package domain

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Store struct {
	bun.BaseModel `bun:"table:stores"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Name      string    `bun:"name,notnull"`
	Address   string    `bun:"address"`
	Phone     string    `bun:"phone"`
	LogoURL   string    `bun:"logo_url"`
	IsActive  bool      `bun:"is_active,notnull,default:true"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func (m *Store) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
