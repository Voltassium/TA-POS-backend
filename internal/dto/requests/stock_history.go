package requests

import "backend-ta/internal/dto"

type ListStockHistory struct {
	dto.PaginationRequest
	ProductID int64 `form:"product_id"`
}
