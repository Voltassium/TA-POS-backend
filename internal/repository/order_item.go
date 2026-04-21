package repository

import (
	"backend-ta/internal/domain"
	"backend-ta/pkg/database"
	"context"
)

type OrderItemRepository interface {
	CreateItem(ctx context.Context, data *domain.OrderItem) error
	DeleteItem(ctx context.Context, id int64) error
	GetItem(ctx context.Context, id int64) (domain.OrderItem, error)
	SumSubtotalByOrder(ctx context.Context, orderID int64) (float64, error)
}

type orderItemRepository struct {
	db *database.Database
}

func NewOrderItemRepository(db *database.Database) OrderItemRepository {
	return &orderItemRepository{db: db}
}

func (r *orderItemRepository) CreateItem(ctx context.Context, data *domain.OrderItem) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	return err
}

func (r *orderItemRepository) DeleteItem(ctx context.Context, id int64) error {
	_, err := r.db.InitQuery(ctx).
		NewDelete().
		Model((*domain.OrderItem)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *orderItemRepository) GetItem(ctx context.Context, id int64) (domain.OrderItem, error) {
	var res domain.OrderItem
	err := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Where("order_item.id = ?", id).
		Scan(ctx)
	return res, err
}

func (r *orderItemRepository) SumSubtotalByOrder(ctx context.Context, orderID int64) (float64, error) {
	var total float64
	err := r.db.InitQuery(ctx).
		NewSelect().
		Model((*domain.OrderItem)(nil)).
		ColumnExpr("COALESCE(SUM(subtotal), 0)").
		Where("order_id = ?", orderID).
		Scan(ctx, &total)
	return total, err
}
