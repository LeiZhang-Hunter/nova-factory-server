// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package entity

// ExportPayload is the internal representation created from OpenTelemetry data.
type ExportPayload struct {
	Signal string
	Format string
	Body   []byte
}
