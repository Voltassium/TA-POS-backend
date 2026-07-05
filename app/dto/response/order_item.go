package response

import (
	"backend-ta/app/domain"
	"time"
)

type OrderItem struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"order_id"`
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity       int       `json:"quantity"`
	UnitPrice      float64   `json:"unit_price"`
	Subtotal       float64   `json:"subtotal"`
	ServedQty      int       `json:"served_qty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewOrderItem(item domain.OrderItem) OrderItem {
	productName := ""
	if item.Product != nil {
		productName = item.Product.Name
	}

	return OrderItem{
		ID:          item.ID,
		OrderID:     item.OrderID,
		ProductID:   item.ProductID,
		ProductName: productName,
		Quantity:       item.Quantity,
		UnitPrice:      item.UnitPrice,
		Subtotal:       item.Subtotal,
		ServedQty:      item.ServedQty,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func NewOrderItemList(items []domain.OrderItem) []OrderItem {
	res := make([]OrderItem, 0, len(items))
	for _, item := range items {
		res = append(res, NewOrderItem(item))
	}
	return res
}
