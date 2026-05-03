package requests

import "backend-ta/internal/constants"

type AddOrderItem struct {
	ProductID     int64                   `json:"product_id" binding:"required"`
	Quantity      int                     `json:"quantity" binding:"required,min=1"`
	DiscountType  *constants.DiscountType `json:"discount_type" binding:"omitempty,valid_enum"`
	DiscountValue float64                 `json:"discount_value" binding:"omitempty,min=0"`
}
