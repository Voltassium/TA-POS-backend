package seeders

import (
	"backend-ta/app/domain"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func SeedAll(db *bun.DB, fresh bool) error {
	ctx := context.Background()

	if fresh {
		fmt.Println("[SEEDER] Fresh mode: truncating tables...")
		models := []interface{}{
			(*domain.Product)(nil),
			(*domain.Category)(nil),
			(*domain.User)(nil),
			(*domain.Store)(nil),
		}
		for _, model := range models {
			_, err := db.NewTruncateTable().Model(model).Cascade().Exec(ctx)
			if err != nil {
				return fmt.Errorf("error truncating table: %w", err)
			}
		}
		fmt.Println("[SEEDER] Tables truncated successfully")
	}

	fmt.Println("[SEEDER] Starting seeding...")

	if err := SeedStores(ctx, db); err != nil {
		return fmt.Errorf("error seeding stores: %w", err)
	}

	if err := SeedUsers(ctx, db); err != nil {
		return fmt.Errorf("error seeding users: %w", err)
	}

	if err := SeedCategories(ctx, db); err != nil {
		return fmt.Errorf("error seeding categories: %w", err)
	}

	if err := SeedProducts(ctx, db); err != nil {
		return fmt.Errorf("error seeding products: %w", err)
	}

	if err := SeedPengeluaran(ctx, db); err != nil {
		return fmt.Errorf("error seeding pengeluaran: %w", err)
	}

	if err := SeedOrders(ctx, db); err != nil {
		return fmt.Errorf("error seeding orders: %w", err)
	}

	fmt.Println("[SEEDER] Seeding completed successfully!")
	return nil
}
