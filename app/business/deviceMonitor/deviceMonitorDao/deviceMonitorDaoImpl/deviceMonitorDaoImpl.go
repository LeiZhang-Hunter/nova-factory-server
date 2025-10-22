package deviceMonitorDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewIDeviceDataReportDaoImpl,
	NewDeviceUtilizationDaoImpl,
)
