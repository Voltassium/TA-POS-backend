package domain

import (
	"backend-ta/internal/constants"
	"context"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID        string             `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StoreID   *int64             `bun:"store_id"`
	Email     string             `bun:"email,notnull,unique"`
	Password  string             `bun:"password,notnull"`
	Role      constants.UserRole `bun:"role,notnull"`
	Store     *Store             `bun:"rel:belongs-to,join:store_id=id"`
	CreatedAt time.Time          `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time          `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

func (m *User) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
