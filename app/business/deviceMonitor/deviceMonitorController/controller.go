package deviceMonitorController

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewDeviceMonitor, wire.Struct(new(DeviceMonitorController), "*"))

type DeviceMonitorController struct {
	DeviceMonitor *DeviceMonitor
}
