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
		{Name: "Minuman", ImageURL: "https://placehold.co/400x400?text=Minuman"},
		{Name: "Makanan", ImageURL: "https://placehold.co/400x400?text=Makanan"},
		{Name: "Cemilan", ImageURL: "https://placehold.co/400x400?text=Cemilan"},
		{Name: "Pencuci Mulut", ImageURL: "https://placehold.co/400x400?text=Pencuci+Mulut"},
	}

	_, err = db.NewInsert().Model(&categories).Exec(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Categories seeded successfully")
	return nil
}
