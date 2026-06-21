package requests

import "backend-ta/internal/dto"

type ListStockHistory struct {
	dto.PaginationRequest
	ProductID string `form:"product_id"`
}
