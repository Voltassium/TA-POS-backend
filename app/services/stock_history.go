package services

import (
	"backend-ta/app/dto"
	"backend-ta/app/dto/requests"
	"backend-ta/app/dto/response"
	"backend-ta/app/repository"
	"context"
)

type StockHistoryService interface {
	List(ctx context.Context, payload requests.ListStockHistory) (dto.PaginationResponse[response.StockHistory], error)
}

type stockHistoryService struct {
	repo repository.StockHistoryRepository
}

func NewStockHistorySrv(repo repository.StockHistoryRepository) StockHistoryService {
	return &stockHistoryService{repo: repo}
}

func (s *stockHistoryService) List(ctx context.Context, payload requests.ListStockHistory) (dto.PaginationResponse[response.StockHistory], error) {
	var paginateRes dto.PaginationResponse[response.StockHistory]
	res, count, err := s.repo.ListStockHistory(ctx, payload)
	if err != nil {
		return paginateRes, err
	}

	paginateRes = dto.NewPaginationResponse(payload.PaginationRequest, count, response.NewStockHistoryList(res))
	return paginateRes, nil
}
