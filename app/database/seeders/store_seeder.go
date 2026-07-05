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
		{
			ID:      2,
			Name:    "Rumah Padang - Kemang",
			Address: "Jl. Kemang Raya No. 45, Bangka, Mampang Prapatan, Jakarta Selatan",
		},
		{
			ID:      3,
			Name:    "Rumah Padang - Senopati",
			Address: "Jl. Senopati No. 72, Selong, Kebayoran Baru, Jakarta Selatan",
		},
		{
			ID:      4,
			Name:    "Rumah Padang - Dago",
			Address: "Jl. Ir. H. Juanda No. 102, Dago, Coblong, Kota Bandung",
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
