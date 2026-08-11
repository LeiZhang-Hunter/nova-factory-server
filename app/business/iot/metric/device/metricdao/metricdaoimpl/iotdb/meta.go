package iotdb

import (
	"nova-factory-server/app/utils/uuid"
	"sort"
	"strings"
)

type iotMetricMeta struct {
	kind       string
	properties map[string]string
	name       string
}

func (m iotMetricMeta) GetKind() string {
	return m.kind
}

func (m iotMetricMeta) GetProperties() map[string]string {
	return m.properties
}

func (m iotMetricMeta) GetName() string {
	if len(m.properties) == 0 {
		return m.kind
	}

	keys := make([]string, 0, len(m.properties))
	for key := range m.properties {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	builder.WriteString(m.kind)
	for _, key := range keys {
		builder.WriteString("|")
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(m.properties[key])
	}
	metricBuildStr := builder.String()
	str := uuid.MakeMd5([]byte(metricBuildStr))
	return str
}
