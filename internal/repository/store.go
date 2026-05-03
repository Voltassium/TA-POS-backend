package repository

import (
	"backend-ta/internal/domain"
	"backend-ta/pkg/database"
	"context"
	"fmt"
)

type StoreRepository interface {
	CreateStore(ctx context.Context, data *domain.Store) error
	UpdateStore(ctx context.Context, data *domain.Store) error
	DeleteStore(ctx context.Context, id int64) error
	GetStore(ctx context.Context, id int64) (domain.Store, error)
	ListStores(ctx context.Context, page, pageSize int, orderBy, orderDir string) ([]domain.Store, int, error)
}

type storeRepository struct {
	db *database.Database
}

func NewStoreRepository(db *database.Database) StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) CreateStore(ctx context.Context, data *domain.Store) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	return err
}

func (r *storeRepository) UpdateStore(ctx context.Context, data *domain.Store) error {
	_, err := r.db.InitQuery(ctx).
		NewUpdate().
		Model(data).
		Where("id = ?", data.ID).
		ExcludeColumn("created_at").
		Returning("id").
		Exec(ctx)
	return err
}

func (r *storeRepository) DeleteStore(ctx context.Context, id int64) error {
	_, err := r.db.InitQuery(ctx).
		NewDelete().
		Model((*domain.Store)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *storeRepository) GetStore(ctx context.Context, id int64) (domain.Store, error) {
	var res domain.Store
	err := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Where("id = ?", id).
		Scan(ctx)
	return res, err
}

func (r *storeRepository) ListStores(ctx context.Context, page, pageSize int, orderBy, orderDir string) ([]domain.Store, int, error) {
	var res []domain.Store
	offset := (page - 1) * pageSize
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Limit(pageSize).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", orderBy, orderDir))

	total, err := q.ScanAndCount(ctx)
	return res, total, err
}
