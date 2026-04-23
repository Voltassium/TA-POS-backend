package services

import (
	"backend-ta/internal/repository"
	"sync"
)

var once = sync.Once{}
var ServicePool *PoolService

type PoolService struct {
	AuthService     AuthService
	UserService     UserService
	CategoryService CategoryService
	ProductService  ProductService
	OrderService    OrderService
	PaymentService  PaymentService
	StockHistoryService StockHistoryService
	StatisticsService   StatisticsService
}

func Init() {
	once.Do(func() {
		repo := repository.RepoPool
		ServicePool = &PoolService{
			AuthService:     NewAuthSrv(repo.UserRepository),
			UserService:     NewUserSrv(repo.UserRepository),
			CategoryService: NewCategorySrv(repo.CategoryRepository),
			ProductService:  NewProductSrv(repo.ProductRepository, repo.CategoryRepository),
			OrderService:    NewOrderSrv(repo.OrderRepository, repo.OrderItemRepository, repo.ProductRepository),
			PaymentService:  NewPaymentSrv(repo.OrderRepository, repo.PaymentRepository, repo.ProductRepository, repo.StockHistoryRepository),
			StockHistoryService: NewStockHistorySrv(repo.StockHistoryRepository),
			StatisticsService:   NewStatisticsSrv(repo.StatisticsRepository),
		}

	})

}
