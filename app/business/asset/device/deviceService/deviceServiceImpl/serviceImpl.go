package deviceServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewDeviceService,
	NewDeviceGroupService,
	NewDeviceTemplateServiceImpl,
	NewISysModbusDeviceConfigDataServiceImpl)
