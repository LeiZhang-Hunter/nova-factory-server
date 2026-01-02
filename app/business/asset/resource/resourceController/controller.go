package resourceController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewResourceFile, wire.Struct(new(ResourceController), "*"))

type ResourceController struct {
	ResourceFile *ResourceFile
}
