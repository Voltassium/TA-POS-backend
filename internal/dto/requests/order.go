package requests

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/dto"
)

type CreateOrder struct {
	TableID int64 `json:"table_id" binding:"required"`
}

type UpdateOrderStatus struct {
	Status constants.OrderStatus `json:"status" binding:"required,valid_enum"`
}

type ListOrder struct {
	dto.PaginationRequest
	Status constants.OrderStatus `form:"status" binding:"omitempty,valid_enum"`
}
