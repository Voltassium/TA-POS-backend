package repository

import (
	"backend-ta/app/domain"
	"backend-ta/app/dto/requests"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/database"
	"context"
	"fmt"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, data *domain.Category) error
	UpdateCategory(ctx context.Context, data *domain.Category) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategory(ctx context.Context, id string) (domain.Category, error)
	ListCategory(ctx context.Context, req requests.ListCategory) ([]domain.Category, int, error)
}

type categoryRepository struct {
	db *database.Database
}

func NewCategoryRepository(db *database.Database) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(ctx context.Context, data *domain.Category) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	return err
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, data *domain.Category) error {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model(data).
		Where("id = ?", data.ID).
		Where("store_id = ?", storeID).
		ExcludeColumn("created_at").
		Returning("id").
		Exec(ctx)
	return err
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id string) error {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	_, err := r.db.InitQuery(ctx).
		NewDelete().
		Model((*domain.Category)(nil)).
		Where("id = ?", id).
		Where("store_id = ?", storeID).
		Exec(ctx)
	return err
}

func (r *categoryRepository) GetCategory(ctx context.Context, id string) (domain.Category, error) {
	var res domain.Category
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	err := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Where("id = ?", id).
		Where("store_id = ?", storeID).
		Scan(ctx)
	return res, err
}

func (r *categoryRepository) ListCategory(ctx context.Context, req requests.ListCategory) ([]domain.Category, int, error) {
	var res []domain.Category
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Where("store_id = ?", storeID)

	if req.Search != "" {
		q.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	q.Limit(req.PageSize).
		Offset(req.CalculateOffset()).
		Order(fmt.Sprintf("%s %s", req.OrderBy, req.OrderDir))

	total, err := q.ScanAndCount(ctx)
	return res, total, err
}
