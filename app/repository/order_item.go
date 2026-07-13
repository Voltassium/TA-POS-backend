package repository

import (
	"backend-ta/app/domain"
	"backend-ta/pkg/database"
	"context"

	"github.com/uptrace/bun"
)

type OrderItemRepository interface {
	CreateItem(ctx context.Context, db bun.IDB, data *domain.OrderItem) error
	DeleteItem(ctx context.Context, db bun.IDB, id string) error
	GetItem(ctx context.Context, db bun.IDB, id string) (domain.OrderItem, error)
	SumSubtotalByOrder(ctx context.Context, db bun.IDB, orderID string) (float64, error)
	UpdateServedQty(ctx context.Context, itemID string, servedQty int) error
	ListItemsByOrder(ctx context.Context, orderID string) ([]domain.OrderItem, error)
}

type orderItemRepository struct {
	db *database.Database
}

func NewOrderItemRepository(db *database.Database) OrderItemRepository {
	return &orderItemRepository{db: db}
}

func (r *orderItemRepository) CreateItem(ctx context.Context, db bun.IDB, data *domain.OrderItem) error {
	_, err := db.NewInsert().Model(data).Returning("id").Exec(ctx)
	return err
}

func (r *orderItemRepository) DeleteItem(ctx context.Context, db bun.IDB, id string) error {
	_, err := db.NewDelete().
		Model((*domain.OrderItem)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *orderItemRepository) GetItem(ctx context.Context, db bun.IDB, id string) (domain.OrderItem, error) {
	var res domain.OrderItem
	err := db.NewSelect().
		Model(&res).
		Where("order_item.id = ?", id).
		Scan(ctx)
	return res, err
}

func (r *orderItemRepository) SumSubtotalByOrder(ctx context.Context, db bun.IDB, orderID string) (float64, error) {
	var total float64
	err := db.NewSelect().
		Model((*domain.OrderItem)(nil)).
		ColumnExpr("COALESCE(SUM(subtotal), 0)").
		Where("order_id = ?", orderID).
		Scan(ctx, &total)
	return total, err
}

func (r *orderItemRepository) UpdateServedQty(ctx context.Context, itemID string, servedQty int) error {
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model((*domain.OrderItem)(nil)).
		Set("served_qty = ?", servedQty).
		Set("updated_at = NOW()").
		Where("id = ?", itemID).
		Exec(ctx)
	return err
}

func (r *orderItemRepository) ListItemsByOrder(ctx context.Context, orderID string) ([]domain.OrderItem, error) {
	var items []domain.OrderItem
	err := r.db.InitQuery(ctx).
		NewSelect().
		Model(&items).
		Where("order_id = ?", orderID).
		Scan(ctx)
	return items, err
}
