package response

import (
	"backend-ta/internal/domain"
	"time"
)

type Store struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	LogoURL   string    `json:"logo_url"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewStore(store domain.Store) Store {
	return Store{
		ID:        store.ID,
		Name:      store.Name,
		Address:   store.Address,
		Phone:     store.Phone,
		LogoURL:   store.LogoURL,
		IsActive:  store.IsActive,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}
}

func NewStoreList(stores []domain.Store) []Store {
	res := make([]Store, 0, len(stores))
	for _, store := range stores {
		res = append(res, NewStore(store))
	}
	return res
}
