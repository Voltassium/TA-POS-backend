package requests

import (
	"backend-ta/app/constants"
	"backend-ta/app/domain"
	"backend-ta/app/dto"
)

type (
	CreateUser struct {
		Email        string             `json:"email" binding:"required,email"`
		Password     string             `json:"password" binding:"required"`
		Role         constants.UserRole `json:"role" binding:"omitempty,valid_enum"`
		StoreName    string             `json:"store_name" binding:"omitempty"`
		StoreAddress string             `json:"store_address" binding:"omitempty"`
	}

	CreateUserByAdmin struct {
		Email    string             `json:"email" binding:"required,email"`
		Password string             `json:"password" binding:"required"`
		Role     constants.UserRole `json:"role" binding:"required,valid_enum"`
	}

	Login struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	UpdateUser struct {
		Email    string             `json:"email" binding:"omitempty,email"`
		Password string             `json:"password" binding:"omitempty"`
		Role     constants.UserRole `json:"role" binding:"omitempty,valid_enum"`
	}

	ListUser struct {
		dto.PaginationRequest
	}

	RefreshToken struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
)

func (r CreateUser) ToDomain() domain.User {
	return domain.User{
		Email: r.Email,
		Role:  r.Role,
	}
}

func (r CreateUserByAdmin) ToDomain() domain.User {
	return domain.User{
		Email: r.Email,
		Role:  r.Role,
	}
}
