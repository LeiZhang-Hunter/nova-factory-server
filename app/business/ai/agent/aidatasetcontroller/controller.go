package aidatasetcontroller

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewDataset, NewException, NewPrediction, NewControl, NewModel, NewAgent,
	NewOCR, wire.Struct(new(AiDataSet), "*"))

type AiDataSet struct {
	Dataset    *Dataset
	Exception  *Exception
	Prediction *Prediction
	Control    *Control
	Model      *Model
	Agent      *Agent
	OCR        *OCR
}
