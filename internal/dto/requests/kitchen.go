package requests

type UpdateServedQty struct {
	ServedQty int `json:"served_qty" binding:"min=0"`
}
