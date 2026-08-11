package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	StatusOnline  = "online"
	StatusOffline = "offline"
)

// Collector 采集器设备
type Collector struct {
	ID               string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name             string         `gorm:"type:varchar(255);not null" json:"name"`
	DeviceID         string         `gorm:"type:varchar(100);index" json:"device_id"`
	Token            string         `gorm:"type:varchar(255)" json:"token"`
	Status           string         `gorm:"type:varchar(20);default:offline;index" json:"status"`
	LastHeartbeat    *time.Time     `json:"last_heartbeat,omitempty"`
	LastConnectedAt  *time.Time     `json:"last_connected_at,omitempty"`
	LastChecksum     string         `gorm:"type:varchar(64)" json:"last_checksum,omitempty"`
	LastDispatchedAt *time.Time     `json:"last_dispatched_at,omitempty"`
	Version          string         `gorm:"type:varchar(50)" json:"version"`
	Tags             string         `gorm:"type:json" json:"tags,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy        string         `gorm:"type:varchar(100)" json:"created_by"`
	UpdatedBy        string         `gorm:"type:varchar(100)" json:"updated_by"`
}

func (Collector) TableName() string { return "data_collectors" }
