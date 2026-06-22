package seeders

import (
	"backend-ta/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

func SeedProducts(ctx context.Context, db *bun.DB) error {
	count, err := db.NewSelect().Model((*domain.Product)(nil)).Count(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		fmt.Println("Products table already has data, skipping...")
		return nil
	}

	var categories []domain.Category
	err = db.NewSelect().Model(&categories).Order("id ASC").Scan(ctx)
	if err != nil {
		return err
	}

	if len(categories) == 0 {
		return fmt.Errorf("no categories found, please seed categories first")
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

	for storeID := int64(1); storeID <= 4; storeID++ {
		makananID := categoryMap[catKey{StoreID: storeID, Name: "Makanan Utama"}]
		laukID := categoryMap[catKey{StoreID: storeID, Name: "Lauk Sampingan"}]
		minumanID := categoryMap[catKey{StoreID: storeID, Name: "Minuman & Jus"}]
		camilanID := categoryMap[catKey{StoreID: storeID, Name: "Camilan & Dessert"}]

		products = append(products,
			// Makanan Utama (Olahan)
			domain.Product{StoreID: storeID, CategoryID: makananID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MAKAN-%d01", storeID)), Name: "Rendang Sapi Minang", Price: 28000, Stock: 100, Description: "Daging sapi pilihan dimasak perlahan dengan santan dan rempah otentik Minang"},
			domain.Product{StoreID: storeID, CategoryID: makananID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MAKAN-%d02", storeID)), Name: "Ayam Pop Khas Padang", Price: 24000, Stock: 80, Description: "Ayam kampung gurih khas Minang dengan tekstur lembut, disajikan dengan sambal pop khas"},
			domain.Product{StoreID: storeID, CategoryID: makananID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MAKAN-%d03", storeID)), Name: "Gulai Tunjang (Kikil)", Price: 30000, Stock: 60, Description: "Kikil sapi tebal nan empuk berkuah gulai kuning gurih kaya bumbu rempah"},
			domain.Product{StoreID: storeID, CategoryID: makananID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MAKAN-%d04", storeID)), Name: "Dendeng Batokok", Price: 28000, Stock: 70, Description: "Daging sapi tipis digoreng garing yang ditumbuk kasar disajikan dengan siraman cabai merah/ijo"},
			domain.Product{StoreID: storeID, CategoryID: makananID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MAKAN-%d05", storeID)), Name: "Nasi Rames Rendang", Price: 38000, Stock: 150, Description: "Nasi padang lengkap dengan rendang sapi, sayur nangka gulai, daun singkong, dan sambal ijo"},

			// Lauk Sampingan (Olahan)
			domain.Product{StoreID: storeID, CategoryID: laukID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("LAUK-%d01", storeID)), Name: "Telur Dadar Barendo", Price: 15000, Stock: 200, Description: "Telur dadar khas padang yang tebal, garing berenda di luar, dan lembut berempah di dalam"},
			domain.Product{StoreID: storeID, CategoryID: laukID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("LAUK-%d02", storeID)), Name: "Perkedel Kentang Padang", Price: 8000, Stock: 120, Description: "Perkedel kentang padat dengan cita rasa gurih bawang goreng dan seledri segar"},
			domain.Product{StoreID: storeID, CategoryID: laukID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("LAUK-%d03", storeID)), Name: "Daun Singkong & Sambal Ijo", Price: 5000, Stock: 150, Description: "Porsi daun singkong rebus empuk dipadu sambal ijo padang legendaris"},
			domain.Product{StoreID: storeID, CategoryID: laukID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("LAUK-%d04", storeID)), Name: "Sayur Nangka Gulai", Price: 7000, Stock: 100, Description: "Sayur gulai nangka muda (cubadak) berkuah santan gurih meresap"},

			// Minuman & Jus (Olahan)
			domain.Product{StoreID: storeID, CategoryID: minumanID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MINUM-%d01", storeID)), Name: "Teh Talua (Teh Telur)", Price: 18000, Stock: 80, Description: "Minuman kesehatan khas Minang dari campuran kuning telur bebek kocok, teh pekat, dan susu kental manis"},
			domain.Product{StoreID: storeID, CategoryID: minumanID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MINUM-%d02", storeID)), Name: "Es Teh Manis Selera", Price: 6000, Stock: 300, Description: "Es teh manis segar pelepas dahaga setelah makan pedas"},
			domain.Product{StoreID: storeID, CategoryID: minumanID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MINUM-%d03", storeID)), Name: "Jus Alpukat Lumer", Price: 15000, Stock: 90, Description: "Jus alpukat mentega kental disiram susu cokelat kental manis melingkar"},
			domain.Product{StoreID: storeID, CategoryID: minumanID, ProductType: "Olahan", SKU: sp(fmt.Sprintf("MINUM-%d04", storeID)), Name: "Es Jeruk Peras Murni", Price: 10000, Stock: 100, Description: "Es jeruk segar dari perasan jeruk manis asli pilihan"},

			// Camilan & Dessert (Kulakan)
			domain.Product{StoreID: storeID, CategoryID: camilanID, ProductType: "Kulakan", HargaBeli: hb(8000), SKU: sp(fmt.Sprintf("CAMIL-%d01", storeID)), Name: "Keripik Sanjai Balado", Price: 15000, Stock: 80, Description: "Keripik singkong khas Bukittinggi dengan balutan bumbu balado basah merah pedas manis"},
			domain.Product{StoreID: storeID, CategoryID: camilanID, ProductType: "Kulakan", HargaBeli: hb(5000), SKU: sp(fmt.Sprintf("CAMIL-%d02", storeID)), Name: "Kerupuk Kulit Jangek", Price: 10000, Stock: 120, Description: "Kerupuk kulit sapi renyah gurih khas Minang yang cocok disantap bersama kuah gulai"},
			domain.Product{StoreID: storeID, CategoryID: camilanID, ProductType: "Kulakan", HargaBeli: hb(6000), SKU: sp(fmt.Sprintf("CAMIL-%d03", storeID)), Name: "Roti Cane Susu", Price: 12000, Stock: 50, Description: "Roti cane khas Padang disiram susu kental manis legit"},
		)
	}

	_, err = db.NewInsert().Model(&products).Exec(ctx)
	if err != nil {
		return err
	}

	var stockHistories []domain.StockHistory
	febFirst := time.Date(2026, 2, 1, 8, 0, 0, 0, time.Local)

	var insertedProducts []domain.Product
	db.NewSelect().Model(&insertedProducts).Scan(ctx)

	for _, p := range insertedProducts {
		if p.Stock > 0 {
			stockHistories = append(stockHistories, domain.StockHistory{
				ProductID: p.ID,
				Change:    p.Stock,
				Reason:    "Inisialisasi Stok Awal (Seeder)",
				CreatedAt: febFirst,
			})
		}
	}

	if len(stockHistories) > 0 {
		_, err = db.NewInsert().Model(&stockHistories).Exec(ctx)
		if err != nil {
			return err
		}
	}

	fmt.Println("Products and initial StockHistory seeded successfully")
	return nil
}
