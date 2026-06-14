package models

import (
	"gorm.io/gorm"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"
)

type StockSyncReq struct {
	Stocks   []StockItem    `json:"stocks" form:"stocks"`
	Cache    cache.Cache    `json:"-"`
	Callback event.Callback `json:"-"`
	DB       *gorm.DB       `json:"-"`
}

func (s *StockSyncReq) GetDB() *gorm.DB {
	return s.DB
}
func (s *StockSyncReq) WithDB(tx *gorm.DB) {
	s.DB = tx
}
func (s *StockSyncReq) ToEvent() event.StockEvent {
	return s
}

func (s *StockSyncReq) GetStocks() []event.StockData {
	stocks := make([]event.StockData, 0, len(s.Stocks))
	for i := range s.Stocks {
		stocks = append(stocks, &s.Stocks[i])
	}
	return stocks
}

func (s *StockSyncReq) Metadata() map[string]any {
	return map[string]any{}
}

func (s *StockSyncReq) Ptr() any {
	return s
}

func (s *StockSyncReq) Config() config.Config {
	return nil
}

func (s *StockSyncReq) Action() event.EventType {
	return event.EventStockChanged
}

func (s *StockSyncReq) GetCache() cache.Cache {
	return s.Cache
}

func (s *StockSyncReq) GetCallback() event.Callback {
	return s.Callback
}

// StockItem 单条库存变更数据，实现 event.StockData 接口
type StockItem struct {
	ProductId     int64   `json:"productid" form:"productid"`
	SkuId         int64   `json:"skuid" form:"skuid"`
	WarehouseId   int64   `json:"warehouseid" form:"warehouseid"`
	Qty           float64 `json:"quantity" form:"quantity"`
	AfterQuantity float64 `json:"afterqty" form:"afterqty"`
}

func (s *StockItem) ProductID() int64 {
	return s.ProductId
}

func (s *StockItem) SkuID() int64 {
	return s.SkuId
}

func (s *StockItem) WarehouseID() int64 {
	return s.WarehouseId
}

func (s *StockItem) Quantity() float64 {
	return s.Qty
}

func (s *StockItem) AfterQty() float64 {
	return s.AfterQuantity
}

func (s *StockItem) Metadata() map[string]any {
	return map[string]any{}
}

func (s *StockItem) Ptr() any {
	return s
}
