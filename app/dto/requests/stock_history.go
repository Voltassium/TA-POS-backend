package requests

import "backend-ta/app/dto"

type ListStockHistory struct {
	dto.PaginationRequest
	ProductID string `form:"product_id"`
}
