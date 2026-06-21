package repository

import (
	"backend-ta/internal/domain"
	"backend-ta/internal/dto/requests"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/database"
	"context"
)

type PengeluaranRepository interface {
	Create(ctx context.Context, data *domain.Pengeluaran) error
	Update(ctx context.Context, data *domain.Pengeluaran) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.Pengeluaran, error)
	List(ctx context.Context, req requests.ListPengeluaran) ([]domain.Pengeluaran, int, error)
}

type pengeluaranRepository struct {
	db *database.Database
}

func NewPengeluaranRepository(db *database.Database) PengeluaranRepository {
	return &pengeluaranRepository{db: db}
}

func (r *pengeluaranRepository) Create(ctx context.Context, data *domain.Pengeluaran) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	return err
}

func (r *pengeluaranRepository) Update(ctx context.Context, data *domain.Pengeluaran) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().Model(data).Where("id = ? AND store_id = ?", data.ID, data.StoreID).Exec(ctx)
	return err
}

func (r *pengeluaranRepository) Delete(ctx context.Context, id string) error {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	_, err := r.db.InitQuery(ctx).NewDelete().Model(&domain.Pengeluaran{}).Where("id = ? AND store_id = ?", id, storeID).Exec(ctx)
	return err
}

func (r *pengeluaranRepository) Get(ctx context.Context, id string) (domain.Pengeluaran, error) {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	var res domain.Pengeluaran
	err := r.db.InitQuery(ctx).NewSelect().Model(&res).Where("id = ? AND store_id = ?", id, storeID).Scan(ctx)
	return res, err
}

func (r *pengeluaranRepository) List(ctx context.Context, req requests.ListPengeluaran) ([]domain.Pengeluaran, int, error) {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	var res []domain.Pengeluaran

	query := r.db.InitQuery(ctx).NewSelect().Model(&res).Where("store_id = ?", storeID)

	if req.StartDate != nil && *req.StartDate != "" {
		query.Where("tanggal >= ?", *req.StartDate)
	}
	if req.EndDate != nil && *req.EndDate != "" {
		query.Where("tanggal <= ?", *req.EndDate)
	}
	if req.Search != "" {
		query.Where("(category ILIKE ? OR description ILIKE ?)", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	count, err := query.Limit(req.PageSize).Offset(req.CalculateOffset()).OrderExpr("tanggal DESC, created_at DESC").ScanAndCount(ctx)
	return res, count, err
}
