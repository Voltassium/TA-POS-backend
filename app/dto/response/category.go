package response

import (
	"backend-ta/app/domain"
	"time"
)

type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCategory(category domain.Category) Category {
	return Category{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func NewCategoryList(categories []domain.Category) []Category {
	res := make([]Category, 0, len(categories))
	for _, category := range categories {
		res = append(res, NewCategory(category))
	}
	return res
}
