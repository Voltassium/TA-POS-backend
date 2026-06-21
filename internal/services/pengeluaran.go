package services

import (
	"backend-ta/internal/dto"
	"backend-ta/internal/dto/requests"
	"backend-ta/internal/dto/response"
	"backend-ta/internal/repository"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/errors"
	"context"
	"net/http"
	"time"
)

type PengeluaranService interface {
	Create(ctx context.Context, payload requests.CreatePengeluaran) (response.Pengeluaran, error)
	Update(ctx context.Context, id string, payload requests.UpdatePengeluaran) error
	Delete(ctx context.Context, id string) error
	Detail(ctx context.Context, id string) (response.Pengeluaran, error)
	List(ctx context.Context, payload requests.ListPengeluaran) (dto.PaginationResponse[response.Pengeluaran], error)
}

type pengeluaranService struct {
	repo repository.PengeluaranRepository
}

func NewPengeluaranSrv(repo repository.PengeluaranRepository) PengeluaranService {
	return &pengeluaranService{repo: repo}
}

func (s *pengeluaranService) Create(ctx context.Context, payload requests.CreatePengeluaran) (response.Pengeluaran, error) {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	userID := authentication.GetUserDataFromToken(ctx).UserID
	if storeID == 0 || userID == "" {
		return response.Pengeluaran{}, errors.NewDefaultError(http.StatusUnauthorized, "Invalid user or store context")
	}

	pengeluaran, err := payload.ToDomain()
	if err != nil {
		return response.Pengeluaran{}, errors.NewDefaultError(http.StatusBadRequest, "Invalid date format")
	}

	pengeluaran.StoreID = storeID
	pengeluaran.CreatedBy = userID

	if err := s.repo.Create(ctx, &pengeluaran); err != nil {
		return response.Pengeluaran{}, err
	}

	return response.NewPengeluaran(pengeluaran), nil
}

func (s *pengeluaranService) Update(ctx context.Context, id string, payload requests.UpdatePengeluaran) error {
	pengeluaran, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if payload.Tanggal != nil {
		t, err := time.Parse("2006-01-02", *payload.Tanggal)
		if err != nil {
			return errors.NewDefaultError(http.StatusBadRequest, "Invalid date format")
		}
		pengeluaran.Tanggal = t
	}
	if payload.Category != nil {
		pengeluaran.Category = *payload.Category
	}
	if payload.Description != nil {
		pengeluaran.Description = payload.Description
	}
	if payload.Amount != nil {
		pengeluaran.Amount = *payload.Amount
	}

	return s.repo.Update(ctx, &pengeluaran)
}

func (s *pengeluaranService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *pengeluaranService) Detail(ctx context.Context, id string) (response.Pengeluaran, error) {
	pengeluaran, err := s.repo.Get(ctx, id)
	if err != nil {
		return response.Pengeluaran{}, err
	}
	return response.NewPengeluaran(pengeluaran), nil
}

func (s *pengeluaranService) List(ctx context.Context, payload requests.ListPengeluaran) (dto.PaginationResponse[response.Pengeluaran], error) {
	var paginateRes dto.PaginationResponse[response.Pengeluaran]
	res, count, err := s.repo.List(ctx, payload)
	if err != nil {
		return paginateRes, err
	}

	paginateRes = dto.NewPaginationResponse(payload.PaginationRequest, count, response.NewPengeluaranList(res))
	return paginateRes, nil
}
