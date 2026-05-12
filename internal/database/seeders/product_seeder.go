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
		// Kopi & Espresso
		{StoreID: 1, CategoryID: categoryMap["Kopi & Espresso"], Name: "Café Latte Panas", Price: 35000, Stock: 150, Description: "Espresso shot premium dengan susu kukus bertekstur lembut dan lapisan busa tipis"},
		{StoreID: 1, CategoryID: categoryMap["Kopi & Espresso"], Name: "Es Kopi Susu Senja", Price: 28000, Stock: 200, Description: "Kopi susu khas Senja dipadukan dengan gula aren murni dan krim lembut yang creamy"},
		{StoreID: 1, CategoryID: categoryMap["Kopi & Espresso"], Name: "Americano Dingin", Price: 25000, Stock: 120, Description: "Ekstraksi ganda biji kopi House Blend disajikan dingin menyegarkan tanpa gula"},
		{StoreID: 1, CategoryID: categoryMap["Kopi & Espresso"], Name: "Caramel Macchiato", Price: 42000, Stock: 90, Description: "Perpaduan sirup vanila, espresso, susu segar, dan siraman saus karamel lezat di atasnya"},

		// Minuman Segar
		{StoreID: 1, CategoryID: categoryMap["Minuman Segar"], Name: "Lychee Tea Minty", Price: 30000, Stock: 110, Description: "Teh hitam rasa leci dengan buah leci utuh segar dan sentuhan daun mint harum"},
		{StoreID: 1, CategoryID: categoryMap["Minuman Segar"], Name: "Matcha Latte Creamy", Price: 38000, Stock: 85, Description: "Bubuk teh hijau Uji Matcha asli Jepang diseduh dengan susu segar kaya rasa"},
		{StoreID: 1, CategoryID: categoryMap["Minuman Segar"], Name: "Tropical Berry Mojito", Price: 36000, Stock: 75, Description: "Mocktail stroberi dan blueberry bersoda dengan perasan jeruk nipis tanpa alkohol"},

		// Makanan Utama
		{StoreID: 1, CategoryID: categoryMap["Makanan Utama"], Name: "Nasi Goreng Buntut Spesial", Price: 65000, Stock: 40, Description: "Nasi goreng gurih berempah disajikan dengan potongan buntut sapi empuk, telur mata sapi, dan kerupuk emping"},
		{StoreID: 1, CategoryID: categoryMap["Makanan Utama"], Name: "Spaghetti Aglio Olio Smoked Beef", Price: 55000, Stock: 50, Description: "Pasta al dente ditumis dengan minyak zaitun, bawang putih, cabai kering, dan irisan daging asap gurih"},
		{StoreID: 1, CategoryID: categoryMap["Makanan Utama"], Name: "Chicken Katsu Curry Rice", Price: 58000, Stock: 60, Description: "Dada ayam fillet goreng tepung panir renyah disiram saus kari khas Jepang otentik dengan potongan kentang dan wortel"},

		// Roti & Pastry
		{StoreID: 1, CategoryID: categoryMap["Roti & Pastry"], Name: "Butter Croissant Hangat", Price: 22000, Stock: 65, Description: "Pastry klasik khas Prancis dengan tekstur berlapis yang renyah di luar dan lumer rasa mentega di dalam"},
		{StoreID: 1, CategoryID: categoryMap["Roti & Pastry"], Name: "Almond Pain au Chocolat", Price: 28000, Stock: 45, Description: "Pastry cokelat berbalut irisan kacang almond panggang dan taburan gula halus"},
		{StoreID: 1, CategoryID: categoryMap["Roti & Pastry"], Name: "Cinnamon Roll Cream Cheese", Price: 26000, Stock: 55, Description: "Roti gulung lembut beraroma kayu manis pekat dengan olesan krim keju manis gurih di bagian atasnya"},

		// Camilan & Dessert
		{StoreID: 1, CategoryID: categoryMap["Camilan & Dessert"], Name: "Truffle French Fries", Price: 32000, Stock: 130, Description: "Kentang goreng potongan tebal disemprot dengan minyak truffle premium dan ditaburi keju parmesan parut"},
		{StoreID: 1, CategoryID: categoryMap["Camilan & Dessert"], Name: "Platter Pisang Goreng Madu", Price: 27000, Stock: 80, Description: "Pisang raja manis yang digoreng garing berlapis karamelisasi madu murni otentik"},
		{StoreID: 1, CategoryID: categoryMap["Camilan & Dessert"], Name: "Classic Tiramisu Cake", Price: 40000, Stock: 35, Description: "Kue penutup berlapis mascarpone cheese dan ladyfingers yang direndam dalam sirup kopi pekat"},
	}

	_, err = db.NewInsert().Model(&products).Exec(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Products seeded successfully")
	return nil
}
