package services

import (
	"backend-ta/app/domain"
	"backend-ta/app/dto/requests"
	"backend-ta/app/dto/response"
	"backend-ta/app/repository"
	"context"
)

type StoreService interface {
	CreateStore(ctx context.Context, req requests.CreateStore) (response.Store, error)
	ListStores(ctx context.Context, page, pageSize int, orderBy, orderDir string) ([]response.Store, int, error)
	UpdateStore(ctx context.Context, id int64, req requests.UpdateStore) (response.Store, error)
	DeleteStore(ctx context.Context, id int64) error
	GetStore(ctx context.Context, id int64) (response.Store, error)
}

type storeService struct {
	repo *repository.PoolRepository
}

func NewStoreService(repo *repository.PoolRepository) StoreService {
	return &storeService{repo: repo}
}

func (s *storeService) CreateStore(ctx context.Context, req requests.CreateStore) (response.Store, error) {
	store := domain.Store{
		Name:    req.Name,
		Address: req.Address,
	}

	err := s.repo.StoreRepository.CreateStore(ctx, &store)
	if err != nil {
		return response.Store{}, err
	}

	return response.NewStore(store), nil
}

func (s *storeService) ListStores(ctx context.Context, page, pageSize int, orderBy, orderDir string) ([]response.Store, int, error) {
	stores, total, err := s.repo.StoreRepository.ListStores(ctx, page, pageSize, orderBy, orderDir)
	if err != nil {
		return nil, 0, err
	}
	return response.NewStoreList(stores), total, nil
}

func (s *storeService) UpdateStore(ctx context.Context, id int64, req requests.UpdateStore) (response.Store, error) {
	store, err := s.repo.StoreRepository.GetStore(ctx, id)
	if err != nil {
		return response.Store{}, err
	}

	if req.Name != "" {
		store.Name = req.Name
	}
	if req.Address != "" {
		store.Address = req.Address
	}

	err = s.repo.StoreRepository.UpdateStore(ctx, &store)
	if err != nil {
		return response.Store{}, err
	}

	return response.NewStore(store), nil
}

func (s *storeService) DeleteStore(ctx context.Context, id int64) error {
	return s.repo.StoreRepository.DeleteStore(ctx, id)
}

func (s *storeService) GetStore(ctx context.Context, id int64) (response.Store, error) {
	store, err := s.repo.StoreRepository.GetStore(ctx, id)
	if err != nil {
		return response.Store{}, err
	}
	return response.NewStore(store), nil
}
