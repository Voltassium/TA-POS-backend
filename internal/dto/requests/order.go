package requests

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/dto"
)

type CreateOrder struct {
	TableID       *int64                `json:"table_id" binding:"omitempty"`
	CustomerName  *string               `json:"customer_name" binding:"omitempty"`
	DiscountType  *constants.DiscountType `json:"discount_type" binding:"omitempty,valid_enum"`
	DiscountValue float64               `json:"discount_value" binding:"omitempty,min=0"`
	Items         []AddOrderItem        `json:"items" binding:"omitempty,dive"`
}

type UpdateOrderStatus struct {
	Status constants.OrderStatus `json:"status" binding:"required,valid_enum"`
}

type ListOrder struct {
	dto.PaginationRequest
	Status        constants.OrderStatus `form:"status" binding:"omitempty,valid_enum"`
	ExcludeStatus constants.OrderStatus `form:"exclude_status" binding:"omitempty,valid_enum"`
}
