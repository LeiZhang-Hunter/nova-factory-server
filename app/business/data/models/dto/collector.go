package dto

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
