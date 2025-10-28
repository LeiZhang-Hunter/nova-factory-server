package deviceController

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewDeviceInfo,
	NewDeviceGroup,
	NewTemplate,
	NewTemplateData,
	NewDeviceSubject,
	wire.Struct(new(Device), "*"))

type Device struct {
	Info          *DeviceInfo
	Group         *DeviceGroup
	Template      *Template
	TemplateData  *TemplateData
	DeviceSubject *DeviceSubject
}
