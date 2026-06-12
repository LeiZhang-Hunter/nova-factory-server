// 定义库存事件相关的数据类型与接口。
// 包含库存数据（StockData）及库存变更事件（StockEvent），
// 用于将网店库存变化同步至第三方 ERP/全渠道系统。
package event

// StockData 库存数据接口，描述某一商品/SKU在特定仓库的库存信息。
// 包含变更前后的数量，便于 ERP 系统进行增量或全量同步。
type StockData interface {
	// ProductID 商品ID
	ProductID() int64
	// SkuID SKU ID，无SKU时为0
	SkuID() int64
	// WarehouseID 仓库ID，标识库存在哪个仓库
	WarehouseID() int64
	// Quantity 变更前库存数量
	Quantity() float64
	// AfterQty 变更后库存数量
	AfterQty() float64

	Base
}

// StockEvent 库存变更事件接口，表示一次库存数量变化事件。
// 聚合了 Event（事件元信息）和 Base（基础能力），并提供库存数据列表。
type StockEvent interface {
	// Stocks 返回本次事件涉及的库存数据列表
	Stocks() []StockData
	Event
	Base
}
