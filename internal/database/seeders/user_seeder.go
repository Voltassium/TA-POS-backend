package seeders

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/domain"
	"backend-ta/pkg/authentication"
	"context"
	"fmt"
	"log"

	"github.com/uptrace/bun"
)

func SeedUsers(ctx context.Context, db *bun.DB) error {
	count, err := db.NewSelect().Model((*domain.User)(nil)).Count(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		fmt.Println("Users table already has data, skipping...")
		return nil
	}

	hashedPassword, err := authentication.HashPassword("password123")
	if err != nil {
		log.Fatal("Failed to hash password", err)
	}

	storeID := int64(1)
	users := []domain.User{
		{
			Email:    "admin@pos.com",
			Password: hashedPassword,
			Role:     constants.UserRoleSuperadmin,
			StoreID:  &storeID,
		},
		{
			Email:    "staff@pos.com",
			Password: hashedPassword,
			Role:     constants.UserRoleStaff,
			StoreID:  &storeID,
		},
	}

	_, err = db.NewInsert().Model(&users).Exec(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Users seeded successfully")
	return nil
}
