package metricController

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewMetric, NewCamera, wire.Struct(new(MetricServer), "*"))

type MetricServer struct {
	Metric     *Metric
	CameraGrpc *CameraGrpc
}
