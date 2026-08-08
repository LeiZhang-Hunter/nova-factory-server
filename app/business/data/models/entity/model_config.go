package entity

import "time"

const (
	// ModelConfigDefaultID 数据平台模型配置固定使用单行记录。
	ModelConfigDefaultID = "default"

	// 默认模型参数（与管理后台智能体「模型参数」一致）。
	DefaultModelTemperature    = 0.7
	DefaultModelTopP           = 0.9
	DefaultModelMaxTokens      = 4096
	DefaultModelMaxContext     = 10
	DefaultModelRetrievalTopK  = 5
	DefaultModelMatchThreshold = 0.5
)

// ModelConfig 数据平台模型配置（全局单行）。
type ModelConfig struct {
	ID                      string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Provider                string    `gorm:"type:varchar(100)" json:"provider"`
	Model                   string    `gorm:"type:varchar(200)" json:"model"`
	Temperature             float64   `gorm:"type:decimal(4,2);default:0.7" json:"temperature"`
	EnableTemperature       bool      `gorm:"type:tinyint(1);default:0" json:"enable_temperature"`
	TopP                    float64   `gorm:"type:decimal(4,2);default:0.9" json:"top_p"`
	EnableTopP              bool      `gorm:"type:tinyint(1);default:0" json:"enable_top_p"`
	MaxTokens               int       `gorm:"type:int;default:4096" json:"max_tokens"`
	EnableMaxTokens         bool      `gorm:"type:tinyint(1);default:0" json:"enable_max_tokens"`
	MaxContextCount         int       `gorm:"type:int;default:10" json:"max_context_count"`
	RetrievalTopK           int       `gorm:"type:int;default:5" json:"retrieval_top_k"`
	RetrievalMatchThreshold float64   `gorm:"type:decimal(4,2);default:0.5" json:"retrieval_match_threshold"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	CreatedBy               string    `gorm:"type:varchar(100)" json:"created_by"`
	UpdatedBy               string    `gorm:"type:varchar(100)" json:"updated_by"`
}

func (ModelConfig) TableName() string { return "data_model_config" }

// DefaultModelConfig 返回未保存时使用的默认模型配置。
func DefaultModelConfig() *ModelConfig {
	return &ModelConfig{
		ID:                      ModelConfigDefaultID,
		Temperature:             DefaultModelTemperature,
		TopP:                    DefaultModelTopP,
		MaxTokens:               DefaultModelMaxTokens,
		MaxContextCount:         DefaultModelMaxContext,
		RetrievalTopK:           DefaultModelRetrievalTopK,
		RetrievalMatchThreshold: DefaultModelMatchThreshold,
	}
}
