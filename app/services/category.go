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

type CategoryService interface {
	Create(ctx context.Context, payload requests.CreateCategory) (response.Category, error)
	Update(ctx context.Context, id string, payload requests.UpdateCategory) error
	Delete(ctx context.Context, id string) error
	Detail(ctx context.Context, id string) (response.Category, error)
	List(ctx context.Context, payload requests.ListCategory) (dto.PaginationResponse[response.Category], error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategorySrv(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) Create(ctx context.Context, payload requests.CreateCategory) (response.Category, error) {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	if storeID == 0 {
		return response.Category{}, errors.NewDefaultError(http.StatusUnauthorized, "Invalid store context")
	}

	category := payload.ToDomain()
	category.StoreID = storeID
	if err := s.categoryRepo.CreateCategory(ctx, &category); err != nil {
		return response.Category{}, err
	}

	return response.NewCategory(category), nil
}

func (s *categoryService) Update(ctx context.Context, id string, payload requests.UpdateCategory) error {
	category, err := s.categoryRepo.GetCategory(ctx, id)
	if err != nil {
		return err
	}

	if payload.Name != "" {
		category.Name = payload.Name
	}
	return s.categoryRepo.UpdateCategory(ctx, &category)
}

func (s *categoryService) Delete(ctx context.Context, id string) error {
	return s.categoryRepo.DeleteCategory(ctx, id)
}

func (s *categoryService) Detail(ctx context.Context, id string) (response.Category, error) {
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
