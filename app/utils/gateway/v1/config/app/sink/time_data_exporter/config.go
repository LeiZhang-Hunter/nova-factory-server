package time_data_exporter

import "time"

// Config sink配置
type Config struct {
	Address string `yaml:"address,omitempty" default:"127.0.0.1:6000"`
	// in general, we use codec.PrintEvents instead.
	PrintEvents bool `yaml:"printEvents,omitempty" default:"false"`
	// Within this period, only one log event is printed for troubleshooting.
	PrintEventsInterval time.Duration `yaml:"printEventsInterval,omitempty"`

	PrintMetrics    bool          `yaml:"printMetrics,omitempty"`
	MetricsInterval time.Duration `yaml:"printMetricsInterval,omitempty" default:"1s"`

	// resultStatus can be used to simulate failure, drop
	ResultStatus string `yaml:"resultStatus,omitempty" default:"success"`
}
