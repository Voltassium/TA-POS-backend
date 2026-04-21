package requests

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/domain"
)

type UserAuth struct {
	UserID int64              `json:"user_id"`
	Email  string             `json:"email"`
	Role   constants.UserRole `json:"role"`
}

func ToTokenPayload(record domain.User) UserAuth {
	return UserAuth{
		UserID: record.ID,
		Email:  record.Email,
		Role:   record.Role,
	}
}
