package response

import (
	"backend-ta/internal/domain"
	"time"
)

type Pengeluaran struct {
	ID          string    `json:"id"`
	StoreID     int64     `json:"store_id"`
	Tanggal     string    `json:"tanggal"`
	Category    string    `json:"category"`
	Description *string   `json:"description"`
	Amount      float64   `json:"amount"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewPengeluaran(p domain.Pengeluaran) Pengeluaran {
	return Pengeluaran{
		ID:          p.ID,
		StoreID:     p.StoreID,
		Tanggal:     p.Tanggal.Format("2006-01-02"),
		Category:    p.Category,
		Description: p.Description,
		Amount:      p.Amount,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   p.CreatedAt,
	}
}

func NewPengeluaranList(list []domain.Pengeluaran) []Pengeluaran {
	res := make([]Pengeluaran, 0, len(list))
	for _, item := range list {
		res = append(res, NewPengeluaran(item))
	}
	return res
}
