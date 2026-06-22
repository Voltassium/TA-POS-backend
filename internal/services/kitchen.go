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
	UpdateItemServedQty(ctx context.Context, orderID string, itemID string, payload requests.UpdateServedQty) (response.OrderDetail, error)
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

func (s *kitchenService) UpdateItemServedQty(ctx context.Context, orderID string, itemID string, payload requests.UpdateServedQty) (response.OrderDetail, error) {
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

	if payload.ServedQty > item.Quantity {
		return response.OrderDetail{}, internal_err.NewDefaultError(
			http.StatusBadRequest,
			fmt.Sprintf("Jumlah disajikan (%d) tidak boleh melebihi jumlah pesanan (%d)", payload.ServedQty, item.Quantity),
		)
	}

	if err := s.orderItemRepo.UpdateServedQty(ctx, itemID, payload.ServedQty); err != nil {
		return response.OrderDetail{}, err
	}

	allItems, err := s.orderItemRepo.ListItemsByOrder(ctx, orderID)
	if err != nil {
		return response.OrderDetail{}, err
	}

	allServed := true
	for _, it := range allItems {
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
		if err := s.orderRepo.UpdateOrderStatus(ctx, orderID, constants.OrderStatusCompleted); err != nil {
			return response.OrderDetail{}, err
		}
	}

	refreshed, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return response.OrderDetail{}, err
	}

	return response.NewOrderDetail(refreshed), nil
}
