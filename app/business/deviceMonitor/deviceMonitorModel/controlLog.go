package deviceMonitorModel

import "nova-factory-server/app/baize"

type ControlLogListReq struct {
	DeviceID int64 `form:"device_id,string"` // 设备id
	DataId   int64 `form:"data_id,string"`   // 设备id
	baize.BaseEntityDQL
}
