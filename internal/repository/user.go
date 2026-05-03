package repository

import (
	"backend-ta/internal/domain"
	"backend-ta/internal/dto/requests"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/database"
	"backend-ta/pkg/errors"
	"context"
	"fmt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, data *domain.User) error
	ListUser(ctx context.Context, req requests.ListUser) ([]domain.User, int, error)
	UpdateUser(ctx context.Context, data *domain.User) error
	DeleteUser(ctx context.Context, id int64) error
	GetUser(ctx context.Context, id int64) (res domain.User, err error)
	GetUserByEmail(ctx context.Context, email string) (res domain.User, err error)
}

type userRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, data *domain.User) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Returning("id").Exec(ctx)
	if err != nil {
		return errors.CheckUniqueViolation(err)
	}
	return err
}

func (r *userRepository) ListUser(ctx context.Context, req requests.ListUser) ([]domain.User, int, error) {
	var res []domain.User
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	q := r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Where("store_id = ?", storeID).
		Limit(req.PageSize).
		Offset(req.CalculateOffset()).
		Order(fmt.Sprintf("%s %s", req.OrderBy, req.OrderDir))

	total, err := q.ScanAndCount(ctx)
	return res, total, err
}

func (r *userRepository) UpdateUser(ctx context.Context, data *domain.User) error {
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

func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	_, err := r.db.InitQuery(ctx).
		NewDelete().
		Model((*domain.User)(nil)).
		Where("id = ?", id).
		Where("store_id = ?", storeID).
		Exec(ctx)
	return err
}

func (r *userRepository) GetUser(ctx context.Context, id int64) (res domain.User, err error) {
	storeID := authentication.GetUserDataFromToken(ctx).StoreID
	err = r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Where("id = ?", id).
		Where("store_id = ?", storeID).
		Scan(ctx)
	return res, err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (res domain.User, err error) {
	err = r.db.InitQuery(ctx).
		NewSelect().
		Model(&res).
		Where("email = ?", email).
		Scan(ctx)
	return res, err
}
