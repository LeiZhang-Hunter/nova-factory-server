package models

import "nova-factory-server/app/baize"

// SeckillActivity 秒杀活动
type SeckillActivity struct {
	ID           int64                         `json:"id,string" db:"id"`               // 活动ID
	Type         int8                          `json:"type" db:"type"`                  // 活动类型，1 表示秒杀
	Title        string                        `json:"title" db:"title"`                // 活动名称
	StartDay     int32                         `json:"startDay" db:"start_day"`         // 开始日期
	EndDay       int32                         `json:"endDay" db:"end_day"`             // 结束日期
	TimeIDs      string                        `json:"timeIds" db:"time_ids"`           // 时间段ID，多个逗号分隔
	OnceNum      int64                         `json:"onceNum" db:"once_num"`           // 活动期间每人每日购买数量，0 不限制
	Num          int64                         `json:"num" db:"num"`                    // 全部活动期间用户购买总数限制，0 不限制
	IsCommission int32                         `json:"isCommission" db:"is_commission"` // 是否参与分佣
	Status       int8                          `json:"status" db:"status"`              // 是否显示
	LinkID       int32                         `json:"linkId" db:"link_id"`             // 关联ID
	AddTime      int64                         `json:"addTime" db:"add_time"`           // 添加时间
	DeptID       int64                         `json:"deptId" db:"dept_id"`             // 部门ID
	ProductInfos []*SeckillActivityProductInfo `json:"productInfos" gorm:"-" db:"-"`    // 参与活动的商品列表
	baize.BaseEntity
	State int32 `json:"state" db:"state"` // 操作状态
}

// SeckillActivitySet 秒杀活动保存参数
type SeckillActivitySet struct {
	ID           int64                         `json:"id,string"`                   // 活动ID，更新时传
	Type         int8                          `json:"type"`                        // 活动类型
	Title        string                        `json:"title" binding:"required"`    // 活动名称
	StartDay     int32                         `json:"startDay" binding:"required"` // 开始日期
	EndDay       int32                         `json:"endDay" binding:"required"`   // 结束日期
	TimeIDs      []string                      `json:"timeIds" binding:"required"`  // 时间段ID，多个逗号分隔
	OnceNum      int64                         `json:"onceNum"`                     // 每人每日购买数量限制
	Num          int64                         `json:"num"`                         // 活动期间总购买数量限制
	IsCommission int32                         `json:"isCommission"`                // 是否参与分佣
	Status       int8                          `json:"status"`                      // 是否显示
	LinkID       int32                         `json:"linkId"`                      // 关联ID，保留字段
	ProductInfos []*SeckillActivityProductInfo `json:"productInfos"`                // 参与活动的商品列表
}

// SeckillActivityProductInfo 秒杀活动中的商品信息。
type SeckillActivityProductInfo struct {
	ID     int64                             `json:"id,string"` // 商品主键ID
	Status int32                             `json:"status"`    // 商品状态
	Sort   int64                             `json:"sort"`      // 排序值
	IsHot  int32                             `json:"isHot"`     // 是否热门推荐
	Attrs  []*SeckillActivityProductInfoAttr `json:"attrs"`     // 商品规格活动信息
}

// SeckillActivityProductInfoAttr 秒杀活动中的商品规格活动信息。
type SeckillActivityProductInfoAttr struct {
	SkuID   string  `json:"skuId"`   // 规格业务ID
	Status  int32   `json:"status"`  // 是否参与活动
	Price   float64 `json:"price"`   // 活动价
	Cost    float64 `json:"cost"`    // 成本价
	OtPrice float64 `json:"otPrice"` // 原价
	Quota   int64   `json:"quota"`   // 限量
}

// SeckillActivityQuery 秒杀活动查询参数
type SeckillActivityQuery struct {
	Title  string `form:"title"`  // 活动名称
	Type   *int8  `form:"type"`   // 活动类型
	Status *int8  `form:"status"` // 显示状态
	LinkID int32  `form:"linkId"` // 关联ID
	Page   int64  `form:"page"`   // 页码
	Size   int64  `form:"size"`   // 每页数量
}

// SeckillActivityListData 秒杀活动列表数据
type SeckillActivityListData struct {
	Rows  []*SeckillActivity `json:"rows"`  // 列表数据
	Total int64              `json:"total"` // 总数
}
