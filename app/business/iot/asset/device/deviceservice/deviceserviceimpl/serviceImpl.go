package deviceserviceimpl

import (
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitorservice/deviceMonitorServiceImpl"

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
