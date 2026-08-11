package dto

import (
	"time"
)

// DispatchLogResponse 下发版本记录项（config 为合并后的 YAML 配置快照）
type DispatchLogResponse struct {
	ID            string    `json:"id"`
	CollectorID   string    `json:"collector_id"`
	CollectorName string    `json:"collector_name"`
	RuleIDs       []string  `json:"rule_ids"`
	RuleNames     []string  `json:"rule_names"`
	DispatchNo    int       `json:"dispatch_no"`
	Config        string    `json:"config"`
	ConfigMD5     string    `json:"config_md5"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
}
