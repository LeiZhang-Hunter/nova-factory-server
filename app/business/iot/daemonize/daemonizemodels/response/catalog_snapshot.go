package response

import "nova-factory-server/app/business/iot/asset/device/devicemodels"

// CatalogSnapshot is an internally consistent gateway configuration read.
type CatalogSnapshot struct {
	Devices   []*devicemodels.DeviceVO
	Templates []*devicemodels.SysDeviceTemplate
	Points    []*devicemodels.SysModbusDeviceConfigData
}
