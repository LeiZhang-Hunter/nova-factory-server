package aiDataSetModels

type SysUserLLM struct {
	UserID     int64  `gorm:"column:user_id;primaryKey" json:"user_id,string"`
	LLMFactory string `gorm:"column:llm_factory;primaryKey" json:"llm_factory"`
	ModelType  string `gorm:"column:model_type" json:"model_type"`
	LLMName    string `gorm:"column:llm_name;primaryKey" json:"llm_name"`
	APIKey     string `gorm:"column:api_key" json:"api_key"`
	APIBase    string `gorm:"column:api_base" json:"api_base"`
	MaxTokens  int64  `gorm:"column:max_tokens" json:"max_tokens"`
	UsedTokens int64  `gorm:"column:used_tokens" json:"used_tokens"`
	Status     string `gorm:"column:status" json:"status"`
}

type SetSysUserLLM struct {
	LLMFactory string `json:"llm_factory"`
	APIKey     string `json:"api_key"`
	APIBase    string `json:"api_base"`
}

type GetSysUserLLMReq struct {
	LLMFactory string `form:"llm_factory"`
}
