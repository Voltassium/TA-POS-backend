package seeders

import (
	"backend-ta/app/domain"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func SeedStores(ctx context.Context, db *bun.DB) error {
	stores := []domain.Store{
		{
			ID:      1,
			Name:    "Rumah Padang - Sudirman",
			Address: "Jl. Jend. Sudirman Kav. 52-53, Senayan, Kebayoran Baru, Jakarta Selatan",
		},
	}

	for _, store := range stores {
		_, err := db.NewInsert().Model(&store).On("CONFLICT (id) DO UPDATE").Exec(ctx)
		if err != nil {
			return err
		}
	}

	fmt.Println("[SEEDER] Stores seeded successfully")
	return nil
}
