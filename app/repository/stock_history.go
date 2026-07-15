package repository

import (
	"backend-ta/app/domain"
	"backend-ta/app/dto/requests"
	"backend-ta/pkg/database"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type StockHistoryRepository interface {
	CreateStockHistory(ctx context.Context, db bun.IDB, data *domain.StockHistory) error
	ListStockHistory(ctx context.Context, req requests.ListStockHistory) ([]domain.StockHistory, int, error)
}

type stockHistoryRepository struct {
	db *database.Database
}

func NewStockHistoryRepository(db *database.Database) StockHistoryRepository {
	return &stockHistoryRepository{db: db}
}

// CreateStockHistory menyimpan riwayat perubahan stok dalam transaksi aktif.
// Membaca stok terkini dari DB (melalui db yang sama agar dalam transaksi yang sama)
// untuk menghitung InitialStock dan FinalStock secara akurat.
func (r *stockHistoryRepository) CreateStockHistory(ctx context.Context, db bun.IDB, data *domain.StockHistory) error {
	var product domain.Product
	err := db.NewSelect().
		Model(&product).
		Column("stock").
		Where("id = ?", data.ProductID).
		Scan(ctx)
	if err != nil {
		return err
	}

	data.FinalStock = product.Stock
	data.InitialStock = product.Stock - data.Change

	_, err = db.NewInsert().Model(data).Returning("id").Exec(ctx)
	return err
}

func (r *stockHistoryRepository) ListStockHistory(ctx context.Context, req requests.ListStockHistory) ([]domain.StockHistory, int, error) {
	var res []domain.StockHistory
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("Product")

	if req.ProductID != "" {
		q.Where("stock_history.product_id = ?", req.ProductID)
	}
	if req.Search != "" {
		q.Where("EXISTS (SELECT 1 FROM products p WHERE p.id = stock_history.product_id AND p.name ILIKE ?)", "%"+req.Search+"%")
	}

	orderBy := req.OrderBy
	if orderBy == "updated_at" || orderBy == "created_at" || orderBy == "" {
		orderBy = "stock_history.created_at"
	} else {
		// Secure prefixing for standard fields
		orderBy = "stock_history." + orderBy
	}

	q.Limit(req.PageSize).
		Offset(req.CalculateOffset()).
		Order(fmt.Sprintf("%s %s", orderBy, req.OrderDir))

	total, err := q.ScanAndCount(ctx)
	return res, total, err
}
