package services

import (
	"backend-ta/app/dto"
	"backend-ta/app/dto/requests"
	"backend-ta/app/dto/response"
	"backend-ta/app/repository"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/errors"
	"context"
	"net/http"
)

type ProductService interface {
	Create(ctx context.Context, payload requests.CreateProduct) (response.Product, error)
	Update(ctx context.Context, id string, payload requests.UpdateProduct) error
	Delete(ctx context.Context, id string) error
	Detail(ctx context.Context, id string) (response.Product, error)
	List(ctx context.Context, payload requests.ListProduct) (dto.PaginationResponse[response.Product], error)
}

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductSrv(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productService{productRepo: productRepo, categoryRepo: categoryRepo}
}

func (s *productService) Create(ctx context.Context, payload requests.CreateProduct) (response.Product, error) {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	if storeID == 0 {
		return response.Product{}, errors.NewDefaultError(http.StatusUnauthorized, "Invalid store context")
	}

	if _, err := s.categoryRepo.GetCategory(ctx, payload.CategoryID); err != nil {
		return response.Product{}, err
	}

	product := payload.ToDomain()
	product.StoreID = storeID
	if err := s.productRepo.CreateProduct(ctx, &product); err != nil {
		return response.Product{}, err
	}

	product, err := s.productRepo.GetProduct(ctx, product.ID)
	if err != nil {
		return response.Product{}, err
	}

	return response.NewProduct(product), nil
}

func (s *productService) Update(ctx context.Context, id string, payload requests.UpdateProduct) error {
	product, err := s.productRepo.GetProduct(ctx, id)
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
	if payload.Description != "" {
		product.Description = payload.Description
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

	return s.productRepo.UpdateProduct(ctx, &product)
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
