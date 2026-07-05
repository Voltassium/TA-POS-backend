package services

import (
	"backend-ta/app/repository"
	"sync"
)

var once = sync.Once{}
var ServicePool *PoolService

type PoolService struct {
	AuthService         AuthService
	UserService         UserService
	StoreService        StoreService
	CategoryService     CategoryService
	ProductService      ProductService
	OrderService        OrderService
	PaymentService      PaymentService
	StockHistoryService StockHistoryService
	StatisticsService   StatisticsService
	KitchenService      KitchenService
	PengeluaranService  PengeluaranService
}

func Init() {
	once.Do(func() {
		repo := repository.RepoPool
		ServicePool = &PoolService{
			AuthService:         NewAuthSrv(repo.UserRepository),
			UserService:         NewUserSrv(repo.UserRepository, repo.StoreRepository),
			StoreService:        NewStoreService(repo),
			CategoryService:     NewCategorySrv(repo.CategoryRepository),
			ProductService:      NewProductSrv(repo.ProductRepository, repo.CategoryRepository, repo.StockHistoryRepository),
			OrderService:        NewOrderSrv(repo.OrderRepository, repo.OrderItemRepository, repo.ProductRepository, repo.StockHistoryRepository),
			PaymentService:      NewPaymentSrv(repo.OrderRepository, repo.PaymentRepository, repo.ProductRepository, repo.StockHistoryRepository),
			StockHistoryService: NewStockHistorySrv(repo.StockHistoryRepository),
			StatisticsService:   NewStatisticsSrv(repo.StatisticsRepository),
			KitchenService:      NewKitchenSrv(repo.OrderRepository, repo.OrderItemRepository),
			PengeluaranService:  NewPengeluaranSrv(repo.PengeluaranRepository),
		}

	})

}
