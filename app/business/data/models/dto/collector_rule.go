package dto

// BindCollectorRulesRequest 下发规则请求
type BindCollectorRulesRequest struct {
	RuleIDs []string `json:"rule_ids" binding:"required,min=1"`
}

// PreviewRuleItem 预览中的单条规则
type PreviewRuleItem struct {
	RuleID     string `json:"rule_id"`
	RuleName   string `json:"rule_name"`
	SourceType string `json:"source_type"`
	Status     string `json:"status"`
	Version    int    `json:"version"`
}

// DispatchPreviewResponse 下发预览响应（config 为 YAML 格式）
type DispatchPreviewResponse struct {
	CollectorID string            `json:"collector_id"`
	Rules       []PreviewRuleItem `json:"rules"`
	Checksum    string            `json:"checksum"`
	Config      string            `json:"config"`
}

// PullRulesResponse 采集器拉取响应（config 为 YAML 格式）
type PullRulesResponse struct {
	DeviceID string `json:"device_id"`
	Checksum string `json:"checksum"`
	Updated  bool   `json:"updated"`
	Config   string `json:"config"`
}
