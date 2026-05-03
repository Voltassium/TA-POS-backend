package requests

type CreateStore struct {
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address" binding:"omitempty"`
	Phone    string `json:"phone" binding:"omitempty"`
	LogoURL  string `json:"logo_url" binding:"omitempty"`
	IsActive *bool  `json:"is_active" binding:"omitempty"`
}

type UpdateStore struct {
	Name     string `json:"name" binding:"omitempty"`
	Address  string `json:"address" binding:"omitempty"`
	Phone    string `json:"phone" binding:"omitempty"`
	LogoURL  string `json:"logo_url" binding:"omitempty"`
	IsActive *bool  `json:"is_active" binding:"omitempty"`
}
