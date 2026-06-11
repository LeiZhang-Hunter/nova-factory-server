package event

type StockData interface {
	ProductID() int64
	SkuID() int64
	WarehouseID() int64
	Quantity() float64
	AfterQty() float64
	Base
}

// StockEvent 库存变更事件
type StockEvent interface {
	// Stocks 库存列表
	Stocks() []StockData
	Event
	Base
}
