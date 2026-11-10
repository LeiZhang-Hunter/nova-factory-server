package iotdb

import (
	"testing"

	metricmodels "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
)

// TestDeviceTemplateDefinitionsUseDataIDMeasurements 验证模板使用测点 ID 作为列名。
func TestDeviceTemplateDefinitionsUseDataIDMeasurements(t *testing.T) {
	got := deviceTemplateDefinitions([]int64{8052648150306816, 8053303506112512})
	want := "d8052648150306816 DOUBLE, d8053303506112512 DOUBLE"
	if got != want {
		t.Fatalf("deviceTemplateDefinitions() = %q, want %q", got, want)
	}
}

// TestDeviceMetricQueryUsesTemplateMeasurements 验证查询按模板的测点列构造 SQL。
func TestDeviceMetricQueryUsesTemplateMeasurements(t *testing.T) {
	req := &metricmodels.MetricDataQueryReq{
		Name:       "ignored",
		Expression: "avg(value)",
		Field:      " as value",
		QueryMetric: []*metricmodels.MetricQueryCondition{
			{DeviceId: 1210552827252736, DataId: 8052648150306816},
			{DeviceId: 1210552827252736, DataId: 8053303506112512},
		},
	}

	name, expression := deviceMetricQuery(req)
	wantName := "root.device.dev1210552827252736"
	wantExpression := "avg(d8052648150306816) as d8052648150306816, avg(d8053303506112512) as d8053303506112512"
	if name != wantName || expression != wantExpression {
		t.Fatalf("deviceMetricQuery() = (%q, %q), want (%q, %q)", name, expression, wantName, wantExpression)
	}
}
