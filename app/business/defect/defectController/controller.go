package defectController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDefectController, wire.Struct(new(DefectRoute), "*"))

type DefectRoute struct {
	DefectController *DefectController
}
