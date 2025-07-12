package deviceDaoImpl

import (
	"github.com/google/wire"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao/deviceMonitorDaoImpl"
)

var ProviderSet = wire.NewSet(
	NewSysDeviceDaoImpl,
	NewSysDeviceGroupDaoImpl,
	NewIDeviceTemplateDaoImpl,
	NewISysModbusDeviceConfigDataDaoImp,
	deviceMonitorDaoImpl.NewIDeviceDataReportDaoImpl,
)
