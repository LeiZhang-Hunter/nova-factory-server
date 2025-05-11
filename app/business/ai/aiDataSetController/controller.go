package aiDataSetController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDataset, wire.Struct(new(AiDataSet), "*"))

type AiDataSet struct {
	Dataset *Dataset
}
