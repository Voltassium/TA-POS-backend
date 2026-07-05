package services

import (
	"backend-ta/app/constants"
	"backend-ta/app/domain"
	"backend-ta/app/dto"
	"backend-ta/app/dto/requests"
	"backend-ta/app/dto/response"
	"backend-ta/app/repository"
	"backend-ta/pkg/authentication"
	"context"
	"fmt"
)

type UserService interface {
	Register(ctx context.Context, payload requests.CreateUser) error
	RegisterByAdmin(ctx context.Context, payload requests.CreateUserByAdmin) error
	GetList(ctx context.Context, payload requests.ListUser) (dto.PaginationResponse[response.User], error)
	Update(ctx context.Context, id string, payload requests.UpdateUser) error
	DeleteSrv(ctx context.Context, id string) error
	Detail(ctx context.Context, id string) (response.User, error)
}

type userService struct {
	userRepo  repository.UserRepository
	storeRepo repository.StoreRepository
}

func NewUserSrv(userRepo repository.UserRepository, storeRepo repository.StoreRepository) UserService {
	return &userService{
		userRepo:  userRepo,
		storeRepo: storeRepo,
	}
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

	if user.Role == constants.UserRoleOwner {
		store := domain.Store{
			Name:    payload.StoreName,
			Address: payload.StoreAddress,
		}
		if store.Name == "" {
			store.Name = "My Store"
		}
		if err := a.storeRepo.CreateStore(ctx, &store); err != nil {
			return err
		}
		user.StoreID = &store.ID
	}

	return a.userRepo.CreateUser(ctx, &user)
}

func (a *userService) RegisterByAdmin(ctx context.Context, payload requests.CreateUserByAdmin) error {
	tokenData := authentication.GetUserDataFromToken(ctx)

	switch tokenData.Role {
	case constants.UserRoleOwner:
		if payload.Role != constants.UserRoleChef && payload.Role != constants.UserRoleStaff {
			return fmt.Errorf("role tidak valid: owner hanya dapat mendaftarkan chef atau staff")
		}
	case constants.UserRoleSuperadmin:
		if !payload.Role.IsValidEnum() {
			return fmt.Errorf("role tidak valid")
		}
	default:
		return fmt.Errorf("anda tidak memiliki izin untuk mendaftarkan akun")
	}

	user := payload.ToDomain()
	hashedPassword, err := authentication.HashPassword(payload.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	if payload.Role != constants.UserRoleSuperadmin {
		storeID := tokenData.StoreID
		if storeID == 0 {
			storeID = 1 // Default fallback if superadmin creates owner/staff without store context
		}
		user.StoreID = &storeID
	} else {
		user.StoreID = nil
	}

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

func (a *userService) Update(ctx context.Context, id string, payload requests.UpdateUser) error {
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

func (a *userService) DeleteSrv(ctx context.Context, id string) error {
	return a.userRepo.DeleteUser(ctx, id)
}

func (a *userService) Detail(ctx context.Context, id string) (response.User, error) {
	var res response.User
	data, err := a.userRepo.GetUser(ctx, id)
	if err != nil {
		return res, err
	}

	res = response.NewUser(data)
	return res, nil
}
