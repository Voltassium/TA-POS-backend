package requests

import (
	"backend-ta/app/constants"
	"backend-ta/app/domain"
)

type UserAuth struct {
	UserID  string             `json:"user_id"`
	StoreID int64              `json:"store_id"`
	Email   string             `json:"email"`
	Role    constants.UserRole `json:"role"`
}

func ToTokenPayload(record domain.User) UserAuth {
	storeID := int64(0)
	if record.StoreID != nil {
		storeID = *record.StoreID
	}
	return UserAuth{
		UserID:  record.ID,
		StoreID: storeID,
		Email:   record.Email,
		Role:    record.Role,
	}
}
