package alertServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAlertTemplateServiceImpl, NewAlertRuleServiceImpl, NewAlertLogServiceImpl)
