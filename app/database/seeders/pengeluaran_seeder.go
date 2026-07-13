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
	count, err := db.NewSelect().Model((*domain.Pengeluaran)(nil)).Count(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		fmt.Println("[SEEDER] Pengeluaran table already has data, skipping...")
		return nil
	}

	var owners []domain.User
	err = db.NewSelect().Model(&owners).Where("role = ?", constants.UserRoleOwner).Scan(ctx)
	if err != nil {
		return fmt.Errorf("owners not found: %w", err)
	}

	ownerMap := make(map[int64]string)
	for _, o := range owners {
		if o.StoreID != nil {
			ownerMap[*o.StoreID] = o.ID
		}
	}

	startDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2026, 5, 31, 23, 59, 59, 0, time.Local)

	var pengeluarans []domain.Pengeluaran

	for storeID := int64(1); storeID <= 1; storeID++ {
		ownerID, ok := ownerMap[storeID]
		if !ok {
			continue
		}

		for d := startDate; !d.After(endDate); d = d.AddDate(0, 1, 0) {
			pengeluarans = append(pengeluarans,
				domain.Pengeluaran{
					StoreID:     storeID,
					Tanggal:     d.Add(time.Hour * 10),
					Category:    "Listrik & Air",
					Description: sp("Tagihan listrik dan air bulanan"),
					Amount:      1500000 + float64(rand.Intn(500000)),
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
					Category:    "Bahan Baku",
					Description: sp("Belanja bahan baku nasi padang"),
					Amount:      3000000 + float64(rand.Intn(1500000)),
					CreatedBy:   ownerID,
					CreatedAt:   d.Add(time.Hour * 9),
					UpdatedAt:   d.Add(time.Hour * 9),
				},
				domain.Pengeluaran{
					StoreID:     storeID,
					Tanggal:     d.Add(time.Hour * 11),
					Category:    "Operasional",
					Description: sp("Kotak kemasan, sendok, plastik, tisu"),
					Amount:      500000 + float64(rand.Intn(300000)),
					CreatedBy:   ownerID,
					CreatedAt:   d.Add(time.Hour * 11),
					UpdatedAt:   d.Add(time.Hour * 11),
				},
			)
		}
	}

	if len(pengeluarans) > 0 {
		_, err = db.NewInsert().Model(&pengeluarans).Exec(ctx)
		if err != nil {
			return err
		}
	}

	fmt.Println("[SEEDER] Pengeluaran seeded successfully")
	return nil
}

func sp(s string) *string {
	return &s
}
