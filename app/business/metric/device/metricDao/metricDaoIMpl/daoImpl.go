package metricDaoIMpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewMetricDaoImpl, NewIControlLogDaoImpl)
