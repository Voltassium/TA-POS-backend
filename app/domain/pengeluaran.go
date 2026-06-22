package domain

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Pengeluaran struct {
	bun.BaseModel `bun:"table:pengeluaran"`

	ID          string    `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StoreID     int64     `bun:"store_id,notnull"`
	Tanggal     time.Time `bun:"tanggal,notnull"`
	Category    string    `bun:"category,notnull"`
	Description *string   `bun:"description"`
	Amount      float64   `bun:"amount,notnull"`
	CreatedBy   string    `bun:"created_by,type:uuid,notnull"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func (m *Pengeluaran) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
