package response

import (
	"backend-ta/internal/domain"
	"time"
)

type StockHistory struct {
	ID          int64     `json:"id"`
	ProductID   int64     `json:"product_id"`
	ProductName string    `json:"product_name"`
	Change      int       `json:"change"`
	Reason      string    `json:"reason"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewStockHistory(history domain.StockHistory) StockHistory {
	var productName string
	if history.Product != nil {
		productName = history.Product.Name
	}

	return StockHistory{
		ID:          history.ID,
		ProductID:   history.ProductID,
		ProductName: productName,
		Change:      history.Change,
		Reason:      history.Reason,
		CreatedAt:   history.CreatedAt,
	}
}

func NewStockHistoryList(histories []domain.StockHistory) []StockHistory {
	res := make([]StockHistory, 0, len(histories))
	for _, history := range histories {
		res = append(res, NewStockHistory(history))
	}
	return res
}
