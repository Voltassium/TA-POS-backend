package requests

type AddOrderItem struct {
	ProductID     string                  `json:"product_id" binding:"required,uuid"`
	Quantity      int                     `json:"quantity" binding:"required,min=1"`
}
