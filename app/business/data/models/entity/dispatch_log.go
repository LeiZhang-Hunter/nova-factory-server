package entity

import "time"

// DispatchLog 下发版本记录（一次下发 = 多选规则合并后的配置快照）
type DispatchLog struct {
	ID            string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	CollectorID   string    `gorm:"type:varchar(36);not null;index" json:"collector_id"`
	CollectorName string    `gorm:"type:varchar(255)" json:"collector_name"`
	RuleIDs       string    `gorm:"type:json" json:"-"`
	RuleNames     string    `gorm:"type:json" json:"-"`
	DispatchNo    int       `gorm:"type:int;not null" json:"dispatch_no"`
	Config        string    `gorm:"type:text;not null" json:"-"`
	ConfigMD5     string    `gorm:"type:char(32);not null" json:"config_md5"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `gorm:"type:varchar(100)" json:"created_by"`
}

func (DispatchLog) TableName() string { return "data_dispatch_logs" }
