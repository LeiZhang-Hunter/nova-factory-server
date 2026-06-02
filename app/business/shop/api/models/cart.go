package models

import (
	"time"
)

type CartItem struct {
}

// CartQuery 购物车查询参数
type CartQuery struct {
	Page int64 `form:"page"`
	Size int64 `form:"size"`
}

// CartListData 购物车列表结果
type CartListData struct {
	Rows  []*CartDto `json:"rows"`
	Total int64      `json:"total"`
}

// CartSetDataReq 购物车新增修改参数
type CartSetDataReq struct {
	GoodsID       int64 `json:"goodsId,string"`       // 商品ID
	SkuID         int64 `json:"skuId,string"`         // SKU ID
	Quantity      int64 `json:"quantity"`             // 数量
	BuyNow        bool  `json:"buyNow"`               // 是否立即购买
	SecKillID     int64 `json:"secKillId,string"`     // 秒杀活动商品ID
	CombinationID int64 `json:"combinationId,string"` // 拼团活动商品ID
	PinkID        int64 `json:"pinkId,string"`        // 拼团记录ID
}

func (r *CartSetDataReq) IsBuyNow() bool {
	if r == nil {
		return false
	}
	return r.BuyNow
}

type CartSetData struct {
	ID          int64   `json:"id,string"` // 主键ID
	UserID      int64   `json:"-"`
	Username    string  `json:"username"`                     // 用户名
	GoodsID     int64   `json:"goodsId" binding:"required"`   // 商品ID
	SkuID       uint64  `json:"skuId" binding:"required"`     // SKU ID
	GoodsName   string  `json:"goodsName" binding:"required"` // 商品名称
	SkuName     string  `json:"skuName"`                      // SKU名称
	ImageURL    string  `json:"imageUrl"`                     // 图片地址
	RetailPrice float64 `json:"retailPrice"`                  // 零售价
	Quantity    int64   `json:"quantity" binding:"required"`  // 数量
	ProductType int32   `json:"productType"`                  // 商品类型：0普通 1秒杀 3拼团
	ActivityID  int64   `json:"activityId,string"`            // 活动商品ID
	PinkID      int64   `json:"pinkId,string"`                // 拼团记录ID
	State       int32   `json:"state"`                        // 购物车状态
}

type CartSetDataResp struct {
	Mode   string   `json:"mode"`             // cart 或 buyNow
	CartID string   `json:"cartId,omitempty"` // 购物车ID或立即购买临时快照key
	Cart   *CartDto `json:"cart,omitempty"`   // 购物车项
}

// CartDto 商城用户购物车项
type CartDto struct {
	ID            int64      `json:"id,string" gorm:"id"`                    // 主键ID
	UserID        int64      `json:"userId" gorm:"user_id"`                  // 用户ID
	GoodsID       int64      `json:"goodsId" gorm:"goods_id"`                // 商品ID
	SkuID         int64      `json:"skuId" gorm:"sku_id"`                    // SKU ID
	GoodsName     string     `json:"goodsName" gorm:"goods_name"`            // 商品名称快照
	SkuName       string     `json:"skuName" gorm:"sku_name"`                // SKU名称快照
	ImageURL      string     `json:"imageUrl" gorm:"image_url"`              // 商品或SKU图片快照
	RetailPrice   float64    `json:"retailPrice" gorm:"retail_price"`        // 加入购物车时零售价快照
	Quantity      int64      `json:"quantity" gorm:"quantity"`               // 购买数量
	ProductType   int32      `json:"productType" gorm:"column:product_type"` // 商品类型：0普通 1秒杀 3拼团
	ActivityID    int64      `json:"activityId,string" gorm:"activity_id"`   // 活动商品ID
	PinkID        int64      `json:"pinkId,string" gorm:"pink_id"`           // 拼团记录ID
	IsStockEnough bool       `json:"isStockEnough" gorm:"-" gorm:"-"`        // 库存是否充足（当前库存 >= 购买数量）
	CreateTime    *time.Time `json:"createTime" gorm:"create_time"`          //创建时间
	UpdateTime    *time.Time `json:"updateTime" gorm:"update_time"`          //修改时间
	State         int32      `json:"state" gorm:"column:state" gorm:"state"` // 操作状态
}
