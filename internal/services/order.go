package services

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/domain"
	"backend-ta/internal/dto"
	"backend-ta/internal/dto/requests"
	"backend-ta/internal/dto/response"
	"backend-ta/internal/repository"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/database"
	internal_err "backend-ta/pkg/errors"
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/uptrace/bun"
)

type OrderService interface {
	Create(ctx context.Context, payload requests.CreateOrder) (response.OrderDetail, error)
	List(ctx context.Context, payload requests.ListOrder) (dto.PaginationResponse[response.Order], error)
	Detail(ctx context.Context, id int64) (response.OrderDetail, error)
	UpdateStatus(ctx context.Context, id int64, payload requests.UpdateOrderStatus) error
	Cancel(ctx context.Context, id int64) error
	AddItem(ctx context.Context, orderID int64, payload requests.AddOrderItem) (response.OrderDetail, error)
	RemoveItem(ctx context.Context, orderID int64, itemID int64) (response.OrderDetail, error)
}

type orderService struct {
	orderRepo        repository.OrderRepository
	orderItemRepo    repository.OrderItemRepository
	productRepo      repository.ProductRepository
	stockHistoryRepo repository.StockHistoryRepository
}

func NewOrderSrv(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	productRepo repository.ProductRepository,
	stockHistoryRepo repository.StockHistoryRepository,
) OrderService {
	return &orderService{
		orderRepo:        orderRepo,
		orderItemRepo:    orderItemRepo,
		productRepo:      productRepo,
		stockHistoryRepo: stockHistoryRepo,
	}
}

func (s *orderService) Create(ctx context.Context, payload requests.CreateOrder) (response.OrderDetail, error) {
	staffID := authentication.GetUserDataFromToken(ctx).UserID
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	if staffID == 0 || storeID == 0 {
		return response.OrderDetail{}, internal_err.NewDefaultError(http.StatusUnauthorized, "Invalid user or store")
	}

	var orderID int64

	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		order := domainOrderFromCreate(payload, storeID, staffID)
		if err := s.orderRepo.CreateOrder(ctx, &order); err != nil {
			return err
		}
		orderID = order.ID

		if len(payload.Items) > 0 {
			var totalAmount float64
			for _, item := range payload.Items {
				product, err := s.productRepo.GetProduct(ctx, item.ProductID)
				if err != nil {
					return err
				}

				if product.Stock < item.Quantity {
					return internal_err.NewDefaultError(http.StatusBadRequest, fmt.Sprintf("Not enough stock for product: %s", product.Name))
				}

				orderItem := domainOrderItemFromCreate(order.ID, product.Price, item)
				if err := s.orderItemRepo.CreateItem(ctx, &orderItem); err != nil {
					return err
				}
				totalAmount += orderItem.Subtotal

				if err := s.productRepo.UpdateStock(ctx, tx, item.ProductID, -item.Quantity); err != nil {
					return err
				}
				history := domain.StockHistory{
					ProductID: item.ProductID,
					Change:    -item.Quantity,
					Reason:    fmt.Sprintf("Order #%d Created", order.ID),
				}
				if err := s.stockHistoryRepo.CreateStockHistory(ctx, tx, &history); err != nil {
					return err
				}
			}

			finalTotal, discountAmount := computeOrderAmounts(totalAmount, &order)
			if err := s.orderRepo.UpdateOrderAmounts(ctx, order.ID, finalTotal, discountAmount); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return response.OrderDetail{}, err
	}

	order, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return response.OrderDetail{}, err
	}

	return response.NewOrderDetail(order), nil
}

func (s *orderService) List(ctx context.Context, payload requests.ListOrder) (dto.PaginationResponse[response.Order], error) {
	var paginateRes dto.PaginationResponse[response.Order]
	res, count, err := s.orderRepo.ListOrders(ctx, payload)
	if err != nil {
		return paginateRes, err
	}

	paginateRes = dto.NewPaginationResponse(payload.PaginationRequest, count, response.NewOrderList(res))
	return paginateRes, nil
}

func (s *orderService) Detail(ctx context.Context, id int64) (response.OrderDetail, error) {
	order, err := s.orderRepo.GetOrder(ctx, id)
	if err != nil {
		return response.OrderDetail{}, err
	}

	return response.NewOrderDetail(order), nil
}

func (s *orderService) UpdateStatus(ctx context.Context, id int64, payload requests.UpdateOrderStatus) error {
	order, err := s.orderRepo.GetOrder(ctx, id)
	if err != nil {
		return err
	}

	if (order.Status == constants.OrderStatusPaid && payload.Status != constants.OrderStatusReady) || order.Status == constants.OrderStatusCancelled || order.Status == constants.OrderStatusReady {
		return internal_err.NewDefaultError(http.StatusBadRequest, "Order cannot be modified")
	}

	return s.orderRepo.UpdateOrderStatus(ctx, id, payload.Status)
}

func (s *orderService) Cancel(ctx context.Context, id int64) error {
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		order, err := s.orderRepo.GetOrder(ctx, id)
		if err != nil {
			return err
		}

		if order.Status == constants.OrderStatusPaid || order.Status == constants.OrderStatusCancelled || order.Status == constants.OrderStatusReady {
			return internal_err.NewDefaultError(http.StatusBadRequest, "Order cannot be modified")
		}

		if err := s.orderRepo.UpdateOrderStatus(ctx, id, constants.OrderStatusCancelled); err != nil {
			return err
		}

		for _, item := range order.OrderItems {
			if err := s.productRepo.UpdateStock(ctx, tx, item.ProductID, item.Quantity); err != nil {
				return err
			}
			history := domain.StockHistory{
				ProductID: item.ProductID,
				Change:    item.Quantity,
				Reason:    fmt.Sprintf("Order #%d Cancelled", order.ID),
			}
			if err := s.stockHistoryRepo.CreateStockHistory(ctx, tx, &history); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (s *orderService) AddItem(ctx context.Context, orderID int64, payload requests.AddOrderItem) (response.OrderDetail, error) {
	var detail response.OrderDetail
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		order, err := s.orderRepo.GetOrder(ctx, orderID)
		if err != nil {
			return err
		}

		if order.Status == constants.OrderStatusPaid || order.Status == constants.OrderStatusCancelled || order.Status == constants.OrderStatusReady {
			return internal_err.NewDefaultError(http.StatusBadRequest, "Order cannot be modified")
		}

		product, err := s.productRepo.GetProduct(ctx, payload.ProductID)
		if err != nil {
			return err
		}

		if product.Stock < payload.Quantity {
			return internal_err.NewDefaultError(http.StatusBadRequest, fmt.Sprintf("Not enough stock for product: %s", product.Name))
		}

		item := domainOrderItemFromCreate(orderID, product.Price, payload)
		if err := s.orderItemRepo.CreateItem(ctx, &item); err != nil {
			return err
		}

		if err := s.productRepo.UpdateStock(ctx, tx, payload.ProductID, -payload.Quantity); err != nil {
			return err
		}
		history := domain.StockHistory{
			ProductID: payload.ProductID,
			Change:    -payload.Quantity,
			Reason:    fmt.Sprintf("Item added to Order #%d", order.ID),
		}
		if err := s.stockHistoryRepo.CreateStockHistory(ctx, tx, &history); err != nil {
			return err
		}

		total, err := s.orderItemRepo.SumSubtotalByOrder(ctx, orderID)
		if err != nil {
			return err
		}

		finalTotal, discountAmount := computeOrderAmounts(total, &order)
		if err := s.orderRepo.UpdateOrderAmounts(ctx, orderID, finalTotal, discountAmount); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return detail, err
	}

	order, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return detail, err
	}

	return response.NewOrderDetail(order), nil
}

func (s *orderService) RemoveItem(ctx context.Context, orderID int64, itemID int64) (response.OrderDetail, error) {
	var detail response.OrderDetail
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		order, err := s.orderRepo.GetOrder(ctx, orderID)
		if err != nil {
			return err
		}

		if order.Status == constants.OrderStatusPaid || order.Status == constants.OrderStatusCancelled || order.Status == constants.OrderStatusReady {
			return internal_err.NewDefaultError(http.StatusBadRequest, "Order cannot be modified")
		}

		item, err := s.orderItemRepo.GetItem(ctx, itemID)
		if err != nil {
			return err
		}
		if item.OrderID != orderID {
			return internal_err.NewDefaultError(http.StatusBadRequest, "Order item does not belong to order")
		}

		if err := s.orderItemRepo.DeleteItem(ctx, itemID); err != nil {
			return err
		}

		if err := s.productRepo.UpdateStock(ctx, tx, item.ProductID, item.Quantity); err != nil {
			return err
		}
		history := domain.StockHistory{
			ProductID: item.ProductID,
			Change:    item.Quantity,
			Reason:    fmt.Sprintf("Item removed from Order #%d", order.ID),
		}
		if err := s.stockHistoryRepo.CreateStockHistory(ctx, tx, &history); err != nil {
			return err
		}

		total, err := s.orderItemRepo.SumSubtotalByOrder(ctx, orderID)
		if err != nil {
			return err
		}

		finalTotal, discountAmount := computeOrderAmounts(total, &order)
		if err := s.orderRepo.UpdateOrderAmounts(ctx, orderID, finalTotal, discountAmount); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return detail, err
	}

	order, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return detail, err
	}

	return response.NewOrderDetail(order), nil
}

func domainOrderFromCreate(payload requests.CreateOrder, storeID int64, staffID int64) domain.Order {
	var discountType *string
	if payload.DiscountType != nil {
		dt := string(*payload.DiscountType)
		discountType = &dt
	}

	return domain.Order{
		StoreID:       storeID,
		TableID:       payload.TableID,
		StaffID:       staffID,
		DiscountType:  discountType,
		DiscountValue: payload.DiscountValue,
		TotalAmount:   0,
		Status:        constants.OrderStatusOpen,
	}
}

func domainOrderItemFromCreate(orderID int64, unitPrice float64, payload requests.AddOrderItem) domain.OrderItem {
	baseSubtotal := unitPrice * float64(payload.Quantity)

	var discountAmount float64
	var dtStr *string
	if payload.DiscountType != nil {
		dt := string(*payload.DiscountType)
		dtStr = &dt
		if *payload.DiscountType == constants.DiscountTypePercentage {
			discountAmount = baseSubtotal * (payload.DiscountValue / 100)
		} else if *payload.DiscountType == constants.DiscountTypeFixed {
			discountAmount = payload.DiscountValue
		}
	}

	if discountAmount > baseSubtotal {
		discountAmount = baseSubtotal
	}

	subtotal := baseSubtotal - discountAmount

	return domain.OrderItem{
		OrderID:        orderID,
		ProductID:      payload.ProductID,
		Quantity:       payload.Quantity,
		UnitPrice:      unitPrice,
		DiscountType:   dtStr,
		DiscountValue:  payload.DiscountValue,
		DiscountAmount: discountAmount,
		Subtotal:       subtotal,
	}
}

func computeOrderAmounts(total float64, order *domain.Order) (finalTotal float64, discountAmount float64) {
	discountAmount = 0.0
	if order.DiscountType != nil {
		if *order.DiscountType == string(constants.DiscountTypePercentage) {
			discountAmount = total * (order.DiscountValue / 100)
		} else if *order.DiscountType == string(constants.DiscountTypeFixed) {
			discountAmount = order.DiscountValue
		}
	}
	if discountAmount > total {
		discountAmount = total
	}

	finalTotal = total - discountAmount
	return finalTotal, discountAmount
}
