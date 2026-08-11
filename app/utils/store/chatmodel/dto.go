package chatmodel

// LlmConnection AI 模块中某个模型（供应商 + 模型名）的连接信息。
type LlmConnection interface {
	GetAPIType() string
	GetAPIKey() string
	GetAPIBase() string
	GetMaxTokens() int64
}

// LlmConnectionData LlmConnection 的默认实现。
type LlmConnectionData struct {
	APIType   string
	APIKey    string
	APIBase   string
	MaxTokens int64
}

func (d *LlmConnectionData) GetAPIType() string {
	return d.APIType
}

func (d *LlmConnectionData) GetAPIKey() string {
	return d.APIKey
}

func (d *LlmConnectionData) GetAPIBase() string {
	return d.APIBase
}

func (d *LlmConnectionData) GetMaxTokens() int64 {
	return d.MaxTokens
}
