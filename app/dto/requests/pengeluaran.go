package requests

import (
	"backend-ta/app/domain"
	"backend-ta/app/dto"
	"time"
)

type CreatePengeluaran struct {
	Tanggal     string  `json:"tanggal" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Description *string `json:"description" binding:"omitempty"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
}

type UpdatePengeluaran struct {
	Tanggal     *string  `json:"tanggal" binding:"omitempty"`
	Category    *string  `json:"category" binding:"omitempty"`
	Description *string  `json:"description" binding:"omitempty"`
	Amount      *float64 `json:"amount" binding:"omitempty,gt=0"`
}

type ListPengeluaran struct {
	dto.PaginationRequest
	StartDate *string `form:"start_date" binding:"omitempty"`
	EndDate   *string `form:"end_date" binding:"omitempty"`
}

func (r CreatePengeluaran) ToDomain() (domain.Pengeluaran, error) {
	t, err := time.Parse("2006-01-02", r.Tanggal)
	if err != nil {
		return domain.Pengeluaran{}, err
	}
	return domain.Pengeluaran{
		Tanggal:     t,
		Category:    r.Category,
		Description: r.Description,
		Amount:      r.Amount,
	}, nil
}
