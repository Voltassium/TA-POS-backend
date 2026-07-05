package requests

import (
	"backend-ta/app/domain"
	"backend-ta/app/dto"
)

type CreateProduct struct {
	CategoryID  string   `json:"category_id" binding:"required,uuid"`
	ProductType string   `json:"product_type" binding:"required,oneof=Kulakan Olahan"`
	SKU         *string  `json:"sku" binding:"omitempty"`
	HargaBeli   *float64 `json:"harga_beli" binding:"omitempty,min=0"`
	Name        string   `json:"name" binding:"required"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	IsAvailable *bool    `json:"is_available" binding:"omitempty"`
	Stock       int      `json:"stock" binding:"omitempty,gte=0"`
}

type UpdateProduct struct {
	CategoryID  string   `json:"category_id" binding:"omitempty,uuid"`
	ProductType string   `json:"product_type" binding:"omitempty,oneof=Kulakan Olahan"`
	SKU         *string  `json:"sku" binding:"omitempty"`
	HargaBeli   *float64 `json:"harga_beli" binding:"omitempty,min=0"`
	Name        string   `json:"name" binding:"omitempty"`
	Price       float64  `json:"price" binding:"omitempty,gt=0"`
	IsAvailable *bool    `json:"is_available" binding:"omitempty"`
	Stock       *int     `json:"stock" binding:"omitempty,gte=0"`
}

type ListProduct struct {
	dto.PaginationRequest
	CategoryID  string `form:"category_id"`
	ProductType string `form:"product_type" binding:"omitempty,oneof=Kulakan Olahan"`
}

func (r CreateProduct) ToDomain() domain.Product {
	isAvailable := true
	if r.IsAvailable != nil {
		isAvailable = *r.IsAvailable
	}

	var hargaBeli *float64
	if r.ProductType == "Kulakan" {
		hargaBeli = r.HargaBeli
	}

	return domain.Product{
		CategoryID:  r.CategoryID,
		ProductType: r.ProductType,
		SKU:         r.SKU,
		HargaBeli:   hargaBeli,
		Name:        r.Name,
		Price:       r.Price,
		IsAvailable: isAvailable,
		Stock:       r.Stock,
	}
}
