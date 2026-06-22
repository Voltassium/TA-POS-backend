package constants

type OrderStatus string

type PaymentMethod string

const (
	OrderStatusNew       OrderStatus = "New"
	OrderStatusPaid      OrderStatus = "Paid"
	OrderStatusCancelled OrderStatus = "Cancelled"
	OrderStatusCompleted OrderStatus = "Completed"
)

const (
	PaymentMethodCash          PaymentMethod = "Cash"
	PaymentMethodCard          PaymentMethod = "Card"
	PaymentMethodDigitalWallet PaymentMethod = "Digital Wallet"
)

func (receiver OrderStatus) IsValidEnum() bool {
	switch receiver {
	case OrderStatusNew, OrderStatusPaid, OrderStatusCancelled, OrderStatusCompleted:
		return true
	default:
		return false
	}
}

func (receiver PaymentMethod) IsValidEnum() bool {
	switch receiver {
	case PaymentMethodCash, PaymentMethodCard, PaymentMethodDigitalWallet:
		return true
	default:
		return false
	}
}
