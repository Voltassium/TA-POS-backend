package requests

import "backend-ta/internal/constants"

type CreatePayment struct {
	OrderID       int64                   `json:"order_id" binding:"required"`
	PaymentMethod constants.PaymentMethod `json:"payment_method" binding:"required,valid_enum"`
	AmountPaid    float64                 `json:"amount_paid" binding:"required,gt=0"`
}
