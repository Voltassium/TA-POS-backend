package requests

import (
	"backend-ta/internal/domain"
	"backend-ta/internal/dto"
)

type CreateProduct struct {
	CategoryID  int64   `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"omitempty"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	IsAvailable *bool   `json:"is_available" binding:"omitempty"`
}

type UpdateProduct struct {
	CategoryID  int64   `json:"category_id" binding:"omitempty"`
	Name        string  `json:"name" binding:"omitempty"`
	Description string  `json:"description" binding:"omitempty"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	IsAvailable *bool   `json:"is_available" binding:"omitempty"`
}

type ListProduct struct {
	dto.PaginationRequest
	CategoryID int64 `form:"category_id"`
}

func (r CreateProduct) ToDomain() domain.Product {
	isAvailable := true
	if r.IsAvailable != nil {
		isAvailable = *r.IsAvailable
	}

	return domain.Product{
		CategoryID:  r.CategoryID,
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		IsAvailable: isAvailable,
	}
}
