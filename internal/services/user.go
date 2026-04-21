package services

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/dto"
	"backend-ta/internal/dto/requests"
	"backend-ta/internal/dto/response"
	"backend-ta/internal/repository"
	"backend-ta/pkg/authentication"
	"context"
)

type UserService interface {
	Register(ctx context.Context, payload requests.CreateUser) error
	GetList(ctx context.Context, payload requests.ListUser) (dto.PaginationResponse[response.User], error)
	Update(ctx context.Context, id int64, payload requests.UpdateUser) error
	DeleteSrv(ctx context.Context, id int64) error
	Detail(ctx context.Context, id int64) (response.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserSrv(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (a *userService) Register(ctx context.Context, payload requests.CreateUser) error {
	user := payload.ToDomain()
	if user.Role == "" {
		user.Role = constants.UserRoleStaff
	}

	hashedPassword, err := authentication.HashPassword(payload.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return a.userRepo.CreateUser(ctx, &user)
}

func (a *userService) GetList(ctx context.Context, payload requests.ListUser) (dto.PaginationResponse[response.User], error) {
	var paginateRes dto.PaginationResponse[response.User]
	res, count, err := a.userRepo.ListUser(ctx, payload)
	if err != nil {
		return paginateRes, err
	}

	paginateRes = dto.NewPaginationResponse(payload.PaginationRequest, count, response.NewListUser(res))
	return paginateRes, nil
}

func (a *userService) Update(ctx context.Context, id int64, payload requests.UpdateUser) error {
	userData, err := a.userRepo.GetUser(ctx, id)
	if err != nil {
		return err
	}

	if payload.Email != "" {
		userData.Email = payload.Email
	}
	if payload.Password != "" {
		hashedPassword, err := authentication.HashPassword(payload.Password)
		if err != nil {
			return err
		}
		userData.Password = hashedPassword
	}
	if payload.Role != "" {
		userData.Role = payload.Role
	}

	return a.userRepo.UpdateUser(ctx, &userData)
}

func (a *userService) DeleteSrv(ctx context.Context, id int64) error {
	return a.userRepo.DeleteUser(ctx, id)
}

func (a *userService) Detail(ctx context.Context, id int64) (response.User, error) {
	var res response.User
	data, err := a.userRepo.GetUser(ctx, id)
	if err != nil {
		return res, err
	}

	res = response.NewUser(data)
	return res, nil
}
