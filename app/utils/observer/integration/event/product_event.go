package event

import "nova-factory-server/app/utils/observer/integration/config"

type ProductData interface {
	ProductID() int64
	SkuID() int64
	GoodsName() string
	Barcode() string
	Price() float64
	CategoryID() int64
	Config() config.Config
	Base
}

// ProductEvent 商品变更事件
type ProductEvent interface {
	Base
	Event
	Products() []ProductData
}
