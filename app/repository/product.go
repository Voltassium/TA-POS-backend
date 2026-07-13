package repository

import (
	"backend-ta/app/domain"
	"backend-ta/app/dto/requests"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/database"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, data *domain.Product) error
	UpdateProduct(ctx context.Context, data *domain.Product) error
	DeleteProduct(ctx context.Context, id string) error
	GetProduct(ctx context.Context, id string) (domain.Product, error)
	GetProductForUpdate(ctx context.Context, tx bun.Tx, id string) (domain.Product, error)
	ListProduct(ctx context.Context, req requests.ListProduct) ([]domain.Product, int, error)
	UpdateStock(ctx context.Context, tx bun.Tx, productID string, change int) error
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

func (r *productRepository) DeleteProduct(ctx context.Context, id string) error {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	_, err := r.db.InitQuery(ctx).
		NewDelete().
		Model((*domain.Product)(nil)).
		Where("id = ?", id).
		Where("store_id = ?", storeID).
		Exec(ctx)
	return err
}

func (r *productRepository) GetProduct(ctx context.Context, id string) (domain.Product, error) {
	var res domain.Product
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	err := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("Category").
		Where("product.id = ?", id).
		Where("product.store_id = ?", storeID).
		Scan(ctx)
	return res, err
}

func (r *productRepository) GetProductForUpdate(ctx context.Context, tx bun.Tx, id string) (domain.Product, error) {
	var res domain.Product
	err := tx.NewRaw(
		`SELECT p.*, c.id AS "category__id", c.name AS "category__name", c.store_id AS "category__store_id"
		 FROM products AS p
		 LEFT JOIN categories AS c ON c.id = p.category_id
		 WHERE p.id = ?
		 FOR UPDATE`, id,
	).Scan(ctx, &res)
	return res, err
}

func (r *productRepository) ListProduct(ctx context.Context, req requests.ListProduct) ([]domain.Product, int, error) {
	var res []domain.Product
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Relation("Category").
		Where("product.store_id = ?", storeID)

	if req.CategoryID != "" {
		q.Where("category_id = ?", req.CategoryID)
	}
	if req.ProductType != "" {
		q.Where("product_type = ?", req.ProductType)
	}
	if req.Search != "" {
		q.Where("product.name ILIKE ? OR product.sku ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	q.Limit(req.PageSize).
		Offset(req.CalculateOffset()).
		Order(fmt.Sprintf("%s %s", req.OrderBy, req.OrderDir))

	total, err := q.ScanAndCount(ctx)
	return res, total, err
}

func (r *productRepository) UpdateStock(ctx context.Context, tx bun.Tx, productID string, change int) error {
	_, err := tx.NewUpdate().
		Model((*domain.Product)(nil)).
		Set("stock = stock + ?", change).
		Set("is_available = CASE WHEN stock + ? <= 0 THEN false WHEN stock = 0 AND ? > 0 THEN true ELSE is_available END", change, change).
		Where("id = ?", productID).
		Exec(ctx)
	return err
}
