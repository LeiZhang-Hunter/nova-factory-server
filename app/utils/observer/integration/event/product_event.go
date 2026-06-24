package event

// ProductSku 商品SKU（规格）接口，描述商品的单个规格变体信息。
// 一种商品可能有多个SKU，每个SKU对应不同颜色、尺码等属性组合。
type ProductSku interface {
	// GetGoodsId 所属商品ID
	GetGoodsId() int64
	// GetOuterId SKU的外部编码（商家自定义编码）
	GetOuterId() string
	// GetPrice SKU价格
	GetPrice() float64
	// GetQuantity SKU库存数量
	GetQuantity() int64
	// GetSkuId SKU唯一ID
	GetSkuId() int64
	// GetSkuName SKU名称（如 "红色-XL"）
	GetSkuName() string
	// GetSkuCode SKU编码
	GetSkuCode() string
	// GetBarcode SKU条码
	GetBarcode() string
	// GetWeight SKU重量
	GetWeight() float64
	// GetSize SKU尺寸
	GetSize() float64
}

// ProductData 商品数据接口，封装一个完整商品的所有字段。
// 包含商品基本信息及下属SKU列表，各业务模块实现此接口以接入商品同步。
type ProductData interface {
	// GetGoodsId 商品ID
	GetGoodsId() int64
	// GetGoodsName 商品名称
	GetGoodsName() string
	// GetGoodsCode 商品编码
	GetGoodsCode() string
	// GetImage 商品主图URL
	GetImage() string
	// GetDesc 商品描述
	GetDesc() string
	// GetUnitName 基本单位名称（如 "件"、"箱"）
	GetUnitName() string
	// GetIsService 是否为服务类商品（如虚拟商品、安装服务等）
	GetIsService() int
	// GetPrice 商品价格
	GetPrice() float64
	// GetQuantity 商品库存数量
	GetQuantity() float64
	// GetOuterId 商品外部编码
	GetOuterId() string
	// GetSkus 商品规格列表，无规格时返回空切片
	GetSkus() []ProductSku

	Base
}

// ProductEvent 商品变更事件接口，表示商品创建、更新或删除事件。
// 聚合了 Event（事件元信息）和 Base（基础能力），并提供商品数据列表。
type ProductEvent interface {
	Base
	Event
	// GetProducts 返回本次事件涉及的商品数据列表
	GetProducts() []ProductData
}
