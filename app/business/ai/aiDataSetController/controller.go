package aiDataSetController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDataset, NewException, NewPrediction, wire.Struct(new(AiDataSet), "*"))

type AiDataSet struct {
	Dataset    *Dataset
	Exception  *Exception
	Prediction *Prediction
}
