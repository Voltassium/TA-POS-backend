package response

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/domain"
	"time"
)

type Payment struct {
	ID            int64                   `json:"id"`
	OrderID       int64                   `json:"order_id"`
	PaymentMethod constants.PaymentMethod `json:"payment_method"`
	AmountPaid    float64                 `json:"amount_paid"`
	Timestamp     time.Time               `json:"timestamp"`
}

func NewPayment(payment domain.Payment) Payment {
	return Payment{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		PaymentMethod: payment.PaymentMethod,
		AmountPaid:    payment.AmountPaid,
		Timestamp:     payment.Timestamp,
	}
}
