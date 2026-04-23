package response

import (
	"backend-ta/internal/domain"
	"time"
)

type StockHistory struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	Change    int       `json:"change"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

func NewStockHistory(history domain.StockHistory) StockHistory {
	return StockHistory{
		ID:        history.ID,
		ProductID: history.ProductID,
		Change:    history.Change,
		Reason:    history.Reason,
		CreatedAt: history.CreatedAt,
	}
}

func NewStockHistoryList(histories []domain.StockHistory) []StockHistory {
	res := make([]StockHistory, 0, len(histories))
	for _, history := range histories {
		res = append(res, NewStockHistory(history))
	}
	return res
}
