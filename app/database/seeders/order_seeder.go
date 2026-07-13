package seeders

import (
	"backend-ta/app/constants"
	"backend-ta/app/domain"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func SeedOrders(ctx context.Context, db *bun.DB) error {
	count, err := db.NewSelect().Model((*domain.Order)(nil)).Count(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		fmt.Println("Orders table already has data, skipping...")
		return nil
	}

	var staffs []domain.User
	err = db.NewSelect().Model(&staffs).Where("role = ?", constants.UserRoleStaff).Scan(ctx)
	if err != nil {
		return fmt.Errorf("staffs not found: %w", err)
	}

	staffMap := make(map[int64]domain.User)
	for _, s := range staffs {
		if s.StoreID != nil {
			staffMap[*s.StoreID] = s
		}
	}

	var allProducts []domain.Product
	err = db.NewSelect().Model(&allProducts).Scan(ctx)
	if err != nil || len(allProducts) == 0 {
		return fmt.Errorf("products not found: %w", err)
	}

	productMap := make(map[int64][]domain.Product)
	for _, p := range allProducts {
		productMap[p.StoreID] = append(productMap[p.StoreID], p)
	}

	startDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2026, 5, 31, 23, 59, 59, 0, time.Local)

	var orders []domain.Order
	var orderItems []domain.OrderItem
	var payments []domain.Payment
	var stockHistories []domain.StockHistory

	indonesianNames := []string{
		"Ahmad", "Budi", "Siti", "Dewi", "Andi", "Joko", "Rian", "Indra", "Agus", "Rina",
		"Sari", "Aditya", "Hendra", "Mega", "Dian", "Putra", "Putri", "Tri", "Sri", "Wahyu",
		"Eko", "Bambang", "Heri", "Lilis", "Yanto", "Yanti", "Rudi", "Tari", "Tono", "Wati",
		"Doni", "Dina", "Hadi", "Gita", "Fajar", "Fitri", "Aris", "Anisa", "Roni", "Susi",
	}

	totalSteps := 120
	stepCount := 0

	for storeID := int64(1); storeID <= 1; storeID++ {
		staff, ok := staffMap[storeID]
		if !ok {
			continue
		}
		products, ok := productMap[storeID]
		if !ok || len(products) == 0 {
			continue
		}

		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			stepCount++
			if stepCount%48 == 0 {
				percent := (stepCount * 100) / totalSteps
				fmt.Printf("[SEEDER] Seeding Orders: %d%% completed (%d orders generated so far)\n", percent, len(orders))
			}

			// Generate 5 to 15 orders daily per store
			numOrders := rand.Intn(11) + 5

			for i := 0; i < numOrders; i++ {
				hour := rand.Intn(14) + 8
				minute := rand.Intn(60)
				orderTime := time.Date(d.Year(), d.Month(), d.Day(), hour, minute, 0, 0, time.Local)

				orderID := uuid.New().String()
				orderCode := fmt.Sprintf("%s-%03d", orderTime.Format("20060102"), i+1)

				status := constants.OrderStatusCompleted
				if rand.Intn(100) < 10 {
					status = constants.OrderStatusCancelled
				}

				numItems := rand.Intn(4) + 1
				var totalAmount float64

				for j := 0; j < numItems; j++ {
					prod := products[rand.Intn(len(products))]
					qty := rand.Intn(3) + 1

					subtotal := float64(qty) * prod.Price
					totalAmount += subtotal

					orderItemID := uuid.New().String()

					servedQty := qty
					if status == constants.OrderStatusCancelled {
						servedQty = 0
					}

					orderItems = append(orderItems, domain.OrderItem{
						ID:        orderItemID,
						OrderID:   orderID,
						ProductID: prod.ID,
						Quantity:  qty,
						UnitPrice: prod.Price,
						Subtotal:  subtotal,
						ServedQty: servedQty,
						CreatedAt: orderTime,
						UpdatedAt: orderTime,
					})

					if status == constants.OrderStatusCompleted {
						stockHistories = append(stockHistories, domain.StockHistory{
							ProductID: prod.ID,
							Change:    -qty,
							Reason:    fmt.Sprintf("Order %s Created", orderCode),
							CreatedAt: orderTime,
						})
					}
				}

				finalAmount := totalAmount

				customerName := indonesianNames[rand.Intn(len(indonesianNames))]

				orders = append(orders, domain.Order{
					ID:           orderID,
					OrderCode:    orderCode,
					StoreID:      storeID,
					CustomerName: &customerName,
					StaffID:      staff.ID,
					TotalAmount:  totalAmount,
					Status:       status,
					CreatedAt:    orderTime,
					UpdatedAt:    orderTime,
				})

				if status == constants.OrderStatusCompleted {
					paymentMethods := []constants.PaymentMethod{constants.PaymentMethodCash, constants.PaymentMethodDigitalWallet, constants.PaymentMethodCard}
					paymentMethod := paymentMethods[rand.Intn(len(paymentMethods))]

					payments = append(payments, domain.Payment{
						ID:            uuid.New().String(),
						OrderID:       orderID,
						PaymentMethod: paymentMethod,
						AmountPaid:    finalAmount,
						Timestamp:     orderTime.Add(time.Minute * 2),
					})
				}
			}
		}
	}

	batchSize := 100

	for i := 0; i < len(orders); i += batchSize {
		end := i + batchSize
		if end > len(orders) {
			end = len(orders)
		}
		batch := orders[i:end]
		_, err = db.NewInsert().Model(&batch).Exec(ctx)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(orderItems); i += batchSize {
		end := i + batchSize
		if end > len(orderItems) {
			end = len(orderItems)
		}
		batch := orderItems[i:end]
		_, err = db.NewInsert().Model(&batch).Exec(ctx)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(payments); i += batchSize {
		end := i + batchSize
		if end > len(payments) {
			end = len(payments)
		}
		batch := payments[i:end]
		_, err = db.NewInsert().Model(&batch).Exec(ctx)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(stockHistories); i += batchSize {
		end := i + batchSize
		if end > len(stockHistories) {
			end = len(stockHistories)
		}
		batch := stockHistories[i:end]
		_, err = db.NewInsert().Model(&batch).Exec(ctx)
		if err != nil {
			return err
		}
	}

	fmt.Printf("[SEEDER] Orders seeded successfully: %d orders created\n", len(orders))
	return nil
}
