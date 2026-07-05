package services

import (
	"backend-ta/app/domain"
	"backend-ta/app/dto/response"
	"backend-ta/app/repository"
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
		startDate = time.Date(2026, 1, 1, 0, 0, 0, 0, loc)
		endDate = time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, loc)
	default:
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
	}

	var dateFormat string
	if timeRange == "all" {
		dateFormat = "YYYY-MM"
	} else {
		dateFormat = "YYYY-MM-DD"
	}

	sales, err := s.statsRepo.GetSalesChart(ctx, startDate, endDate, dateFormat)
	if err != nil {
		return response.DashboardResponse{}, err
	}

	expenses, err := s.statsRepo.GetExpensesChart(ctx, startDate, endDate, dateFormat)
	if err != nil {
		return response.DashboardResponse{}, err
	}

	cogs, err := s.statsRepo.GetCogsChart(ctx, startDate, endDate, dateFormat)
	if err != nil {
		return response.DashboardResponse{}, err
	}

	financeMap := make(map[string]*domain.FinanceChartData)

	for _, s := range sales {
		if _, ok := financeMap[s.Date]; !ok {
			financeMap[s.Date] = &domain.FinanceChartData{Date: s.Date}
		}
		financeMap[s.Date].Revenue = s.Total
	}

	for _, e := range expenses {
		if _, ok := financeMap[e.Date]; !ok {
			financeMap[e.Date] = &domain.FinanceChartData{Date: e.Date}
		}
		financeMap[e.Date].Expenses = e.Total
	}

	for _, c := range cogs {
		if _, ok := financeMap[c.Date]; !ok {
			financeMap[c.Date] = &domain.FinanceChartData{Date: c.Date}
		}

		financeMap[c.Date].Expenses += c.Total
	}

	var financeChart []domain.FinanceChartData

	if timeRange == "all" {
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 1, 0) {
			dateStr := d.Format("2006-01")
			if f, ok := financeMap[dateStr]; ok {
				f.Profit += f.Revenue - f.Expenses
				financeChart = append(financeChart, *f)
			} else {
				financeChart = append(financeChart, domain.FinanceChartData{Date: dateStr})
			}
		}
	} else {
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			dateStr := d.Format("2006-01-02")
			if f, ok := financeMap[dateStr]; ok {
				f.Profit += f.Revenue - f.Expenses
				financeChart = append(financeChart, *f)
			} else {
				financeChart = append(financeChart, domain.FinanceChartData{Date: dateStr})
			}
		}
	}

	topProducts, err := s.statsRepo.GetTopProducts(ctx, 10, startDate, endDate)
	if err != nil {
		return response.DashboardResponse{}, err
	}

	stats, err := s.statsRepo.GetDashboardStats(ctx, startDate, endDate)
	if err != nil {
		return response.DashboardResponse{}, err
	}

	return response.NewDashboardResponse(stats, sales, financeChart, topProducts), nil
}
