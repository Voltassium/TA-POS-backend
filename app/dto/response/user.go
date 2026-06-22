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
		CreatedAt time.Time          `json:"created_at"`
		UpdatedAt time.Time          `json:"updated_at"`
	}
)

func NewListUser(users []domain.User) []User {
	var res []User
	for _, user := range users {
		res = append(res, User{
			ID:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return res
}

func NewUser(user domain.User) User {
	return User{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

}
