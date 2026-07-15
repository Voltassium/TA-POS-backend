package seeders

import (
	"backend-ta/app/domain"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func SeedCategories(ctx context.Context, db *bun.DB) error {
	var stores []domain.Store
	err := db.NewSelect().Model(&stores).Scan(ctx)
	if err != nil {
		return err
	}

	var categories []domain.Category
	for _, store := range stores {
		count, err := db.NewSelect().
			Model((*domain.Category)(nil)).
			Where("store_id = ?", store.ID).
			Count(ctx)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		categories = append(categories,
			domain.Category{StoreID: store.ID, Name: "Makanan Utama"},
			domain.Category{StoreID: store.ID, Name: "Lauk Sampingan"},
			domain.Category{StoreID: store.ID, Name: "Minuman & Jus"},
			domain.Category{StoreID: store.ID, Name: "Camilan & Dessert"},
		)
	}

	if len(categories) > 0 {
		_, err = db.NewInsert().Model(&categories).Exec(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("[SEEDER] Categories seeded successfully: %d categories created\n", len(categories))
	} else {
		fmt.Println("[SEEDER] Categories seeding skipped (already has data for all stores)")
	}

	return nil
}
