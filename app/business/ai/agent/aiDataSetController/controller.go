package aiDataSetController

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewDataset, NewException, NewPrediction, NewControl, NewModel,
	NewOCR, wire.Struct(new(AiDataSet), "*"))

type AiDataSet struct {
	Dataset    *Dataset
	Exception  *Exception
	Prediction *Prediction
	Control    *Control
	Model      *Model
	OCR        *OCR
}
