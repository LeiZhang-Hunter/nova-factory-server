package deviceDaoImpl

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewSysDeviceDaoImpl,
	NewSysDeviceGroupDaoImpl,
	NewIDeviceTemplateDaoImpl,
	NewISysModbusDeviceConfigDataDaoImp,
	NewIDeviceSubjectDaoImpl,
	NewIDeviceCheckPlanDaoImpl,
	NewIDeviceCheckMachineryDaoImpl,
	NewIDeviceCheckSubjectDaoImpl,
)
