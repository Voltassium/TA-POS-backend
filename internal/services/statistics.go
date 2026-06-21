package services

import (
	"backend-ta/internal/dto/response"
	"backend-ta/internal/repository"
	"context"
	"time"
)

type StatisticsService interface {
	GetDashboardData(ctx context.Context, timeRange string) (response.DashboardResponse, error)
}

type statisticsService struct {
	statsRepo repository.StatisticsRepository
}

func NewStatisticsSrv(statsRepo repository.StatisticsRepository) StatisticsService {
	return &statisticsService{statsRepo: statsRepo}
}

func (s *statisticsService) GetDashboardData(ctx context.Context, timeRange string) (response.DashboardResponse, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	now := time.Now().In(loc)

	var startDate, endDate time.Time

	switch timeRange {
	case "weekly":
		offset := int(time.Monday - now.Weekday())
		if offset > 0 {
			offset = -6
		}
		startOfWeek := now.AddDate(0, 0, offset)
		startDate = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, loc)
		endOfWeek := startOfWeek.AddDate(0, 0, 6)
		endDate = time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 999999999, loc)
	case "monthly":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		nextMonth := startDate.AddDate(0, 1, 0)
		endDate = nextMonth.Add(-time.Nanosecond)
	case "all":
	    startDate = time.Date(2000, 1, 1, 0, 0, 0, 0, loc)
	    endDate = time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, loc)
	default: // daily
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
	}

	sales, err := s.statsRepo.GetSalesChart(ctx, startDate, endDate)
	if err != nil {
		return response.DashboardResponse{}, err
	}

	topProducts, err := s.statsRepo.GetTopProducts(ctx, 10, startDate, endDate) // Top 10 products
	if err != nil {
		return response.DashboardResponse{}, err
	}

	stats, err := s.statsRepo.GetDashboardStats(ctx, startDate, endDate)
	if err != nil {
		return response.DashboardResponse{}, err
	}

	return response.NewDashboardResponse(stats, sales, topProducts), nil
}
