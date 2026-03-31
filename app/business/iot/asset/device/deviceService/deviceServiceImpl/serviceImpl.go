package deviceServiceImpl

import (
	"nova-factory-server/app/business/iot/deviceMonitor/deviceMonitorService/deviceMonitorServiceImpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewDeviceService,
	NewDeviceGroupService,
	NewDeviceTemplateServiceImpl,
	NewISysModbusDeviceConfigDataServiceImpl,
	deviceMonitorServiceImpl.NewIDeviceDataReportServiceImpl,
	NewIDeviceSubjectServiceImpl,
	NewIDeviceCheckMachineryServiceImpl,
	NewIDeviceCheckSubjectServiceImpl,
	NewIDeviceCheckPlanServiceImpl)
