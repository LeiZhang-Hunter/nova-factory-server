package deviceMonitorModel

import "nova-factory-server/app/baize"

type ControlLogListReq struct {
	DeviceID int64  `form:"device_id,string"` // 设备id
	DataId   int64  `form:"data_id,string"`   // 设备id
	Start    uint64 `form:"start"`
	End      uint64 `form:"end"`
	baize.BaseEntityDQL
}
