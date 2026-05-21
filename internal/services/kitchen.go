package services

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/dto/requests"
	"backend-ta/internal/dto/response"
	"backend-ta/internal/repository"
	internal_err "backend-ta/pkg/errors"
	"context"
	"fmt"
	"net/http"
)

type KitchenService interface {
	UpdateItemServedQty(ctx context.Context, orderID int64, itemID int64, payload requests.UpdateServedQty) (response.OrderDetail, error)
}

type kitchenService struct {
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
}

func NewKitchenSrv(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
) KitchenService {
	return &kitchenService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

func (s *kitchenService) UpdateItemServedQty(ctx context.Context, orderID int64, itemID int64, payload requests.UpdateServedQty) (response.OrderDetail, error) {
	// 1. Fetch and validate the order
	order, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return response.OrderDetail{}, err
	}

	if order.Status != constants.OrderStatusPaid {
		return response.OrderDetail{}, internal_err.NewDefaultError(
			http.StatusBadRequest,
			fmt.Sprintf("Hanya pesanan dengan status 'Paid' yang bisa diperbarui (status saat ini: %s)", order.Status),
		)
	}

	// 2. Fetch and validate the item belongs to this order
	item, err := s.orderItemRepo.GetItem(ctx, itemID)
	if err != nil {
		return response.OrderDetail{}, err
	}

	if item.OrderID != orderID {
		return response.OrderDetail{}, internal_err.NewDefaultError(
			http.StatusBadRequest,
			"Item tidak ditemukan dalam pesanan ini",
		)
	}

	// 3. Validate served_qty bounds
	if payload.ServedQty > item.Quantity {
		return response.OrderDetail{}, internal_err.NewDefaultError(
			http.StatusBadRequest,
			fmt.Sprintf("Jumlah disajikan (%d) tidak boleh melebihi jumlah pesanan (%d)", payload.ServedQty, item.Quantity),
		)
	}

	// 4. Update served_qty
	if err := s.orderItemRepo.UpdateServedQty(ctx, itemID, payload.ServedQty); err != nil {
		return response.OrderDetail{}, err
	}

	// 5. Check if all items are now fully served → auto-complete order
	allItems, err := s.orderItemRepo.ListItemsByOrder(ctx, orderID)
	if err != nil {
		return response.OrderDetail{}, err
	}

	allServed := true
	for _, it := range allItems {
		// Use the updated value for the item we just changed
		servedQty := it.ServedQty
		if it.ID == itemID {
			servedQty = payload.ServedQty
		}
		if servedQty < it.Quantity {
			allServed = false
			break
		}
	}

	if allServed {
		if err := s.orderRepo.UpdateOrderStatus(ctx, orderID, constants.OrderStatusReady); err != nil {
			return response.OrderDetail{}, err
		}
	}

	// 6. Return refreshed order detail
	refreshed, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return response.OrderDetail{}, err
	}

	return response.NewOrderDetail(refreshed), nil
}
