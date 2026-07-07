package iotdb

// MetricSeries represents a queried time series.
type MetricSeries interface{}

// MetricResult contains a set of series, it can also produce final result like aggregation value
type MetricResult interface {
	MetricMeta
	// AddSeries receives series and saves in MetricResult, which can be used for generating the final value
	AddSeries(MetricSeries) error
}
