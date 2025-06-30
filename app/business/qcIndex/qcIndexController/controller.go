package qcIndexController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewQcIndexController, wire.Struct(new(QcIndexRoute), "*"))

type QcIndexRoute struct {
	QcIndexController *QcIndexController
}
