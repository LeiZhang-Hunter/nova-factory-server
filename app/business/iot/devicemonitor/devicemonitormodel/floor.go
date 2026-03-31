package devicemonitormodel

import (
	"nova-factory-server/app/business/iot/asset/building/buildingmodels"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
)

// DeviceLayoutData 楼层布局
type DeviceLayoutData struct {
	// layout 楼层
	Layout    *buildingmodels.SysFloor          `json:"layout"`
	DeviceMap map[string]*devicemodels.DeviceVO `json:"deviceMap"`
}
