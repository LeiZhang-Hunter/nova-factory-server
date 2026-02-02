package deviceMonitorModel

type ControlStatusItem struct {
	DeviceId uint64 `json:"device_id,string"`
	DataId   uint64 `json:"data_id,string"`
}

type ControlStatusReq struct {
	Items []ControlStatusItem `json:"items"`
}

type ControlStatusItemRes struct {
	DeviceId uint64 `json:"device_id"`
	DataId   uint64 `json:"data_id"`
	Status   int    `json:"status"` // 0: completed, 1: in progress
}

type ControlStatusRes struct {
	Items []ControlStatusItemRes `json:"items"`
}
