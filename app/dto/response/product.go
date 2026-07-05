package response

import (
	"backend-ta/app/domain"
	"time"
)

type Product struct {
	ID           string    `json:"id"`
	CategoryID   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	ProductType  string    `json:"product_type"`
	SKU          *string   `json:"sku"`
	HargaBeli    *float64  `json:"harga_beli"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	IsAvailable  bool      `json:"is_available"`
	Stock        int       `json:"stock"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewProduct(product domain.Product) Product {
	categoryName := ""
	if product.Category != nil {
		categoryName = product.Category.Name
	}

	return Product{
		ID:           product.ID,
		CategoryID:   product.CategoryID,
		CategoryName: categoryName,
		ProductType:  product.ProductType,
		SKU:          product.SKU,
		HargaBeli:    product.HargaBeli,
		Name:         product.Name,
		Price:        product.Price,
		IsAvailable:  product.IsAvailable,
		Stock:        product.Stock,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

func NewProductList(products []domain.Product) []Product {
	res := make([]Product, 0, len(products))
	for _, product := range products {
		res = append(res, NewProduct(product))
	}
	return res
}
