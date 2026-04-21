package repository

import (
	"backend-ta/pkg/database"
	"sync"
)

var once = &sync.Once{}
var RepoPool *PoolRepository

type PoolRepository struct {
	UserRepository      UserRepository
	CategoryRepository  CategoryRepository
	ProductRepository   ProductRepository
	OrderRepository     OrderRepository
	OrderItemRepository OrderItemRepository
	PaymentRepository   PaymentRepository
}

func Init(db *database.Database) {
	once.Do(func() {
		RepoPool = &PoolRepository{
			UserRepository:      NewUserRepository(db),
			CategoryRepository:  NewCategoryRepository(db),
			ProductRepository:   NewProductRepository(db),
			OrderRepository:     NewOrderRepository(db),
			OrderItemRepository: NewOrderItemRepository(db),
			PaymentRepository:   NewPaymentRepository(db),
		}
	})

}
