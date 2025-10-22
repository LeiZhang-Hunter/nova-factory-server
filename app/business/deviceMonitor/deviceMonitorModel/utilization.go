package deviceMonitorModel

import "nova-factory-server/app/constant/device"

// DeviceUtilizationReq 稼动率请求
type DeviceUtilizationReq struct {
	Start uint64 `json:"start"`
	End   uint64 `json:"end"`
}

type DeviceStatus struct {
	DeviceId int64             `json:"deviceId"`
	Value    float64           `json:"value"`
	Status   device.RUN_STATUS `json:"status"`
}

type DeviceStatusList struct {
	List []DeviceStatus
}

func NewDeviceStatusList() *DeviceStatusList {
	return &DeviceStatusList{
		List: make([]DeviceStatus, 0),
	}
}
