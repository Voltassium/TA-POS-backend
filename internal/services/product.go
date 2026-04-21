package services

import (
	"backend-ta/internal/dto"
	"backend-ta/internal/dto/requests"
	"backend-ta/internal/dto/response"
	"backend-ta/internal/repository"
	"context"
)

type ProductService interface {
	Create(ctx context.Context, payload requests.CreateProduct) (response.Product, error)
	Update(ctx context.Context, id int64, payload requests.UpdateProduct) error
	Delete(ctx context.Context, id int64) error
	Detail(ctx context.Context, id int64) (response.Product, error)
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
	if _, err := s.categoryRepo.GetCategory(ctx, payload.CategoryID); err != nil {
		return response.Product{}, err
	}

	product := payload.ToDomain()
	if err := s.productRepo.CreateProduct(ctx, &product); err != nil {
		return response.Product{}, err
	}

	product, err := s.productRepo.GetProduct(ctx, product.ID)
	if err != nil {
		return response.Product{}, err
	}

	return response.NewProduct(product), nil
}

func (s *productService) Update(ctx context.Context, id int64, payload requests.UpdateProduct) error {
	product, err := s.productRepo.GetProduct(ctx, id)
	if err != nil {
		return err
	}

	if payload.CategoryID != 0 && payload.CategoryID != product.CategoryID {
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

	return s.productRepo.UpdateProduct(ctx, &product)
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	return s.productRepo.DeleteProduct(ctx, id)
}

func (s *productService) Detail(ctx context.Context, id int64) (response.Product, error) {
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
