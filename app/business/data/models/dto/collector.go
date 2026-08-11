package dto

import "time"

type CreateCollectorRequest struct {
	Name     string   `json:"name" binding:"required,max=255"`
	DeviceID string   `json:"device_id" binding:"max=100"`
	Token    string   `json:"token"`
	Tags     []string `json:"tags"`
}

type UpdateCollectorRequest struct {
	Name     string   `json:"name" binding:"max=255"`
	DeviceID string   `json:"device_id" binding:"max=100"`
	Token    string   `json:"token"`
	Tags     []string `json:"tags"`
}

// CollectorResponse 采集器响应
type CollectorResponse struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	DeviceID         string     `json:"device_id"`
	Token            string     `json:"token"`
	Status           string     `json:"status"`
	LastHeartbeat    *time.Time `json:"last_heartbeat,omitempty"`
	LastConnectedAt  *time.Time `json:"last_connected_at,omitempty"`
	LastChecksum     string     `json:"last_checksum,omitempty"`
	LastDispatchedAt *time.Time `json:"last_dispatched_at,omitempty"`
	Version          string     `json:"version"`
	Tags             []string   `json:"tags"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
