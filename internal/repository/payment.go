package repository

import (
	"backend-ta/internal/domain"
	"backend-ta/pkg/database"
	"context"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, data *domain.Payment) error
	GetPaymentByOrder(ctx context.Context, orderID string) (domain.Payment, error)
}

type paymentRepository struct {
	db *database.Database
}

func NewPaymentRepository(db *database.Database) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreatePayment(ctx context.Context, data *domain.Payment) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	return err
}

func (r *paymentRepository) GetPaymentByOrder(ctx context.Context, orderID string) (domain.Payment, error) {
	var res domain.Payment
	err := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Where("order_id = ?", orderID).
		Scan(ctx)
	return res, err
}
