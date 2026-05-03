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
		{StoreID: 1, Name: "Minuman", ImageURL: "https://placehold.co/400x400?text=Minuman"},
		{StoreID: 1, Name: "Makanan", ImageURL: "https://placehold.co/400x400?text=Makanan"},
		{StoreID: 1, Name: "Cemilan", ImageURL: "https://placehold.co/400x400?text=Cemilan"},
		{StoreID: 1, Name: "Pencuci Mulut", ImageURL: "https://placehold.co/400x400?text=Pencuci+Mulut"},
	}

	_, err = db.NewInsert().Model(&categories).Exec(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Categories seeded successfully")
	return nil
}
