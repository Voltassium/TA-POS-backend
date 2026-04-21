package repository

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/domain"
	"backend-ta/internal/dto/requests"
	"backend-ta/pkg/database"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, data *domain.Order) error
	UpdateOrderStatus(ctx context.Context, id int64, status constants.OrderStatus) error
	UpdateOrderTotal(ctx context.Context, id int64, total float64) error
	GetOrder(ctx context.Context, id int64) (domain.Order, error)
	ListOrders(ctx context.Context, req requests.ListOrder) ([]domain.Order, int, error)
}

type orderRepository struct {
	db *database.Database
}

func NewOrderRepository(db *database.Database) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, data *domain.Order) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	return err
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, id int64, status constants.OrderStatus) error {
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model((*domain.Order)(nil)).
		Set("status = ?", status).
		Set("updated_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *orderRepository) UpdateOrderTotal(ctx context.Context, id int64, total float64) error {
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model((*domain.Order)(nil)).
		Set("total_amount = ?", total).
		Set("updated_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *orderRepository) GetOrder(ctx context.Context, id int64) (domain.Order, error) {
	var res domain.Order
	err := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("OrderItems", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Relation("Product")
		}).
		Relation("Payment").
		Where("\"order\".id = ?", id).
		Scan(ctx)
	return res, err
}

func (r *orderRepository) ListOrders(ctx context.Context, req requests.ListOrder) ([]domain.Order, int, error) {
	var res []domain.Order
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Limit(req.PageSize).
		Offset(req.CalculateOffset()).
		Order(fmt.Sprintf("%s %s", req.OrderBy, req.OrderDir))

	if req.Status != "" {
		q.Where("status = ?", req.Status)
	}

	total, err := q.ScanAndCount(ctx)
	return res, total, err
}
