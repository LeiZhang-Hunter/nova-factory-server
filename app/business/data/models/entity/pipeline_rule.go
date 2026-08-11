package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	SourceTypeFile  = "file"
	SourceTypeMySQL = "mysql"
	SourceTypeAPI   = "api"

	RuleStatusEnabled  = "enabled"
	RuleStatusDisabled = "disabled"
)

type PipelineRule struct {
	ID          string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null;index" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	SourceType  string         `gorm:"type:varchar(20);not null;index" json:"source_type"`
	Status      string         `gorm:"type:varchar(20);not null;default:enabled;index" json:"status"`
	Version     int            `gorm:"type:int;not null;default:1" json:"version"`
	Config      string         `gorm:"type:json;not null" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy   string         `gorm:"type:varchar(100)" json:"created_by"`
	UpdatedBy   string         `gorm:"type:varchar(100)" json:"updated_by"`
}

func (PipelineRule) TableName() string { return "data_pipeline_rules" }

type ServiceConnection struct {
	ID                   string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name                 string         `gorm:"type:varchar(255);not null;index" json:"name"`
	Description          string         `gorm:"type:text" json:"description"`
	SourceType           string         `gorm:"type:varchar(20);not null;index" json:"source_type"`
	Status               string         `gorm:"type:varchar(20);not null;default:enabled;index" json:"status"`
	Config               string         `gorm:"type:json;not null" json:"-"`
	EncryptedCredentials string         `gorm:"type:text" json:"-"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy            string         `gorm:"type:varchar(100)" json:"created_by"`
	UpdatedBy            string         `gorm:"type:varchar(100)" json:"updated_by"`
}

func (ServiceConnection) TableName() string { return "data_service_connections" }
