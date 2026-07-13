package response

import (
	"backend-ta/app/constants"
	"backend-ta/app/domain"
	"time"
)

type (
	User struct {
		ID        string             `json:"id"`
		Email     string             `json:"email"`
		Role      constants.UserRole `json:"role"`
		StoreName string             `json:"store_name,omitempty"`
		CreatedAt time.Time          `json:"created_at"`
		UpdatedAt time.Time          `json:"updated_at"`
	}
)

func NewListUser(users []domain.User) []User {
	var res []User
	for _, user := range users {
		storeName := ""
		if user.Store != nil {
			storeName = user.Store.Name
		}
		res = append(res, User{
			ID:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			StoreName: storeName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return res
}

func NewUser(user domain.User) User {
	storeName := ""
	if user.Store != nil {
		storeName = user.Store.Name
	}
	return User{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		StoreName: storeName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

}
