package qcTemplateDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewQcTemplateDaoImpl)
