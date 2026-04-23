package seeders

import (
	"backend-ta/internal/domain"
	"context"
	"fmt"

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

	categoryMap := make(map[string]int64)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
	}

	products := []domain.Product{
		// Minuman
		{CategoryID: categoryMap["Minuman"], Name: "Kopi Hitam", Price: 15000, Stock: 100, Description: "Kopi hitam murni panas"},
		{CategoryID: categoryMap["Minuman"], Name: "Es Teh Manis", Price: 5000, Stock: 200, Description: "Teh manis dengan es segar"},
		{CategoryID: categoryMap["Minuman"], Name: "Jus Jeruk", Price: 12000, Stock: 50, Description: "Perasan jeruk asli segar"},
		
		// Makanan
		{CategoryID: categoryMap["Makanan"], Name: "Nasi Goreng Spesial", Price: 25000, Stock: 50, Description: "Nasi goreng dengan telur dan ayam"},
		{CategoryID: categoryMap["Makanan"], Name: "Mie Ayam", Price: 15000, Stock: 60, Description: "Mie ayam dengan pangsit"},
		{CategoryID: categoryMap["Makanan"], Name: "Sate Ayam", Price: 20000, Stock: 40, Description: "Sate ayam isi 10 tusuk"},
		
		// Cemilan
		{CategoryID: categoryMap["Cemilan"], Name: "Kentang Goreng", Price: 10000, Stock: 100, Description: "Kentang goreng renyah"},
		{CategoryID: categoryMap["Cemilan"], Name: "Pisang Goreng", Price: 8000, Stock: 70, Description: "Pisang goreng manis isi 3"},
		
		// Pencuci Mulut
		{CategoryID: categoryMap["Pencuci Mulut"], Name: "Es Krim Cokelat", Price: 10000, Stock: 30, Description: "Es krim cokelat lembut"},
		{CategoryID: categoryMap["Pencuci Mulut"], Name: "Puding Susu", Price: 7000, Stock: 50, Description: "Puding susu manis"},
	}

	_, err = db.NewInsert().Model(&products).Exec(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Products seeded successfully")
	return nil
}
