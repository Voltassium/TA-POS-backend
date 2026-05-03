package repository

import (
	"backend-ta/pkg/database"
	"sync"
)

var once = &sync.Once{}
var RepoPool *PoolRepository

type PoolRepository struct {
	UserRepository         UserRepository
	StoreRepository        StoreRepository
	CategoryRepository     CategoryRepository
	ProductRepository      ProductRepository
	OrderRepository        OrderRepository
	OrderItemRepository    OrderItemRepository
	PaymentRepository      PaymentRepository
	StockHistoryRepository StockHistoryRepository
	StatisticsRepository   StatisticsRepository
}

func Init(db *database.Database) {
	once.Do(func() {
		RepoPool = &PoolRepository{
			UserRepository:         NewUserRepository(db),
			StoreRepository:        NewStoreRepository(db),
			CategoryRepository:     NewCategoryRepository(db),
			ProductRepository:      NewProductRepository(db),
			OrderRepository:        NewOrderRepository(db),
			OrderItemRepository:    NewOrderItemRepository(db),
			PaymentRepository:      NewPaymentRepository(db),
			StockHistoryRepository: NewStockHistoryRepository(db),
			StatisticsRepository:   NewStatisticsRepository(db),
		}
	})

}
