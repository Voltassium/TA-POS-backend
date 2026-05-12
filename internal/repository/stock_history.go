package repository

import (
	"backend-ta/internal/domain"
	"backend-ta/internal/dto/requests"
	"backend-ta/pkg/database"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type StockHistoryRepository interface {
	CreateStockHistory(ctx context.Context, tx bun.Tx, data *domain.StockHistory) error
	ListStockHistory(ctx context.Context, req requests.ListStockHistory) ([]domain.StockHistory, int, error)
}

type stockHistoryRepository struct {
	db *database.Database
}

func NewStockHistoryRepository(db *database.Database) StockHistoryRepository {
	return &stockHistoryRepository{db: db}
}

func (r *stockHistoryRepository) CreateStockHistory(ctx context.Context, tx bun.Tx, data *domain.StockHistory) error {
	_, err := tx.NewInsert().Model(data).Returning("id").Exec(ctx)
	return err
}

func (r *stockHistoryRepository) ListStockHistory(ctx context.Context, req requests.ListStockHistory) ([]domain.StockHistory, int, error) {
	var res []domain.StockHistory
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("Product")

	if req.ProductID != 0 {
		q.Where("stock_history.product_id = ?", req.ProductID)
	}

	q.Limit(req.PageSize).
		Offset(req.CalculateOffset()).
		Order(fmt.Sprintf("%s %s", req.OrderBy, req.OrderDir))

	total, err := q.ScanAndCount(ctx)
	return res, total, err
}
