package metricserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewIMetricServiceImpl, NewIDevMapServiceImpl, NewIControlLogServiceImpl, NewICameraServiceImpl)
