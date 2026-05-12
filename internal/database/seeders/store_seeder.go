package seeders

import (
	"backend-ta/internal/domain"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func SeedStores(ctx context.Context, db *bun.DB) error {
	stores := []domain.Store{
		{
			ID:      1,
			Name:    "Kopi Senja - Sudirman",
			Address: "Jl. Jend. Sudirman Kav. 52-53, Senayan, Kebayoran Baru, Jakarta Selatan",
		},
		{
			ID:      2,
			Name:    "Kopi Senja - Kemang",
			Address: "Jl. Kemang Raya No. 45, Bangka, Mampang Prapatan, Jakarta Selatan",
		},
	}

	for _, store := range stores {
		_, err := db.NewInsert().Model(&store).On("CONFLICT (id) DO UPDATE").Exec(ctx)
		if err != nil {
			return err
		}
	}

	fmt.Println("Stores seeded successfully")
	return nil
}
