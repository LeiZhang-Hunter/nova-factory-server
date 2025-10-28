package deviceServiceImpl

import (
	"github.com/google/wire"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService/deviceMonitorServiceImpl"
)

var ProviderSet = wire.NewSet(
	NewDeviceService,
	NewDeviceGroupService,
	NewDeviceTemplateServiceImpl,
	NewISysModbusDeviceConfigDataServiceImpl,
	deviceMonitorServiceImpl.NewIDeviceDataReportServiceImpl,
	NewIDeviceSubjectServiceImpl)
