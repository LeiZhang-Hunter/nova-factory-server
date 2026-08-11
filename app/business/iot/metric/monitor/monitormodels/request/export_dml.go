// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package request

import "time"

// ExportDML is the payload sent by the monitor exporter.
type ExportDML struct {
	Signal    string    `json:"signal"`
	Dataset   string    `json:"dataset,omitempty"`
	Format    string    `json:"format"`
	CreatedAt time.Time `json:"created_at"`
	Body      []byte    `json:"body"`
}
