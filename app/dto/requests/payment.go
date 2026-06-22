package requests

import "backend-ta/app/constants"

type CreatePayment struct {
	OrderID       string                  `json:"order_id" binding:"required"`
	PaymentMethod constants.PaymentMethod `json:"payment_method" binding:"required,valid_enum"`
	AmountPaid    float64                 `json:"amount_paid" binding:"required,gt=0"`
}
