package response

import (
	"backend-ta/internal/domain"
	"time"
)

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCategory(category domain.Category) Category {
	return Category{
		ID:        category.ID,
		Name:      category.Name,
		ImageURL:  category.ImageURL,
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
