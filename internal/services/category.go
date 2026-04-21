package services

import (
	"backend-ta/internal/dto"
	"backend-ta/internal/dto/requests"
	"backend-ta/internal/dto/response"
	"backend-ta/internal/repository"
	"context"
)

type CategoryService interface {
	Create(ctx context.Context, payload requests.CreateCategory) (response.Category, error)
	Update(ctx context.Context, id int64, payload requests.UpdateCategory) error
	Delete(ctx context.Context, id int64) error
	Detail(ctx context.Context, id int64) (response.Category, error)
	List(ctx context.Context, payload requests.ListCategory) (dto.PaginationResponse[response.Category], error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategorySrv(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) Create(ctx context.Context, payload requests.CreateCategory) (response.Category, error) {
	category := payload.ToDomain()
	if err := s.categoryRepo.CreateCategory(ctx, &category); err != nil {
		return response.Category{}, err
	}

	return response.NewCategory(category), nil
}

func (s *categoryService) Update(ctx context.Context, id int64, payload requests.UpdateCategory) error {
	category, err := s.categoryRepo.GetCategory(ctx, id)
	if err != nil {
		return err
	}

	if payload.Name != "" {
		category.Name = payload.Name
	}
	if payload.ImageURL != "" {
		category.ImageURL = payload.ImageURL
	}

	return s.categoryRepo.UpdateCategory(ctx, &category)
}

func (s *categoryService) Delete(ctx context.Context, id int64) error {
	return s.categoryRepo.DeleteCategory(ctx, id)
}

func (s *categoryService) Detail(ctx context.Context, id int64) (response.Category, error) {
	category, err := s.categoryRepo.GetCategory(ctx, id)
	if err != nil {
		return response.Category{}, err
	}

	return response.NewCategory(category), nil
}

func (s *categoryService) List(ctx context.Context, payload requests.ListCategory) (dto.PaginationResponse[response.Category], error) {
	var paginateRes dto.PaginationResponse[response.Category]
	res, count, err := s.categoryRepo.ListCategory(ctx, payload)
	if err != nil {
		return paginateRes, err
	}

	paginateRes = dto.NewPaginationResponse(payload.PaginationRequest, count, response.NewCategoryList(res))
	return paginateRes, nil
}
