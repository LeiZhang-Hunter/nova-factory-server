package iotdb

// MetricMeta is the meta info of metric
type MetricMeta interface {
	// GetKind should returns the metric kind like pod_cpu_usage, pod_cpu_throttled
	GetKind() string
	// GetProperties should return the property of metric like pod_uid, container_id, gpu_device_name
	GetProperties() map[string]string
}

// MetricSample is a sample of specified metric, e.g. '{__name__: node_cpu_usage} = <2023-04-18:20:00:00, 4.1 core>'
type MetricSample interface {
	MetricMeta

	// timestamp returns the ts of metric
	timestamp() int64

	// value returns the metric value
	value() float64
}

// SelectHints specifies hints passed for query.
// It is only an option for implementation of MetricResult to use, e.g. GroupedResult
type QueryHints struct {
	// GroupBy []string
}
