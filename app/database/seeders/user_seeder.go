package seeders

import (
	"backend-ta/app/constants"
	"backend-ta/app/domain"
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
		fmt.Println("[SEEDER] Users table already has data, skipping...")
		return nil
	}

	hashedPassword, err := authentication.HashPassword("password123")
	if err != nil {
		log.Fatal("Failed to hash password", err)
	}

	var users []domain.User

	store1 := int64(1)
	users = append(users, domain.User{
		Email:    "admin@pos.com",
		Password: hashedPassword,
		Role:     constants.UserRoleSuperadmin,
		StoreID:  &store1,
	})

	for i := int64(1); i <= 4; i++ {
		currStore := i
		var ownerEmail, chefEmail, staffEmail string
		if i == 1 {
			ownerEmail = "owner@pos.com"
			chefEmail = "chef@pos.com"
			staffEmail = "staff@pos.com"
		} else {
			ownerEmail = fmt.Sprintf("owner%d@pos.com", i)
			chefEmail = fmt.Sprintf("chef%d@pos.com", i)
			staffEmail = fmt.Sprintf("staff%d@pos.com", i)
		}

		users = append(users, domain.User{
			Email:    ownerEmail,
			Password: hashedPassword,
			Role:     constants.UserRoleOwner,
			StoreID:  &currStore,
		}, domain.User{
			Email:    chefEmail,
			Password: hashedPassword,
			Role:     constants.UserRoleChef,
			StoreID:  &currStore,
		}, domain.User{
			Email:    staffEmail,
			Password: hashedPassword,
			Role:     constants.UserRoleStaff,
			StoreID:  &currStore,
		})
	}

	_, err = db.NewInsert().Model(&users).Exec(ctx)
	if err != nil {
		return err
	}

	fmt.Println("[SEEDER] Users seeded successfully")
	return nil
}
