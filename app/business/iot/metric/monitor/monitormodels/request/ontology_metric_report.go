// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package request

import "time"

// OntologyMetricReport is posted by nova-ontology-collector metrics sink.
type OntologyMetricReport struct {
	Sink           string                  `json:"sink"`
	Pipeline       string                  `json:"pipeline,omitempty"`
	EntityID       string                  `json:"entityId"`
	EntityType     string                  `json:"entityType"`
	SourceRecordID string                  `json:"sourceRecordId,omitempty"`
	Source         string                  `json:"source,omitempty"`
	FileName       string                  `json:"fileName,omitempty"`
	Path           string                  `json:"path,omitempty"`
	Timestamp      time.Time               `json:"timestamp"`
	Metrics        []OntologyMetricValue   `json:"metrics"`
	Summary        []OntologyMetricSummary `json:"summary,omitempty"`
	Items          []OntologyMetricItem    `json:"items,omitempty"`
	Attributes     map[string]any          `json:"attributes,omitempty"`
}

type OntologyMetricItem struct {
	SheetName   string                     `json:"sheetName"`
	RowIndex    string                     `json:"rowIndex"`
	Module      string                     `json:"module"`
	ModuleDesc  string                     `json:"moduleDesc,omitempty"`
	BlockLayout string                     `json:"blockLayout,omitempty"`
	Fields      []OntologyMetricFieldValue `json:"fields"`
}

type OntologyMetricFieldValue struct {
	TargetIndex   string `json:"targetIndex"`
	ExcelColName  string `json:"excelColName,omitempty"`
	ExcelColIndex *int   `json:"excelColIndex,omitempty"`
	DataType      string `json:"dataType,omitempty"`
	Desc          string `json:"desc,omitempty"`
	Value         any    `json:"value"`
}

type OntologyMetricValue struct {
	IndexName      string `json:"indexName"`
	IndexDesc      string `json:"indexDesc,omitempty"`
	DataType       string `json:"dataType,omitempty"`
	Value          any    `json:"value"`
	SheetName      string `json:"sheetName,omitempty"`
	ComputeFormula string `json:"computeFormula,omitempty"`
}

type OntologyMetricSummary struct {
	IndexName string `json:"indexName"`
	Value     any    `json:"value"`
}
