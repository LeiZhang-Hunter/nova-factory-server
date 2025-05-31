package metricController

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewMetric, wire.Struct(new(MetricServer), "*"))

type MetricServer struct {
	Metric *Metric
}
