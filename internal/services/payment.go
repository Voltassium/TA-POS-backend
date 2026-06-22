package services

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/domain"
	"backend-ta/internal/dto/requests"
	"backend-ta/internal/dto/response"
	"backend-ta/internal/repository"
	"backend-ta/pkg/database"
	internal_err "backend-ta/pkg/errors"
	"context"
	"database/sql"
	"net/http"

	"github.com/uptrace/bun"
)

type PaymentService interface {
	Process(ctx context.Context, payload requests.CreatePayment) (response.Payment, error)
	GetByOrder(ctx context.Context, orderID string) (response.Payment, error)
}

type paymentService struct {
	orderRepo        repository.OrderRepository
	paymentRepo      repository.PaymentRepository
	productRepo      repository.ProductRepository
	stockHistoryRepo repository.StockHistoryRepository
}

func NewPaymentSrv(
	orderRepo repository.OrderRepository,
	paymentRepo repository.PaymentRepository,
	productRepo repository.ProductRepository,
	stockHistoryRepo repository.StockHistoryRepository,
) PaymentService {
	return &paymentService{
		orderRepo:        orderRepo,
		paymentRepo:      paymentRepo,
		productRepo:      productRepo,
		stockHistoryRepo: stockHistoryRepo,
	}
}

func (s *paymentService) Process(ctx context.Context, payload requests.CreatePayment) (response.Payment, error) {
	var paymentRes response.Payment
	var payment domain.Payment

	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		order, err := s.orderRepo.GetOrder(ctx, payload.OrderID)
		if err != nil {
			return err
		}

		if order.Status != constants.OrderStatusNew {
			return internal_err.NewDefaultError(http.StatusBadRequest, "Payment allowed only for open orders")
		}

		if payload.AmountPaid < order.TotalAmount {
			return internal_err.NewDefaultError(http.StatusBadRequest, "Amount paid is less than total")
		}

		payment = domain.Payment{
			OrderID:       order.ID,
			PaymentMethod: payload.PaymentMethod,
			AmountPaid:    payload.AmountPaid,
		}
		if err := s.paymentRepo.CreatePayment(ctx, &payment); err != nil {
			return err
		}

		if err := s.orderRepo.UpdateOrderStatus(ctx, order.ID, constants.OrderStatusPaid); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return paymentRes, err
	}

	paymentRes = response.NewPayment(payment)
	return paymentRes, nil
}

func (s *paymentService) GetByOrder(ctx context.Context, orderID string) (response.Payment, error) {
	payment, err := s.paymentRepo.GetPaymentByOrder(ctx, orderID)
	if err != nil {
		return response.Payment{}, err
	}
	return response.NewPayment(payment), nil
}
