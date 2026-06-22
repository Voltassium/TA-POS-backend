package seeders

import (
	"backend-ta/internal/domain"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func SeedCategories(ctx context.Context, db *bun.DB) error {
	count, err := db.NewSelect().Model((*domain.Category)(nil)).Count(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		fmt.Println("Categories table already has data, skipping...")
		return nil
	}

	var categories []domain.Category
	for storeID := int64(1); storeID <= 4; storeID++ {
		categories = append(categories,
			domain.Category{StoreID: storeID, Name: "Makanan Utama"},
			domain.Category{StoreID: storeID, Name: "Lauk Sampingan"},
			domain.Category{StoreID: storeID, Name: "Minuman & Jus"},
			domain.Category{StoreID: storeID, Name: "Camilan & Dessert"},
		)
	}

	_, err = db.NewInsert().Model(&categories).Exec(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Categories seeded successfully")
	return nil
}
