package requests

type CreateStore struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"omitempty"`
}

type UpdateStore struct {
	Name    string `json:"name" binding:"omitempty"`
	Address string `json:"address" binding:"omitempty"`
}
