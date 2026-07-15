package seeders

import (
	"backend-ta/app/constants"
	"backend-ta/app/domain"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/uptrace/bun"
)

func SeedPengeluaran(ctx context.Context, db *bun.DB) error {
	var stores []domain.Store
	err := db.NewSelect().Model(&stores).Scan(ctx)
	if err != nil {
		return err
	}

	var users []domain.User
	err = db.NewSelect().Model(&users).Where("role = ? OR role = ?", constants.UserRoleOwner, constants.UserRoleSuperadmin).Scan(ctx)
	if err != nil {
		return fmt.Errorf("users not found: %w", err)
	}

	ownerMap := make(map[int64]string)
	var fallbackUserID string
	for _, u := range users {
		if fallbackUserID == "" {
			fallbackUserID = u.ID
		}
		if u.StoreID != nil {
			ownerMap[*u.StoreID] = u.ID
		}
	}

	startDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2026, 5, 31, 23, 59, 59, 0, time.Local)

	var pengeluarans []domain.Pengeluaran

	for _, store := range stores {
		storeID := store.ID
		count, err := db.NewSelect().
			Model((*domain.Pengeluaran)(nil)).
			Where("store_id = ?", storeID).
			Count(ctx)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		ownerID, ok := ownerMap[storeID]
		if !ok {
			ownerID = fallbackUserID
		}
		if ownerID == "" {
			continue
		}

		for d := startDate; !d.After(endDate); d = d.AddDate(0, 1, 0) {
			pengeluarans = append(pengeluarans,
				domain.Pengeluaran{
					StoreID:     storeID,
					Tanggal:     d.Add(time.Hour * 10),
					Category:    "Listrik & Air",
					Description: sp("Tagihan listrik dan air bulanan"),
					Amount:      1500000 + float64(rand.Intn(500)*1000),
					CreatedBy:   ownerID,
					CreatedAt:   d.Add(time.Hour * 10),
					UpdatedAt:   d.Add(time.Hour * 10),
				},
				domain.Pengeluaran{
					StoreID:     storeID,
					Tanggal:     d.Add(time.Hour * 10),
					Category:    "Gaji Karyawan",
					Description: sp("Gaji staff dan chef"),
					Amount:      10000000,
					CreatedBy:   ownerID,
					CreatedAt:   d.Add(time.Hour * 10),
					UpdatedAt:   d.Add(time.Hour * 10),
				},
				domain.Pengeluaran{
					StoreID:     storeID,
					Tanggal:     d.Add(time.Hour * 10),
					Category:    "Sewa Tempat",
					Description: sp("Sewa bulanan ruko"),
					Amount:      5000000,
					CreatedBy:   ownerID,
					CreatedAt:   d.Add(time.Hour * 10),
					UpdatedAt:   d.Add(time.Hour * 10),
				},
			)
		}

		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 7) {
			pengeluarans = append(pengeluarans,
				domain.Pengeluaran{
					StoreID:     storeID,
					Tanggal:     d.Add(time.Hour * 9),
					Category:    "Belanja Harian",
					Description: sp("Belanja bahan baku nasi padang"),
					Amount:      3000000 + float64(rand.Intn(1500)*1000),
					CreatedBy:   ownerID,
					CreatedAt:   d.Add(time.Hour * 9),
					UpdatedAt:   d.Add(time.Hour * 9),
				},
				domain.Pengeluaran{
					StoreID:     storeID,
					Tanggal:     d.Add(time.Hour * 11),
					Category:    "Operasional",
					Description: sp("Kotak kemasan, sendok, plastik, tisu"),
					Amount:      500000 + float64(rand.Intn(300)*1000),
					CreatedBy:   ownerID,
					CreatedAt:   d.Add(time.Hour * 11),
					UpdatedAt:   d.Add(time.Hour * 11),
				},
			)
		}

		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			pengeluarans = append(pengeluarans,
				domain.Pengeluaran{
					StoreID:     storeID,
					Tanggal:     d.Add(time.Hour * 8),
					Category:    "Belanja Harian",
					Description: sp("Belanja kebutuhan harian operasional toko"),
					Amount:      30000 + float64(rand.Intn(771)*1000),
					CreatedBy:   ownerID,
					CreatedAt:   d.Add(time.Hour * 8),
					UpdatedAt:   d.Add(time.Hour * 8),
				},
			)
		}
	}

	if len(pengeluarans) > 0 {
		_, err = db.NewInsert().Model(&pengeluarans).Exec(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("[SEEDER] Pengeluaran seeded successfully: %d items created\n", len(pengeluarans))
	} else {
		fmt.Println("[SEEDER] Pengeluaran seeding skipped (already has data for all stores)")
	}

	return nil
}

func sp(s string) *string {
	return &s
}
