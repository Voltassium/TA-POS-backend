package services

import (
	"backend-ta/app/domain"
	"backend-ta/app/dto"
	"backend-ta/app/dto/requests"
	"backend-ta/app/dto/response"
	"backend-ta/app/repository"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/database"
	"backend-ta/pkg/errors"
	"context"
	"database/sql"
	"net/http"

	"github.com/uptrace/bun"
)

type ProductService interface {
	Create(ctx context.Context, payload requests.CreateProduct) (response.Product, error)
	Update(ctx context.Context, id string, payload requests.UpdateProduct) error
	Delete(ctx context.Context, id string) error
	Detail(ctx context.Context, id string) (response.Product, error)
	List(ctx context.Context, payload requests.ListProduct) (dto.PaginationResponse[response.Product], error)
	Restock(ctx context.Context, id string, payload requests.RestockProduct) error
}

type productService struct {
	productRepo      repository.ProductRepository
	categoryRepo     repository.CategoryRepository
	stockHistoryRepo repository.StockHistoryRepository
}

func NewProductSrv(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	stockHistoryRepo repository.StockHistoryRepository,
) ProductService {
	return &productService{
		productRepo:      productRepo,
		categoryRepo:     categoryRepo,
		stockHistoryRepo: stockHistoryRepo,
	}
}

func (s *productService) Create(ctx context.Context, payload requests.CreateProduct) (response.Product, error) {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	if storeID == 0 {
		return response.Product{}, errors.NewDefaultError(http.StatusUnauthorized, "Invalid store context")
	}

	if _, err := s.categoryRepo.GetCategory(ctx, payload.CategoryID); err != nil {
		return response.Product{}, err
	}

	var product domain.Product
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		product = payload.ToDomain()
		product.StoreID = storeID
		if product.Stock <= 0 {
			product.IsAvailable = false
		}
		if err := s.productRepo.CreateProduct(ctx, &product); err != nil {
			return err
		}

		if product.Stock > 0 {
			history := domain.StockHistory{
				ProductID: product.ID,
				Change:    product.Stock,
				Reason:    "Stok Awal Produk Baru",
			}
			if err := s.stockHistoryRepo.CreateStockHistory(ctx, tx, &history); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return response.Product{}, err
	}

	product, err = s.productRepo.GetProduct(ctx, product.ID)
	if err != nil {
		return response.Product{}, err
	}

	return response.NewProduct(product), nil
}

func (s *productService) Update(ctx context.Context, id string, payload requests.UpdateProduct) error {
	var product domain.Product
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var err error
		product, err = s.productRepo.GetProduct(ctx, id)
		if err != nil {
			return err
		}

		if payload.CategoryID != "" && payload.CategoryID != product.CategoryID {
			if _, err := s.categoryRepo.GetCategory(ctx, payload.CategoryID); err != nil {
				return err
			}
			product.CategoryID = payload.CategoryID
		}
		if payload.Name != "" {
			product.Name = payload.Name
		}
		if payload.Price != 0 {
			product.Price = payload.Price
		}
		if payload.IsAvailable != nil {
			product.IsAvailable = *payload.IsAvailable
		}
		if payload.SKU != nil {
			product.SKU = payload.SKU
		}
		if payload.ProductType != "" {
			product.ProductType = payload.ProductType
			if payload.ProductType == "Olahan" {
				product.HargaBeli = nil
			}
		}
		if payload.HargaBeli != nil && product.ProductType == "Kulakan" {
			product.HargaBeli = payload.HargaBeli
		}

		var stockChange int
		if payload.Stock != nil && *payload.Stock != product.Stock {
			stockChange = *payload.Stock - product.Stock
			product.Stock = *payload.Stock
		}

		if product.Stock <= 0 {
			product.IsAvailable = false
		} else if payload.IsAvailable == nil && product.Stock > 0 && stockChange != 0 && product.Stock - stockChange <= 0 {
			product.IsAvailable = true
		}

		if err := s.productRepo.UpdateProduct(ctx, &product); err != nil {
			return err
		}

		if stockChange != 0 {
			history := domain.StockHistory{
				ProductID: product.ID,
				Change:    stockChange,
				Reason:    "Penyesuaian Stok (Manual)",
			}
			if err := s.stockHistoryRepo.CreateStockHistory(ctx, tx, &history); err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

func (s *productService) Delete(ctx context.Context, id string) error {
	return s.productRepo.DeleteProduct(ctx, id)
}

func (s *productService) Detail(ctx context.Context, id string) (response.Product, error) {
	product, err := s.productRepo.GetProduct(ctx, id)
	if err != nil {
		return response.Product{}, err
	}
	return response.NewProduct(product), nil
}

func (s *productService) List(ctx context.Context, payload requests.ListProduct) (dto.PaginationResponse[response.Product], error) {
	var paginateRes dto.PaginationResponse[response.Product]
	res, count, err := s.productRepo.ListProduct(ctx, payload)
	if err != nil {
		return paginateRes, err
	}

	paginateRes = dto.NewPaginationResponse(payload.PaginationRequest, count, response.NewProductList(res))
	return paginateRes, nil
}

func (s *productService) Restock(ctx context.Context, id string, payload requests.RestockProduct) error {
	err := database.RunInTx(ctx, database.GetDB(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		product, err := s.productRepo.GetProduct(ctx, id)
		if err != nil {
			return err
		}

		if product.ProductType != "Kulakan" {
			return errors.NewDefaultError(http.StatusBadRequest, "Hanya produk tipe Kulakan yang dapat direstock")
		}

		product.HargaBeli = &payload.HargaBeli
		product.Stock += payload.JumlahStok
		product.IsAvailable = true

		if err := s.productRepo.UpdateProduct(ctx, &product); err != nil {
			return err
		}

		history := domain.StockHistory{
			ProductID: product.ID,
			Change:    payload.JumlahStok,
			Reason:    "Pembelian Stok (Kulakan)",
		}
		if err := s.stockHistoryRepo.CreateStockHistory(ctx, tx, &history); err != nil {
			return err
		}

		return nil
	})
	return err
}

