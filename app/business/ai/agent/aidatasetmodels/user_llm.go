package aidatasetmodels

type SysUserLLM struct {
	UserID     int64  `gorm:"column:user_id;primaryKey" json:"user_id,string"`
	LLMFactory string `gorm:"column:llm_factory;primaryKey" json:"llm_factory"`
	ModelType  string `gorm:"column:model_type" json:"model_type"`
	LLMName    string `gorm:"column:llm_name;primaryKey" json:"llm_name"`
	APIType    string `gorm:"column:api_type" json:"api_type"`
	APIKey     string `gorm:"column:api_key" json:"api_key"`
	APIBase    string `gorm:"column:api_base" json:"api_base"`
	MaxTokens  int64  `gorm:"column:max_tokens" json:"max_tokens"`
	UsedTokens int64  `gorm:"column:used_tokens" json:"used_tokens"`
	Status     string `gorm:"column:status" json:"status"`
}

func (e *SysUserLLM) GetUserID() int64 {
	return e.UserID
}
func (e *SysUserLLM) GetLLMFactory() string {
	return e.LLMFactory
}
func (e *SysUserLLM) GetModelType() string {
	return e.ModelType
}
func (e *SysUserLLM) GetLLMName() string {
	return e.LLMName
}
func (e *SysUserLLM) GetAPIType() string {
	return e.APIType
}
func (e *SysUserLLM) GetAPIKey() string {
	return e.APIKey
}
func (e *SysUserLLM) GetAPIBase() string {
	return e.APIBase
}
func (e *SysUserLLM) GetMaxTokens() int64 {
	return e.MaxTokens
}
func (e *SysUserLLM) GetUsedTokens() int64 {
	return e.UsedTokens
}
func (e *SysUserLLM) GetStatus() string {
	return e.Status
}

type SetSysUserLLM struct {
	LLMFactory string `json:"llm_factory" binding:"max=128"`
	APIKey     string `json:"api_key"`
	APIBase    string `json:"api_base" binding:"max=255"`
	APIType    string `json:"api_type" binding:"max=128"`
}

type GetSysUserLLMReq struct {
	LLMFactory string `form:"llm_factory"`
}
