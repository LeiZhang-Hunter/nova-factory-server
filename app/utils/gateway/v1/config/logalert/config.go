package logalert

import "time"

type Config struct {
	Addr         []string      `yaml:"addr,omitempty"`
	BufferSize   int           `yaml:"bufferSize,omitempty" default:"100"`
	BatchTimeout time.Duration `yaml:"batchTimeout,omitempty" default:"10s"`
	BatchSize    int           `yaml:"batchSize,omitempty" default:"10"`
	AlertConfig  `yaml:",inline"`
}

type AlertConfig struct {
	Template               string            `yaml:"template,omitempty"`
	Timeout                time.Duration     `yaml:"timeout,omitempty" default:"30s"`
	Headers                map[string]string `yaml:"headers,omitempty"`
	Method                 string            `yaml:"method,omitempty"`
	LineLimit              int               `yaml:"lineLimit,omitempty" default:"10"`
	GroupKey               string            `yaml:"groupKey,omitempty" default:"${_meta.pipelineName}-${_meta.sourceName}"`
	AlertSendingThreshold  int               `yaml:"alertSendingThreshold,omitempty" default:"1"`
	SendLogAlertAtOnce     bool              `yaml:"sendLogAlertAtOnce"`
	SendNoDataAlertAtOnce  bool              `yaml:"sendNoDataAlertAtOnce" default:"true"`
	SendGatewayError       bool              `yaml:"sendGatewayError"`
	SendGatewayErrorAtOnce bool              `yaml:"sendGatewayErrorAtOnce"`
}
