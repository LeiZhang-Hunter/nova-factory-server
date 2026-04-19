package models

import "nova-factory-server/app/baize"

// Cart 商城用户购物车项
type Cart struct {
	ID          int64   `json:"id,string" db:"id"`                         // 主键ID
	UserID      int64   `json:"userId" db:"user_id"`                       // 用户ID
	GoodsID     string  `json:"goodsId" db:"goods_id"`                     // 商品ID
	SkuID       string  `json:"skuId" db:"sku_id"`                         // SKU ID
	GoodsName   string  `json:"goodsName" db:"goods_name"`                 // 商品名称快照
	SkuName     string  `json:"skuName" db:"sku_name"`                     // SKU名称快照
	ImageURL    string  `json:"imageUrl" db:"image_url"`                   // 商品或SKU图片快照
	RetailPrice float64 `json:"retailPrice" db:"retail_price"`             // 加入购物车时零售价快照
	Quantity    int64   `json:"quantity" db:"quantity"`                    // 购买数量
	Selected    int32   `json:"selected" db:"selected"`                    // 是否选中：1选中，0未选中
	Status      int32   `json:"status" db:"status"`                        // 状态：1有效，0失效
	DeptID      int64   `json:"deptId" gorm:"column:dept_id" db:"dept_id"` // 部门ID
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state" db:"state"` // 操作状态
}

// CartSetReq 购物车新增修改参数
type CartSetReq struct {
	ID          int64   `json:"id,string"` // 主键ID
	UserID      int64   `json:"-"`
	Username    string  `json:"username"`                     // 用户名
	GoodsID     string  `json:"goodsId" binding:"required"`   // 商品ID
	SkuID       string  `json:"skuId" binding:"required"`     // SKU ID
	GoodsName   string  `json:"goodsName" binding:"required"` // 商品名称
	SkuName     string  `json:"skuName"`                      // SKU名称
	ImageURL    string  `json:"imageUrl"`                     // 图片地址
	RetailPrice float64 `json:"retailPrice"`                  // 零售价
	Quantity    int64   `json:"quantity" binding:"required"`  // 数量
	Selected    *int32  `json:"selected"`                     // 是否选中
	Status      *int32  `json:"status"`                       // 状态
}

// CartQuery 购物车查询参数
type CartQuery struct {
	UserID   int64  `form:"userId"`   // 用户ID
	GoodsID  string `form:"goodsId"`  // 商品ID
	SkuID    string `form:"skuId"`    // SKU ID
	Selected *int32 `form:"selected"` // 是否选中
	Status   *int32 `form:"status"`   // 状态
	Page     int64  `form:"page"`     // 页码
	Size     int64  `form:"size"`     // 每页数量
}

// CartListData 购物车列表结果
type CartListData struct {
	Rows  []*Cart `json:"rows"`  // 数据列表
	Total int64   `json:"total"` // 总数
}
