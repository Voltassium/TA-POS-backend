package repository

import (
	"backend-ta/internal/domain"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/database"
	"context"
	"time"
)

type StatisticsRepository interface {
	GetSalesChart(ctx context.Context, startDate, endDate time.Time) ([]domain.SalesData, error)
	GetTopProducts(ctx context.Context, limit int, startDate, endDate time.Time) ([]domain.TopSellingProduct, error)
	GetDashboardStats(ctx context.Context, startDate, endDate time.Time) (domain.DashboardStats, error)
}

type statisticsRepository struct{
	db *database.Database
}

func NewStatisticsRepository(db *database.Database) StatisticsRepository {
	return &statisticsRepository{db: db}
}

func (r *statisticsRepository) GetSalesChart(ctx context.Context, startDate, endDate time.Time) ([]domain.SalesData, error) {
	var sales []domain.SalesData

	err := r.db.DB.NewSelect().
		TableExpr("orders").
		ColumnExpr("TO_CHAR(created_at, 'YYYY-MM-DD') AS date").
		ColumnExpr("SUM(total_amount) AS total").
		Where("status = ?", "Paid").
		Where("store_id = ?", authentication.GetUserDataFromToken(ctx).StoreID).
		Where("created_at >= ?", startDate).
		Where("created_at <= ?", endDate).
		GroupExpr("TO_CHAR(created_at, 'YYYY-MM-DD')").
		OrderExpr("date ASC").
		Scan(ctx, &sales)

	return sales, err
}

func (r *statisticsRepository) GetTopProducts(ctx context.Context, limit int, startDate, endDate time.Time) ([]domain.TopSellingProduct, error) {
	var products []domain.TopSellingProduct

	err := r.db.DB.NewSelect().
		TableExpr("order_items AS oi").
		ColumnExpr("oi.product_id").
		ColumnExpr("p.name AS product_name").
		ColumnExpr("c.name AS category_name").
		ColumnExpr("SUM(oi.quantity) AS quantity").
		Join("JOIN orders o ON oi.order_id = o.id").
		Join("JOIN products p ON oi.product_id = p.id").
		Join("LEFT JOIN categories c ON p.category_id = c.id").
		Where("o.status = ?", "Paid").
		Where("o.store_id = ?", authentication.GetUserDataFromToken(ctx).StoreID).
		Where("o.created_at >= ?", startDate).
		Where("o.created_at <= ?", endDate).
		GroupExpr("oi.product_id, p.name, c.name").
		OrderExpr("quantity DESC").
		Limit(limit).
		Scan(ctx, &products)

	return products, err
}

func (r *statisticsRepository) GetDashboardStats(ctx context.Context, startDate, endDate time.Time) (domain.DashboardStats, error) {
	var stats domain.DashboardStats

	err := r.db.DB.NewSelect().
		TableExpr("orders").
		ColumnExpr("COUNT(*) AS total_orders").
		ColumnExpr("COALESCE(SUM(total_amount), 0) AS total_revenue").
		Where("status = ?", "Paid").
		Where("store_id = ?", authentication.GetUserDataFromToken(ctx).StoreID).
		Where("created_at >= ?", startDate).
		Where("created_at <= ?", endDate).
		Scan(ctx, &stats)
	if err != nil {
		return stats, err
	}

	// Calculate Total Expenses
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	err = r.db.DB.NewSelect().
		TableExpr("pengeluaran").
		ColumnExpr("COALESCE(SUM(amount), 0)").
		Where("store_id = ?", storeID).
		Where("tanggal >= ?", startDate.Format("2006-01-02")).
		Where("tanggal <= ?", endDate.Format("2006-01-02")).
		Scan(ctx, &stats.TotalExpenses)
	if err != nil {
		return stats, err
	}

	// Calculate Total Product Kulakan Purchase Costs
	var totalPurchaseCosts float64
	err = r.db.DB.NewSelect().
		TableExpr("order_items AS oi").
		ColumnExpr("COALESCE(SUM(oi.quantity * p.harga_beli), 0)").
		Join("JOIN orders o ON oi.order_id = o.id").
		Join("JOIN products p ON oi.product_id = p.id").
		Where("o.status = ?", "Paid").
		Where("o.store_id = ?", storeID).
		Where("o.created_at >= ?", startDate).
		Where("o.created_at <= ?", endDate).
		Where("p.harga_beli > 0"). // Only count products that have purchase cost
		Scan(ctx, &totalPurchaseCosts)
	if err != nil {
		return stats, err
	}

	// Calculate Profit
	stats.TotalProfit = stats.TotalRevenue - (totalPurchaseCosts + stats.TotalExpenses)

	return stats, nil
}
