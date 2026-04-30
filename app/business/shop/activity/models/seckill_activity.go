package models

import "nova-factory-server/app/baize"

// SeckillActivity 秒杀活动
type SeckillActivity struct {
	ID           int64  `json:"id,string" db:"id"`               // 活动ID
	Type         int32  `json:"type" db:"type"`                  // 活动类型，1秒杀
	Title        string `json:"title" db:"title"`                // 活动名称
	StartDay     int64  `json:"startDay" db:"start_day"`         // 开始日期
	EndDay       int64  `json:"endDay" db:"end_day"`             // 结束日期
	TimeIDs      string `json:"timeIds" db:"time_ids"`           // 时间段ID，多个逗号分隔
	OnceNum      int32  `json:"onceNum" db:"once_num"`           // 每人每日购买数量限制
	Num          int32  `json:"num" db:"num"`                    // 活动期间总购买数量限制
	IsCommission int32  `json:"isCommission" db:"is_commission"` // 是否参与分佣
	Status       int32  `json:"status" db:"status"`              // 是否显示
	LinkID       int32  `json:"linkId" db:"link_id"`             // 关联ID
	AddTime      int64  `json:"addTime" db:"add_time"`           // 添加时间
	DeptID       int64  `json:"deptId" db:"dept_id"`             // 部门ID
	baize.BaseEntity
	State int32 `json:"state" db:"state"` // 操作状态
}

// SeckillActivitySet 秒杀活动保存参数
type SeckillActivitySet struct {
	ID           int64  `json:"id,string"`                   // 活动ID，更新时传
	Type         int32  `json:"type"`                        // 活动类型
	Title        string `json:"title" binding:"required"`    // 活动名称
	StartDay     int64  `json:"startDay" binding:"required"` // 开始日期
	EndDay       int64  `json:"endDay" binding:"required"`   // 结束日期
	TimeIDs      string `json:"timeIds" binding:"required"`  // 时间段ID，多个逗号分隔
	OnceNum      int32  `json:"onceNum"`                     // 每人每日购买数量限制
	Num          int32  `json:"num"`                         // 活动期间总购买数量限制
	IsCommission int32  `json:"isCommission"`                // 是否参与分佣
	Status       int32  `json:"status"`                      // 是否显示
	LinkID       int32  `json:"linkId"`                      // 关联ID
}

// SeckillActivityQuery 秒杀活动查询参数
type SeckillActivityQuery struct {
	Title  string `form:"title"`  // 活动名称
	Type   *int32 `form:"type"`   // 活动类型
	Status *int32 `form:"status"` // 显示状态
	LinkID int32  `form:"linkId"` // 关联ID
	Page   int64  `form:"page"`   // 页码
	Size   int64  `form:"size"`   // 每页数量
}

// SeckillActivityListData 秒杀活动列表数据
type SeckillActivityListData struct {
	Rows  []*SeckillActivity `json:"rows"`  // 列表数据
	Total int64              `json:"total"` // 总数
}
