package seeders

import (
	"backend-ta/app/domain"
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

func SeedProducts(ctx context.Context, db *bun.DB) error {
	var stores []domain.Store
	err := db.NewSelect().Model(&stores).Scan(ctx)
	if err != nil {
		return err
	}

	var categories []domain.Category
	err = db.NewSelect().Model(&categories).Order("id ASC").Scan(ctx)
	if err != nil {
		return err
	}

	type catKey struct {
		StoreID int64
		Name    string
	}
	categoryMap := make(map[catKey]string)
	for _, cat := range categories {
		categoryMap[catKey{StoreID: cat.StoreID, Name: cat.Name}] = cat.ID
	}

	hb := func(v float64) *float64 { return &v }
	sp := func(v string) *string { return &v }

	var products []domain.Product

	for _, store := range stores {
		storeID := store.ID
		count, err := db.NewSelect().
			Model((*domain.Product)(nil)).
			Where("store_id = ?", storeID).
			Count(ctx)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		makananID := categoryMap[catKey{StoreID: storeID, Name: "Makanan Utama"}]
		laukID := categoryMap[catKey{StoreID: storeID, Name: "Lauk Sampingan"}]
		minumanID := categoryMap[catKey{StoreID: storeID, Name: "Minuman & Jus"}]
		camilanID := categoryMap[catKey{StoreID: storeID, Name: "Camilan & Dessert"}]

		if makananID == "" || laukID == "" || minumanID == "" || camilanID == "" {
			continue
		}

		products = append(products,
			// Makanan Utama (Olahan)
			domain.Product{StoreID: storeID, CategoryID: makananID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MAKAN-%d01", storeID)), Name: "Rendang Sapi Minang", Price: 28000, Stock: 100},
			domain.Product{StoreID: storeID, CategoryID: makananID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MAKAN-%d02", storeID)), Name: "Ayam Pop Khas Padang", Price: 24000, Stock: 80},
			domain.Product{StoreID: storeID, CategoryID: makananID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MAKAN-%d03", storeID)), Name: "Gulai Tunjang (Kikil)", Price: 30000, Stock: 60},
			domain.Product{StoreID: storeID, CategoryID: makananID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MAKAN-%d04", storeID)), Name: "Dendeng Batokok", Price: 28000, Stock: 70},
			domain.Product{StoreID: storeID, CategoryID: makananID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MAKAN-%d05", storeID)), Name: "Nasi Rames Rendang", Price: 38000, Stock: 150},

			// Lauk Sampingan (Olahan)
			domain.Product{StoreID: storeID, CategoryID: laukID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("LAUK-%d01", storeID)), Name: "Telur Dadar Barendo", Price: 15000, Stock: 200},
			domain.Product{StoreID: storeID, CategoryID: laukID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("LAUK-%d02", storeID)), Name: "Perkedel Kentang Padang", Price: 8000, Stock: 120},
			domain.Product{StoreID: storeID, CategoryID: laukID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("LAUK-%d03", storeID)), Name: "Daun Singkong & Sambal Ijo", Price: 5000, Stock: 150},
			domain.Product{StoreID: storeID, CategoryID: laukID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("LAUK-%d04", storeID)), Name: "Sayur Nangka Gulai", Price: 7000, Stock: 100},

			// Minuman & Jus (Olahan)
			domain.Product{StoreID: storeID, CategoryID: minumanID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MINUM-%d01", storeID)), Name: "Teh Talua (Teh Telur)", Price: 18000, Stock: 80},
			domain.Product{StoreID: storeID, CategoryID: minumanID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MINUM-%d02", storeID)), Name: "Es Teh Manis Selera", Price: 6000, Stock: 300},
			domain.Product{StoreID: storeID, CategoryID: minumanID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MINUM-%d03", storeID)), Name: "Jus Alpukat Lumer", Price: 15000, Stock: 90},
			domain.Product{StoreID: storeID, CategoryID: minumanID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MINUM-%d04", storeID)), Name: "Es Jeruk Peras Murni", Price: 10000, Stock: 100},

			// Camilan & Dessert (Kulakan)
			domain.Product{StoreID: storeID, CategoryID: camilanID, ProductType: "Kulakan", HargaBeli: hb(8000), SKU: sp(fmt.Sprintf("CAMIL-%d01", storeID)), Name: "Keripik Sanjai Balado", Price: 15000, Stock: 80},
			domain.Product{StoreID: storeID, CategoryID: camilanID, ProductType: "Kulakan", HargaBeli: hb(5000), SKU: sp(fmt.Sprintf("CAMIL-%d02", storeID)), Name: "Kerupuk Kulit Jangek", Price: 10000, Stock: 120},
			domain.Product{StoreID: storeID, CategoryID: camilanID, ProductType: "Kulakan", HargaBeli: hb(6000), SKU: sp(fmt.Sprintf("CAMIL-%d03", storeID)), Name: "Roti Cane Susu", Price: 12000, Stock: 50},
		)
	}

	if len(products) > 0 {
		_, err = db.NewInsert().Model(&products).Exec(ctx)
		if err != nil {
			return err
		}

		var stockHistories []domain.StockHistory
		febFirst := time.Date(2026, 2, 1, 8, 0, 0, 0, time.Local)

		for _, p := range products {
			if p.Stock > 0 {
				stockHistories = append(stockHistories, domain.StockHistory{
					ProductID:    p.ID,
					Change:       p.Stock,
					InitialStock: 0,
					FinalStock:   p.Stock,
					Reason:       "Inisialisasi Stok Awal (Seeder)",
					CreatedAt:    febFirst,
				})
			}
		}

		if len(stockHistories) > 0 {
			_, err = db.NewInsert().Model(&stockHistories).Exec(ctx)
			if err != nil {
				return err
			}
		}
		fmt.Printf("[SEEDER] Products and initial StockHistory seeded successfully: %d products created\n", len(products))
	} else {
		fmt.Println("[SEEDER] Products seeding skipped (already has data for all stores)")
	}

	return nil
}
