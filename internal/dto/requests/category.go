package requests

import (
	"backend-ta/internal/domain"
	"backend-ta/internal/dto"
)

type CreateCategory struct {
	Name     string `json:"name" binding:"required"`
}

type UpdateCategory struct {
	Name     string `json:"name" binding:"omitempty"`
}

type ListCategory struct {
	dto.PaginationRequest
}

func (r CreateCategory) ToDomain() domain.Category {
	return domain.Category{
		Name:     r.Name,
	}
}
