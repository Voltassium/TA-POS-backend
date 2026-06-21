package domain

type SalesData struct {
	Date  string  `bun:"date"`
	Total float64 `bun:"total"`
}

type TopSellingProduct struct {
	ProductID    string  `bun:"product_id"`
	ProductName  string  `bun:"product_name"`
	CategoryName string  `bun:"category_name"`
	Quantity     int     `bun:"quantity"`
}

type DashboardStats struct {
	TotalOrders   int64   `bun:"total_orders"`
	TotalRevenue  float64 `bun:"total_revenue"`
	TotalProfit   float64 `bun:"total_profit"`
	TotalExpenses float64 `bun:"total_expenses"`
}
