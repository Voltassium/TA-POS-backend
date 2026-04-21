package repository

import (
	"backend-ta/internal/domain"
	"backend-ta/internal/dto/requests"
	"backend-ta/pkg/database"
	"context"
	"fmt"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, data *domain.Product) error
	UpdateProduct(ctx context.Context, data *domain.Product) error
	DeleteProduct(ctx context.Context, id int64) error
	GetProduct(ctx context.Context, id int64) (domain.Product, error)
	ListProduct(ctx context.Context, req requests.ListProduct) ([]domain.Product, int, error)
}

type productRepository struct {
	db *database.Database
}

func NewProductRepository(db *database.Database) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(ctx context.Context, data *domain.Product) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	return err
}

func (r *productRepository) UpdateProduct(ctx context.Context, data *domain.Product) error {
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model(data).
		Where("id = ?", data.ID).
		ExcludeColumn("created_at").
		Returning("id").
		Exec(ctx)
	return err
}

func (r *productRepository) DeleteProduct(ctx context.Context, id int64) error {
	_, err := r.db.InitQuery(ctx).
		NewDelete().
		Model((*domain.Product)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *productRepository) GetProduct(ctx context.Context, id int64) (domain.Product, error) {
	var res domain.Product
	err := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("Category").
		Where("product.id = ?", id).
		Scan(ctx)
	return res, err
}

func (r *productRepository) ListProduct(ctx context.Context, req requests.ListProduct) ([]domain.Product, int, error) {
	var res []domain.Product
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("Category").
		Limit(req.PageSize).
		Offset(req.CalculateOffset()).
		Order(fmt.Sprintf("%s %s", req.OrderBy, req.OrderDir))

	if req.CategoryID != 0 {
		q.Where("category_id = ?", req.CategoryID)
	}

	total, err := q.ScanAndCount(ctx)
	return res, total, err
}
