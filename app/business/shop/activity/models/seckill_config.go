package models

import "nova-factory-server/app/baize"

// SeckillConfig 商品秒杀配置。
type SeckillConfig struct {
	ID            int64  `json:"id,string" gorm:"id"`                 // 秒杀配置ID
	BeginClock    int64  `json:"beginClock" gorm:"begin_clock"`       // 开启时间
	ContinueClock int64  `json:"continueClock" gorm:"continue_clock"` // 持续时间
	Images        string `json:"images" gorm:"images"`                // 轮播图
	Sort          int64  `json:"sort" gorm:"sort"`                    // 排序值
	Status        bool   `json:"status" gorm:"status"`
	DeptID        int64  `json:"deptId" gorm:"dept_id"` // 部门ID
	baize.BaseEntity
	State int32 `json:"state" gorm:"state"` // 操作状态
}

// SeckillConfigSet 商品秒杀配置保存参数。
type SeckillConfigSet struct {
	ID            int64  `json:"id,string"`     // 秒杀配置ID，更新时传
	BeginClock    int64  `json:"beginClock"`    // 开启时间
	ContinueClock int64  `json:"continueClock"` // 持续时间
	Images        string `json:"images"`        // 轮播图
	Sort          int64  `json:"sort"`          // 排序值
	Status        bool   `json:"status"`
}

// SeckillConfigQuery 商品秒杀配置查询参数。
type SeckillConfigQuery struct {
	BeginClock *int64 `form:"beginClock"` // 开启时间
	Status     *bool  `form:"status"`
	Page       int64  `form:"page"` // 页码
	Size       int64  `form:"size"` // 每页数量
}

// SeckillConfigListData 商品秒杀配置列表结果。
type SeckillConfigListData struct {
	Rows  []*SeckillConfig `json:"rows"`  // 列表数据
	Total int64            `json:"total"` // 总数
}
