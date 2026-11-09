package daemonizecontroller

import (
	"context"
	"errors"
	"math"
	"strconv"

	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
)

// ValidateTelemetry authenticates the gateway and checks every cloud identity before writing.
func (c *Config) ValidateTelemetry(ctx context.Context, request *v1.ExportMetricsServiceRequest) error {
	id, err := c.authorizeGateway(ctx, "")
	if err != nil {
		return errors.New("gateway authentication failed")
	}
	if request == nil {
		return errors.New("empty telemetry request")
	}
	catalog, err := c.service.Catalog(ctx, int64(id))
	if err != nil {
		return errors.New("cannot validate gateway catalog")
	}
	allowed := make(map[[3]uint64]bool)
	for _, d := range catalog.Devices {
		deviceID, _ := strconv.ParseUint(d.ID, 10, 64)
		templateID, _ := strconv.ParseUint(d.TemplateID, 10, 64)
		for _, p := range d.Points {
			pointID, _ := strconv.ParseUint(p.ID, 10, 64)
			allowed[[3]uint64{deviceID, templateID, pointID}] = true
		}
	}
	count := 0
	for _, resource := range request.ResourceMetrics {
		if resource == nil {
			continue
		}
		for _, scope := range resource.ScopeMetrics {
			if scope == nil {
				continue
			}
			for _, sample := range scope.Metrics {
				count++
				if count > 1000 || sample == nil {
					return errors.New("invalid telemetry batch")
				}
				if !allowed[[3]uint64{sample.DeviceId, sample.TemplateId, sample.DataId}] || sample.DeviceId == 0 || sample.TemplateId == 0 ||
					sample.DataId == 0 {
					return errors.New("telemetry identity outside gateway catalog")
				}
				if sample.TimeUnixNano == 0 || sample.TimeUnixNano > math.MaxInt64 || sample.StartTimeUnixNano > math.MaxInt64 {
					return errors.New("invalid telemetry timestamp")
				}
				switch v := sample.Value.(type) {
				case *v1.DeviceMetric_AsInt:
				case *v1.DeviceMetric_AsDouble:
					if math.IsNaN(v.AsDouble) || math.IsInf(v.AsDouble, 0) {
						return errors.New("non-finite telemetry")
					}
				default:
					return errors.New("unsupported telemetry value")
				}
			}
		}
	}
	return nil
}
