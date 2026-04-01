package systemdaoimpl

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewIDeviceElectricDaoImpl)
