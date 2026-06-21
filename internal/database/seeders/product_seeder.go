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

	categoryMap := make(map[string]string)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
	}

	hb := func(v float64) *float64 { return &v }
	sp := func(v string) *string { return &v }

	products := []domain.Product{
		// Kopi & Espresso (Olahan — dibuat sendiri)
		{StoreID: 1, CategoryID: categoryMap["Kopi & Espresso"], ProductType: "Olahan", SKU: sp("KOPI-001"), Name: "Kopi Tubruk Robusta Gayo", Price: 18000, Stock: 150, Description: "Kopi hitam tradisional tubruk menggunakan biji kopi Robusta Gayo pilihan dengan aroma kuat dan body tebal"},
		{StoreID: 1, CategoryID: categoryMap["Kopi & Espresso"], ProductType: "Olahan", SKU: sp("KOPI-002"), Name: "Es Kopi Susu Gula Aren", Price: 22000, Stock: 200, Description: "Double shot espresso blend, susu segar, dan gula aren cair premium khas Nusantara"},
		{StoreID: 1, CategoryID: categoryMap["Kopi & Espresso"], ProductType: "Olahan", SKU: sp("KOPI-003"), Name: "Es Kopi Susu Pandan", Price: 24000, Stock: 120, Description: "Kopi susu espresso creamy dengan sirup pandan wangi alami yang segar dan aromatik"},
		{StoreID: 1, CategoryID: categoryMap["Kopi & Espresso"], ProductType: "Olahan", SKU: sp("KOPI-004"), Name: "Kopi Tarik Khas Aceh", Price: 20000, Stock: 90, Description: "Perpaduan kopi robusta pekat dan susu kental manis yang ditarik secara tradisional hingga berbusa tebal"},
		{StoreID: 1, CategoryID: categoryMap["Kopi & Espresso"], ProductType: "Olahan", SKU: sp("KOPI-005"), Name: "Es Kopi Hitam Kerinci", Price: 18000, Stock: 100, Description: "Kopi saring dingin menggunakan single origin Kerinci dengan rasa asam manis alami buah"},

		// Minuman Segar (Olahan — dibuat sendiri)
		{StoreID: 1, CategoryID: categoryMap["Minuman Segar"], ProductType: "Olahan", SKU: sp("SEGAR-001"), Name: "Es Cendol Durian Klasik", Price: 28000, Stock: 110, Description: "Minuman santan gurih dengan cendol pandan kenyal, gula merah sisir, dan toping daging buah durian asli"},
		{StoreID: 1, CategoryID: categoryMap["Minuman Segar"], ProductType: "Olahan", SKU: sp("SEGAR-002"), Name: "Es Doger Spesial", Price: 25000, Stock: 85, Description: "Es serut kelapa muda merah muda dengan ketan hitam, pacar cina, roti tawar, dan siraman susu kental manis"},
		{StoreID: 1, CategoryID: categoryMap["Minuman Segar"], ProductType: "Olahan", SKU: sp("SEGAR-003"), Name: "Es Kelapa Muda Jeruk", Price: 22000, Stock: 75, Description: "Air kelapa muda segar dengan serutan kelapa muda, perasan jeruk peras murni, dan es batu"},
		{StoreID: 1, CategoryID: categoryMap["Minuman Segar"], ProductType: "Olahan", SKU: sp("SEGAR-004"), Name: "Es Kunyit Asam Dingin", Price: 15000, Stock: 100, Description: "Jamu kunyit asam segar tradisional yang disajikan dingin, berkhasiat menyegarkan tubuh"},
		{StoreID: 1, CategoryID: categoryMap["Minuman Segar"], ProductType: "Olahan", SKU: sp("SEGAR-005"), Name: "Es Selasih Lemon Selera", Price: 18000, Stock: 120, Description: "Minuman perasan lemon segar dipadu biji selasih dan sirup gula tebu murni"},

		// Makanan Utama (Olahan — dimasak sendiri)
		{StoreID: 1, CategoryID: categoryMap["Makanan Utama"], ProductType: "Olahan", SKU: sp("MAKAN-001"), Name: "Nasi Goreng Buntut Spesial", Price: 65000, Stock: 40, Description: "Nasi goreng gurih berempah disajikan dengan potongan buntut sapi empuk, telur mata sapi, emping, dan acar"},
		{StoreID: 1, CategoryID: categoryMap["Makanan Utama"], ProductType: "Olahan", SKU: sp("MAKAN-002"), Name: "Nasi Ayam Geprek Sambal Korek", Price: 28000, Stock: 120, Description: "Nasi putih hangat dengan dada ayam goreng tepung renyah yang dimemarkan bersama sambal korek bawang pedas mantap"},
		{StoreID: 1, CategoryID: categoryMap["Makanan Utama"], ProductType: "Olahan", SKU: sp("MAKAN-003"), Name: "Mie Goreng Jawa Nyemek", Price: 25000, Stock: 90, Description: "Mie kuning basah ditumis dengan bumbu kemiri, kol, sawi, potongan ayam, bakso, telur, disajikan sedikit berkuah"},
		{StoreID: 1, CategoryID: categoryMap["Makanan Utama"], ProductType: "Olahan", SKU: sp("MAKAN-004"), Name: "Sate Ayam Madura (10 Tusuk)", Price: 30000, Stock: 70, Description: "Sate daging ayam empuk dipanggang arang, disiram bumbu kacang gurih kental, kecap manis, dan irisan bawang merah"},
		{StoreID: 1, CategoryID: categoryMap["Makanan Utama"], ProductType: "Olahan", SKU: sp("MAKAN-005"), Name: "Rendang Sapi Minang Plate", Price: 38000, Stock: 50, Description: "Nasi hangat dengan rendang daging sapi empuk bumbu hitam otentik Padang, daun singkong rebus, dan sambal ijo"},

		// Roti & Pastry (Kulakan — dibeli dari supplier)
		{StoreID: 1, CategoryID: categoryMap["Roti & Pastry"], ProductType: "Kulakan", HargaBeli: hb(12000), SKU: sp("ROTI-001"), Name: "Roti Bakar Bandung Cokelat Keju", Price: 22000, Stock: 65, Description: "Roti tawar tebal bakar mentega dengan isi cokelat meses melimpah dan parutan keju cheddar gurih"},
		{StoreID: 1, CategoryID: categoryMap["Roti & Pastry"], ProductType: "Kulakan", HargaBeli: hb(10000), SKU: sp("ROTI-002"), Name: "Kue Pancong Lumer Keju", Price: 18000, Stock: 80, Description: "Kue pancong tradisional gurih kelapa parut disajikan setengah matang lumer dengan parutan keju cheddar"},
		{StoreID: 1, CategoryID: categoryMap["Roti & Pastry"], ProductType: "Kulakan", HargaBeli: hb(11000), SKU: sp("ROTI-003"), Name: "Pisang Goreng Pasir Keju Aren", Price: 20000, Stock: 95, Description: "Pisang kepok goreng tepung roti krispi, disiram gula aren premium dan taburan keju cheddar parut"},
		{StoreID: 1, CategoryID: categoryMap["Roti & Pastry"], ProductType: "Kulakan", HargaBeli: hb(8000), SKU: sp("ROTI-004"), Name: "Roti Bun Mentega Kopi", Price: 15000, Stock: 110, Description: "Roti manis lembut beraroma kopi mentega yang dipanggang segar dengan kulit luar krispi ala kedai kopi modern"},

		// Camilan & Dessert (Kulakan — dibeli dari supplier)
		{StoreID: 1, CategoryID: categoryMap["Camilan & Dessert"], ProductType: "Kulakan", HargaBeli: hb(7000), SKU: sp("MILAN-001"), Name: "Cireng Renyah Rujak Pedas", Price: 15000, Stock: 130, Description: "Camilan tepung kanji goreng renyah kenyal khas Sunda dengan cocolan bumbu rujak gula merah pedas manis"},
		{StoreID: 1, CategoryID: categoryMap["Camilan & Dessert"], ProductType: "Kulakan", HargaBeli: hb(10000), SKU: sp("MILAN-002"), Name: "Tahu Walik Banyuwangi", Price: 18000, Stock: 80, Description: "Tahu goreng isi adonan bakso ayam gurih dibalik hingga kulit tahu krispi, disajikan dengan cabai rawit hijau"},
		{StoreID: 1, CategoryID: categoryMap["Camilan & Dessert"], ProductType: "Kulakan", HargaBeli: hb(9000), SKU: sp("MILAN-003"), Name: "Singkong Goreng Keju Merekah", Price: 17000, Stock: 90, Description: "Singkong gurih empuk yang digoreng hingga merekah, ditaburi keju parut gurih khas kedai"},
		{StoreID: 1, CategoryID: categoryMap["Camilan & Dessert"], ProductType: "Kulakan", HargaBeli: hb(5000), SKU: sp("MILAN-004"), Name: "Jasuke Jagung Susu Keju", Price: 12000, Stock: 100, Description: "Jagung manis pipil hangat dicampur margarin, disiram susu kental manis manis gurih, dan parutan keju melimpah"},
	}

	_, err = db.NewInsert().Model(&products).Exec(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Products seeded successfully")
	return nil
}
