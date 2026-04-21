package response

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/domain"
	"time"
)

type Order struct {
	ID          int64                 `json:"id"`
	TableID     int64                 `json:"table_id"`
	StaffID     int64                 `json:"staff_id"`
	TotalAmount float64               `json:"total_amount"`
	Status      constants.OrderStatus `json:"status"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type OrderDetail struct {
	Order
	Items   []OrderItem `json:"items"`
	Payment *Payment    `json:"payment"`
}

func NewOrder(order domain.Order) Order {
	return Order{
		ID:          order.ID,
		TableID:     order.TableID,
		StaffID:     order.StaffID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

func NewOrderList(orders []domain.Order) []Order {
	res := make([]Order, 0, len(orders))
	for _, order := range orders {
		res = append(res, NewOrder(order))
	}
	return res
}

func NewOrderDetail(order domain.Order) OrderDetail {
	var payment *Payment
	if order.Payment != nil {
		converted := NewPayment(*order.Payment)
		payment = &converted
	}

	return OrderDetail{
		Order:   NewOrder(order),
		Items:   NewOrderItemList(order.OrderItems),
		Payment: payment,
	}
}
