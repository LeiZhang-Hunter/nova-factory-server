package models

import "nova-factory-server/app/business/shop/activity/models"

// OrderCacheReq 预下单缓存请求
type OrderCacheReq struct {
	CombinationId int64    `json:"combinationId,string"` // 拼团活动ID
	SecKillId     int64    `json:"secKillId,string"`     // 秒杀活动ID
	PinkId        int64    `json:"pinkId,string"`        // 拼团记录ID
	SkuID         int64    `json:"skuId,string"`         // 立即购买商品规格ID
	Quantity      int64    `json:"quantity"`             // 购买数量
	CartId        []string `json:"cartId"`               // 购物车ID列表
}

// OrderCacheItemReq 预下单商品请求
type OrderCacheItemReq struct {
	GoodsID  int64 `json:"goodsId,string" binding:"required"` // 商品ID
	SkuID    int64 `json:"skuId,string" binding:"required"`   // SKU ID
	Quantity int64 `json:"quantity" binding:"required"`       // 购买数量

}

// OrderCacheItem 预下单商品快照
type OrderCacheItem struct {
	CombinationId     int64                       `json:"combinationId,string"`
	SecKillId         int64                       `json:"secKillId,string"`
	SeckillInfo       *models.SeckillMainInfo     `json:"seckillInfo"`
	CombinationInfo   *models.CombinationMainInfo `json:"combinationMainInfo"`
	PinkId            int64                       `json:"pinkId,string"`
	GoodsID           int64                       `json:"goodsId,string"`    // 商品ID
	SkuID             int64                       `json:"skuId,string"`      // SKU ID
	GoodsName         string                      `json:"goodsName"`         // 商品名称
	SkuName           string                      `json:"skuName"`           // SKU 名称
	ImageURL          string                      `json:"imageUrl"`          // 图片地址
	Price             float64                     `json:"price"`             // 当前单价快照
	Quantity          int64                       `json:"quantity"`          // 购买数量
	AvailableStock    int64                       `json:"availableStock"`    // 当前可用库存
	StockInsufficient bool                        `json:"stockInsufficient"` // 是否库存不足
	TotalAmount       float64                     `json:"totalAmount"`       // 商品小计
}

// OrderCacheData 预下单缓存数据
type OrderCacheData struct {
	OrderKey       string            `json:"orderKey"`         // 预订单 key
	UserID         int64             `json:"userId,string"`    // 用户ID
	Items          []*OrderCacheItem `json:"items"`            // 商品快照列表
	AddressID      int64             `json:"addressId,string"` // 地址ID
	DeliveryType   string            `json:"deliveryType"`     // 配送方式
	Remark         string            `json:"remark"`           // 备注
	GoodsAmount    float64           `json:"goodsAmount"`      // 商品金额
	FreightAmount  float64           `json:"freightAmount"`    // 运费金额
	DiscountAmount float64           `json:"discountAmount"`   // 优惠金额
	PayAmount      float64           `json:"payAmount"`        // 应付金额
}

// OrderCacheResp 预下单缓存响应
type OrderCacheResp struct {
	OrderKey      string            `json:"orderKey"`      // 预订单 key
	Items         []*OrderCacheItem `json:"items"`         // 商品快照列表
	ExpireSeconds int64             `json:"expireSeconds"` // 过期秒数
}

// OrderConfirmReq 确认单请求
type OrderConfirmReq struct {
	OrderKey     string `json:"orderKey" binding:"required"` // 预订单 key
	AddressID    int64  `json:"addressId,string"`            // 地址ID
	CouponCode   string `json:"couponCode"`                  // 优惠券编码
	UsePoints    bool   `json:"usePoints"`                   // 是否使用积分
	DeliveryType string `json:"deliveryType"`                // 配送方式
}

// OrderConfirmResp 确认单响应
type OrderConfirmResp struct {
	OrderKey       string              `json:"orderKey"`       // 预订单 key
	Address        *ShopUserAddressApp `json:"address"`        // 收货地址
	Items          []*OrderCacheItem   `json:"items"`          // 商品快照列表
	GoodsAmount    float64             `json:"goodsAmount"`    // 商品金额
	FreightAmount  float64             `json:"freightAmount"`  // 运费金额
	DiscountAmount float64             `json:"discountAmount"` // 优惠金额
	PayAmount      float64             `json:"payAmount"`      // 应付金额
	CouponCode     string              `json:"couponCode"`     // 优惠券编码
	UsePoints      bool                `json:"usePoints"`      // 是否使用积分
	DeliveryType   string              `json:"deliveryType"`   // 配送方式
	ExpireSeconds  int64               `json:"expireSeconds"`  // 过期秒数
}

// OrderCreateReq 正式创建订单请求
type OrderCreateReq struct {
	OrderKey string `json:"orderKey" binding:"required"` // 预订单 key
	Remark   string `json:"remark"`                      // 订单备注
}
