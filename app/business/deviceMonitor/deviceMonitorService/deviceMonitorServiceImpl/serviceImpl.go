package deviceMonitorServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDeviceMonitorServiceImpl, NewDeviceUtilizationServiceImpl, NewControlLogServiceImpl)
