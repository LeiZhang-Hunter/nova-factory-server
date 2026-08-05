// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitorexporter

import (
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configoptional"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for the monitor exporter.
type Config struct {
	confighttp.ClientConfig `mapstructure:",squash"`

	QueueSettings configoptional.Optional[exporterhelper.QueueBatchConfig] `mapstructure:"sending_queue,omitempty"`
	BackOffConfig configretry.BackOffConfig                                `mapstructure:"retry_on_failure,omitempty"`

	Dataset string `mapstructure:"dataset,omitempty"`
}
