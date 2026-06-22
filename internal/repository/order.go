package repository

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/domain"
	"backend-ta/internal/dto/requests"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/database"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, data *domain.Order) error
	CountTodayOrders(ctx context.Context, storeID int64) (int, error)
	UpdateOrderStatus(ctx context.Context, id string, status constants.OrderStatus) error
	UpdateOrderAmounts(ctx context.Context, id string, total float64, discountAmount float64) error
	GetOrder(ctx context.Context, id string) (domain.Order, error)
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

func (r *orderRepository) CountTodayOrders(ctx context.Context, storeID int64) (int, error) {
	count, err := r.db.InitQuery(ctx).
		NewSelect().
		Model((*domain.Order)(nil)).
		Where("store_id = ?", storeID).
		Where("created_at >= CURRENT_DATE").
		Where("created_at < CURRENT_DATE + INTERVAL '1 day'").
		Count(ctx)
	return count, err
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, id string, status constants.OrderStatus) error {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model((*domain.Order)(nil)).
		Set("status = ?", status).
		Set("updated_at = NOW()").
		Where("id = ?", id).
		Where("store_id = ?", storeID).
		Exec(ctx)
	return err
}

func (r *orderRepository) UpdateOrderAmounts(ctx context.Context, id string, total float64, discountAmount float64) error {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model((*domain.Order)(nil)).
		Set("total_amount = ?", total).
		Set("discount_amount = ?", discountAmount).
		Set("updated_at = NOW()").
		Where("id = ?", id).
		Where("store_id = ?", storeID).
		Exec(ctx)
	return err
}

func (r *orderRepository) GetOrder(ctx context.Context, id string) (domain.Order, error) {
	var res domain.Order
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	err := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("OrderItems", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Relation("Product")
		}).
		Relation("Payment").
		Relation("Staff").
		Where("\"order\".id = ?", id).
		Where("\"order\".store_id = ?", storeID).
		Scan(ctx)
	return res, err
}

func (r *orderRepository) ListOrders(ctx context.Context, req requests.ListOrder) ([]domain.Order, int, error) {
	var res []domain.Order
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("Staff").
		Where("\"order\".store_id = ?", storeID)

	if req.Status != "" {
		q.Where("status = ?", req.Status)
	}
	if req.ExcludeStatus != "" {
		q.Where("status != ?", req.ExcludeStatus)
	}
	if req.Search != "" {
		q.Where("order_code ILIKE ?", "%"+req.Search+"%")
	}

	q.Limit(req.PageSize).
		Offset(req.CalculateOffset()).
		Order(fmt.Sprintf("%s %s", req.OrderBy, req.OrderDir))

	total, err := q.ScanAndCount(ctx)
	return res, total, err
}
