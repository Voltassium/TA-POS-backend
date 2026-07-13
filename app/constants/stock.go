package constants

type StockSourceType = string

const (
	StockSourcePurchase StockSourceType = "purchase"
	StockSourceSale     StockSourceType = "sale"
	StockSourceReturn   StockSourceType = "return"
	StockSourceManual   StockSourceType = "manual"
)
