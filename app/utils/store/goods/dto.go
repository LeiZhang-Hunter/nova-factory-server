package goods

// Request 商品查询
type Request interface {
	// GetSellerCIds 卖家店铺自定义类目ID
	//（多个之间用“,”分隔，可以根据 selfmall.sellercats.list.get 获得;目前最多支持32个。
	//类目ID为空或0时，返回店铺所有商品信息）
	GetSellerCIds() string

	// GetStartModified 起始的修改时间
	GetStartModified() string

	// GetEndModified 结束的修改时间
	GetEndModified() string

	// GetName 商品名称（搜索商品的名称）
	GetName() string

	// GetOuterId 规格商家编码（搜索商品的商家编码）
	GetOuterId() string

	// GetStatus 	查询商品状态（up=线上，down=库存中）
	GetStatus() string

	// GetPageNo 页码
	GetPageNo() int

	// GetPageSize 	每页条数
	GetPageSize() int
}

type ProductData interface {
	GetCid() int
	GetCatName() string
	GetProductId() int64
	GetName() string
	GetOuterId() string
	GetPicPath() string
	GetPrice() int
	GetBarcodeStr() string
	GetCreated() string
	GetDesc() string
	GetModified() string
	GetStatus() string
	GetQuantity() int
	GetSkus() []ProductSku
}

type ProductSku interface {
	GetSkuId() int64
	GetSkuName() string
	GetProductId() int
	GetOuterId() string
	GetPrice() int
	GetQuantity() int
	GetCreated() string
	GetModified() string
	GetProperties() string
}

// DataResult 查询商品结果
type DataResult interface {
	GetIsError() bool
	GetErrorMsg() string
	GetTotalResults() int
	GetProductInfo() []ProductData
}

// EmptyDataResult 查询空的商品结果
type EmptyDataResult struct{}

func (*EmptyDataResult) GetIsError() bool {
	return false
}

func (*EmptyDataResult) GetErrorMsg() string {
	return ""
}

func (*EmptyDataResult) GetTotalResults() int {
	return 0
}

func (*EmptyDataResult) GetProductInfo() []ProductData {
	return []ProductData{}
}
