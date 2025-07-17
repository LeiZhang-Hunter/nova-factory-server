package alertDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAlertSinkTemplateDaoImpl, NewAlertRuleDaoImpl)
