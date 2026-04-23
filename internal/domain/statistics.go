package domain

type SalesData struct {
	Date  string  `bun:"date"`
	Total float64 `bun:"total"`
}

type TopSellingProduct struct {
	ProductID    int64   `bun:"product_id"`
	ProductName  string  `bun:"product_name"`
	CategoryName string  `bun:"category_name"`
	Quantity     int     `bun:"quantity"`
}

type DashboardStats struct {
	TotalOrders  int64   `bun:"total_orders"`
	TotalRevenue float64 `bun:"total_revenue"`
}
