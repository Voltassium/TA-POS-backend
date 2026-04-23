package response

import "backend-ta/internal/domain"

type SalesData struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type TopSellingProduct struct {
	ProductID    int64   `json:"product_id"`
	ProductName  string  `json:"product_name"`
	CategoryName string  `json:"category_name"`
	Quantity     int     `json:"quantity"`
}

type DashboardStats struct {
	TotalOrders  int64   `json:"total_orders"`
	TotalRevenue float64 `json:"total_revenue"`
}

type DashboardResponse struct {
	Stats        DashboardStats      `json:"stats"`
	SalesChart   []SalesData         `json:"sales_chart"`
	TopProducts  []TopSellingProduct `json:"top_products"`
}

func NewDashboardResponse(stats domain.DashboardStats, sales []domain.SalesData, products []domain.TopSellingProduct) DashboardResponse {
	salesData := make([]SalesData, 0)
	for _, s := range sales {
		salesData = append(salesData, SalesData{
			Date:  s.Date,
			Total: s.Total,
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
			TotalOrders:  stats.TotalOrders,
			TotalRevenue: stats.TotalRevenue,
		},
		SalesChart:  salesData,
		TopProducts: topProducts,
	}
}
