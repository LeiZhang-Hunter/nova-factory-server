package models

// SeckillConfig 秒杀时间段配置
type SeckillConfig struct {
	ID            int64  `json:"id,string" gorm:"column:id"`
	BeginClock    int64  `json:"beginClock" gorm:"column:begin_clock"`
	ContinueClock int64  `json:"continueClock" gorm:"column:continue_clock"`
	Images        string `json:"images" gorm:"column:images"`
	Sort          int64  `json:"sort" gorm:"column:sort"`
	Status        bool   `json:"status" gorm:"column:status"`
}

// SeckillConfigQuery 秒杀配置查询参数
type SeckillConfigQuery struct {
	Status *bool `form:"status"`
	Page   int64 `form:"page"`
	Size   int64 `form:"size"`
}

// SeckillConfigListData 秒杀配置列表数据
type SeckillConfigListData struct {
	Rows  []*SeckillConfig `json:"rows"`
	Total int64            `json:"total"`
}
