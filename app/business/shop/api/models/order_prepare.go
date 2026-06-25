package models

import "nova-factory-server/app/business/shop/activity/models"

// OrderCacheItem 是确认单和创建订单共享的商品快照。
type OrderCacheItem struct {
	CombinationId     int64                       `json:"combinationId,string"`
	SecKillId         int64                       `json:"secKillId,string"`
	SeckillInfo       *models.SeckillMainInfo     `json:"seckillInfo"`
	CombinationInfo   *models.CombinationMainInfo `json:"combinationMainInfo"`
	PinkId            int64                       `json:"pinkId,string"`
	OuterIid          string                      `json:"outerIid"`
	GoodsID           int64                       `json:"goodsId,string"`
	SkuID             int64                       `json:"skuId,string"`
	GoodsName         string                      `json:"goodsName"`
	SkuName           string                      `json:"skuName"`
	ImageURL          string                      `json:"imageUrl"`
	Price             float64                     `json:"price"`
	Quantity          int64                       `json:"quantity"`
	AvailableStock    int64                       `json:"availableStock"`
	StockInsufficient bool                        `json:"stockInsufficient"`
	TotalAmount       float64                     `json:"totalAmount"`
	CartID            int64                       `json:"cartId,string"`
}

// OrderCacheData 是确认单缓存，Confirm 写入，Create 消费。
type OrderCacheData struct {
	OrderKey       string            `json:"orderKey"`
	UserID         int64             `json:"userId,string"`
	Items          []*OrderCacheItem `json:"items"`
	AddressID      int64             `json:"addressId,string"`
	DeliveryType   string            `json:"deliveryType"`
	Remark         string            `json:"remark"`
	GoodsAmount    float64           `json:"goodsAmount"`
	FreightAmount  float64           `json:"freightAmount"`
	DiscountAmount float64           `json:"discountAmount"`
	PayAmount      float64           `json:"payAmount"`
	CartIDs        []int64           `json:"cartIds"`
	BuyNow         bool              `json:"buyNow"`
}

// OrderConfirmReq 是确认单请求。确认单从 cartId 构建，不再接收 orderKey。
type OrderConfirmReq struct {
	CartID       string `json:"cartId" binding:"required"`
	BuyNow       bool   `json:"buyNow"`
	AddressID    int64  `json:"addressId,string"`
	DeliveryType string `json:"deliveryType"`
	ShippingType int32  `json:"shippingType"`
}

// OrderConfirmResp 是确认单响应，包含后续创建订单所需的 orderKey。
type OrderConfirmResp struct {
	OrderKey       string              `json:"orderKey"`
	Address        *ShopUserAddressApp `json:"address"`
	Items          []*OrderCacheItem   `json:"items"`
	GoodsAmount    float64             `json:"goodsAmount"`
	FreightAmount  float64             `json:"freightAmount"`
	DiscountAmount float64             `json:"discountAmount"`
	PayAmount      float64             `json:"payAmount"`
	DeliveryType   string              `json:"deliveryType"`
	ExpireSeconds  int64               `json:"expireSeconds"`
}

func (r *OrderConfirmReq) CartIDValue() string {
	if r == nil {
		return ""
	}
	return r.CartID
}

// OrderCreateReq 是正式创建订单请求，仅保留主链路字段。
type OrderCreateReq struct {
	OrderKey string `json:"orderKey" binding:"required"`
	Remark   string `json:"remark"`
}
