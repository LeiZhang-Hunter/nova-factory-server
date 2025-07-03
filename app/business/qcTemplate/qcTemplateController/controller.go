package qcTemplateController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewQcTemplateController, wire.Struct(new(QcTemplateRoute), "*"))

type QcTemplateRoute struct {
	QcTemplateController *QcTemplateController
}
