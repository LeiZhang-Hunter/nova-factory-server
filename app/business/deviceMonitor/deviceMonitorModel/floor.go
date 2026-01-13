package deviceMonitorModel

import (
	"nova-factory-server/app/business/asset/building/buildingModels"
	"nova-factory-server/app/business/asset/device/deviceModels"
)

// DeviceLayoutData 楼层布局
type DeviceLayoutData struct {
	// layout 楼层
	Layout    *buildingModels.SysFloor          `json:"layout"`
	DeviceMap map[string]*deviceModels.DeviceVO `json:"deviceMap"`
}
