package services

import (
	"backend-ta/internal/dto/response"
	"backend-ta/internal/repository"
	"context"
	"time"
)

type StatisticsService interface {
	GetDashboardData(ctx context.Context) (response.DashboardResponse, error)
}

type statisticsService struct {
	statsRepo repository.StatisticsRepository
}

func NewStatisticsSrv(statsRepo repository.StatisticsRepository) StatisticsService {
	return &statisticsService{statsRepo: statsRepo}
}

func (s *statisticsService) GetDashboardData(ctx context.Context) (response.DashboardResponse, error) {
	// Let's get data for the last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	sales, err := s.statsRepo.GetSalesChart(ctx, startDate, endDate)
	if err != nil {
		return response.DashboardResponse{}, err
	}

	topProducts, err := s.statsRepo.GetTopProducts(ctx, 10) // Top 10 products
	if err != nil {
		return response.DashboardResponse{}, err
	}

	stats, err := s.statsRepo.GetDashboardStats(ctx)
	if err != nil {
		return response.DashboardResponse{}, err
	}

	return response.NewDashboardResponse(stats, sales, topProducts), nil
}
