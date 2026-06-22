package response

import "backend-ta/internal/domain"

type SalesData struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type TopSellingProduct struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	CategoryName string  `json:"category_name"`
	Quantity     int     `json:"quantity"`
}

type DashboardStats struct {
	TotalOrders   int64   `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
	TotalProfit   float64 `json:"total_profit"`
	TotalExpenses float64 `json:"total_expenses"`
}

type FinanceChartData struct {
	Date     string  `json:"date"`
	Revenue  float64 `json:"revenue"`
	Expenses float64 `json:"expenses"`
	Profit   float64 `json:"profit"`
}

type DashboardResponse struct {
	Stats        DashboardStats      `json:"stats"`
	SalesChart   []SalesData         `json:"sales_chart"`
	FinanceChart []FinanceChartData  `json:"finance_chart"`
	TopProducts  []TopSellingProduct `json:"top_products"`
}

func NewDashboardResponse(stats domain.DashboardStats, sales []domain.SalesData, finance []domain.FinanceChartData, products []domain.TopSellingProduct) DashboardResponse {
	salesData := make([]SalesData, 0)
	for _, s := range sales {
		salesData = append(salesData, SalesData{
			Date:  s.Date,
			Total: s.Total,
		})
	}

	financeData := make([]FinanceChartData, 0)
	for _, f := range finance {
		financeData = append(financeData, FinanceChartData{
			Date:     f.Date,
			Revenue:  f.Revenue,
			Expenses: f.Expenses,
			Profit:   f.Profit,
		})
	}

	topProducts := make([]TopSellingProduct, 0)
	for _, p := range products {
		topProducts = append(topProducts, TopSellingProduct{
			ProductID:    p.ProductID,
			ProductName:  p.ProductName,
			CategoryName: p.CategoryName,
			Quantity:     p.Quantity,
		})
	}

	return DashboardResponse{
		Stats: DashboardStats{
			TotalOrders:   stats.TotalOrders,
			TotalRevenue:  stats.TotalRevenue,
			TotalProfit:   stats.TotalProfit,
			TotalExpenses: stats.TotalExpenses,
		},
		SalesChart:   salesData,
		FinanceChart: financeData,
		TopProducts:  topProducts,
	}
}
