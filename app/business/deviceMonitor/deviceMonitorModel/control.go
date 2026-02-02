package deviceMonitorModel

// ControlReq 控制请求
type ControlReq struct {
	DeviceId uint64      `json:"device_id,string"`
	DataId   uint64      `json:"data_id,string"`
	AgentId  uint64      `json:"agent_id,string"`
	Value    interface{} `json:"value"`
}

// ControlRes 控制结果
type ControlRes struct {
	Code int    `json:"code"` // 0: success, 201: in progress, 404: agent not found, 504: timeout
	Msg  string `json:"msg"`
}
