package seeders

import (
	"backend-ta/internal/domain"
	"context"

	"github.com/uptrace/bun"
)

func SeedStores(ctx context.Context, db *bun.DB) error {
	stores := []domain.Store{
		{
			ID:       1,
			Name:     "Main Store",
			Address:  "123 Main Street, City",
			Phone:    "+6281234567890",
			IsActive: true,
		},
		{
			ID:       2,
			Name:     "Branch Store",
			Address:  "456 Branch Avenue, City",
			Phone:    "+6280987654321",
			IsActive: true,
		},
	}

	for _, store := range stores {
		_, err := db.NewInsert().Model(&store).On("CONFLICT (id) DO UPDATE").Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
