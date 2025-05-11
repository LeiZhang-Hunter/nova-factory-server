package materialController

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewMaterialInfo, wire.Struct(new(Material), "*"))

type Material struct {
	Material *MaterialInfo
}
