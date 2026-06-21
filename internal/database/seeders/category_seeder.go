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

	categories := []domain.Category{
		{StoreID: 1, Name: "Kopi & Espresso"},
		{StoreID: 1, Name: "Minuman Segar"},
		{StoreID: 1, Name: "Makanan Utama"},
		{StoreID: 1, Name: "Roti & Pastry"},
		{StoreID: 1, Name: "Camilan & Dessert"},
	}

	_, err = db.NewInsert().Model(&categories).Exec(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Categories seeded successfully")
	return nil
}
