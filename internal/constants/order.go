package constants

type OrderStatus string

type PaymentMethod string

const (
	OrderStatusOpen      OrderStatus = "Open"
	OrderStatusPaid      OrderStatus = "Paid"
	OrderStatusCancelled OrderStatus = "Cancelled"
)

const (
	PaymentMethodCash          PaymentMethod = "Cash"
	PaymentMethodCard          PaymentMethod = "Card"
	PaymentMethodDigitalWallet PaymentMethod = "Digital Wallet"
)

func (receiver OrderStatus) IsValidEnum() bool {
	switch receiver {
	case OrderStatusOpen, OrderStatusPaid, OrderStatusCancelled:
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
