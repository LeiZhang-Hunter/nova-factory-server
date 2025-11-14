package deviceMonitorController

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewDeviceMonitor, NewDeviceReport, NewDeviceUtilization, NewControlLog,
	wire.Struct(new(DeviceMonitorController), "*"))

type DeviceMonitorController struct {
	DeviceMonitor     *DeviceMonitor
	DeviceReport      *DeviceReport
	DeviceUtilization *DeviceUtilization
	ControlLog        *ControlLog
}
