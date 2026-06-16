package result

// ProductRelationQueryResponse 商品对应关系查询响应接口，对应 emall.goodsrelation.get 返回。
type ProductRelationQueryResponse interface {
	GetCode() int64
	GetMessage() string
	GetTotal() int64
	GetGoodsRelation() []ProductRelationQueryData
}

// ProductRelationQueryData 商品对应关系数据接口，对应 goodsrelation[{...}]。
type ProductRelationQueryData interface {
	GetEshopGoodsID() string
	GetEshopSkuID() string
	GetGoodsID() string
	GetGoodsCode() string
	GetGoodsName() string
	GetSkuID() int64
	GetSkuCode() string
	GetSkuName() string
	GetUnitID() int64
	GetUnitName() string
}
