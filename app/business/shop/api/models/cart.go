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
	GoodsID  int64 `json:"goodsId,string" binding:"required"` // 商品ID
	SkuID    int64 `json:"skuId,string" binding:"required"`   // SKU ID
	Quantity int64 `json:"quantity" binding:"required"`       // 数量
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
}

// CartDto 商城用户购物车项
type CartDto struct {
	ID          int64      `json:"id,string" db:"id"`                    // 主键ID
	UserID      int64      `json:"userId" db:"user_id"`                  // 用户ID
	GoodsID     int64      `json:"goodsId" db:"goods_id"`                // 商品ID
	SkuID       int64      `json:"skuId" db:"sku_id"`                    // SKU ID
	GoodsName   string     `json:"goodsName" db:"goods_name"`            // 商品名称快照
	SkuName     string     `json:"skuName" db:"sku_name"`                // SKU名称快照
	ImageURL    string     `json:"imageUrl" db:"image_url"`              // 商品或SKU图片快照
	RetailPrice float64    `json:"retailPrice" db:"retail_price"`        // 加入购物车时零售价快照
	Quantity    int64      `json:"quantity" db:"quantity"`               // 购买数量
	CreateTime  *time.Time `json:"createTime" db:"create_time"`          //创建时间
	UpdateTime  *time.Time `json:"updateTime" db:"update_time"`          //修改时间
	State       int32      `json:"state" gorm:"column:state" db:"state"` // 操作状态
}
