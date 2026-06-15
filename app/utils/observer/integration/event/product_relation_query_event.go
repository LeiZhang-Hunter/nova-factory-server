package event

// ZProductRelationQueryReqEvent 商品对应关系查询请求事件，对应 emall.goodsrelation.get。
type ZProductRelationQueryReqEvent struct {
	Event
	Base
	ProductRelationQueryReqDataEvent
}

// ProductRelationQueryReqDataEvent 商品对应关系查询数据载体。
type ProductRelationQueryReqDataEvent interface {
	GetPage() int64
	GetPageSize() int64
	GetGoodsCode() *[]string
	GetGoodsName() *[]string
}
