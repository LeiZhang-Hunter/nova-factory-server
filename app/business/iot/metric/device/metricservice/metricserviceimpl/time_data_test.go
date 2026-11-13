package metricserviceimpl

import (
	"context"
	"math"
	"testing"

	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"

	"nova-factory-server/app/business/iot/metric/device/metricdao"
)

type timeDataCaptureDAO struct {
	metricdao.IMetricDao
	data map[string][]*v1.ResourceTimeMetrics
}

// ExportTimeData 捕获服务层提交到 DAO 的数据。
func (d *timeDataCaptureDAO) ExportTimeData(_ context.Context, data map[string][]*v1.ResourceTimeMetrics) error {
	d.data = data
	return nil
}

// TestExportTimeDataIntegerFields 验证运行模板字段转换、整数精度及输入隔离。
func TestExportTimeDataIntegerFields(t *testing.T) {
	const table = "root.run_status_device.dev1898806979334144"
	const largeInteger int64 = 9007199254740993
	resource := &v1.ResourceTimeMetrics{Table: table, TimeUnixNano: 1794590079000000000, Metrics: []*v1.TimeDataMetric{
		{Field: "duration", Value: &v1.TimeDataMetric_AsDouble{AsDouble: 60}},
		{Field: "status", Value: &v1.TimeDataMetric_AsInt{AsInt: largeInteger}},
		{Field: "other", Value: &v1.TimeDataMetric_AsDouble{AsDouble: 1.5}},
	}}
	dao := &timeDataCaptureDAO{}
	service := &IMetricServiceImpl{dao: dao}
	if err := service.ExportTimeData(context.Background(), &v1.ExportTimeDataRequest{ResourceMetrics: []*v1.ResourceTimeMetrics{nil, resource}}); err != nil {
		t.Fatal(err)
	}
	got := dao.data[table][0]
	if _, ok := got.Metrics[0].Value.(*v1.TimeDataMetric_AsInt); !ok || got.Metrics[0].GetAsInt() != 60 {
		t.Fatalf("duration was not converted to integer: %v", got.Metrics[0])
	}
	if got.Metrics[1].GetAsInt() != largeInteger || got.Metrics[2].GetAsDouble() != 1.5 {
		t.Fatalf("other values changed: %v", got.Metrics)
	}
	if _, ok := resource.Metrics[0].Value.(*v1.TimeDataMetric_AsDouble); !ok {
		t.Fatal("original request was modified")
	}
}

// TestExportTimeDataRejectsLossyInteger 验证非法整数在调用 DAO 前被拒绝。
func TestExportTimeDataRejectsLossyInteger(t *testing.T) {
	for _, value := range []float64{60.5, math.NaN(), math.Inf(1), 0x1p63, -0x1p63 - 2048} {
		dao := &timeDataCaptureDAO{}
		service := &IMetricServiceImpl{dao: dao}
		request := &v1.ExportTimeDataRequest{ResourceMetrics: []*v1.ResourceTimeMetrics{{
			Table:   "root.run_status_device.dev1",
			Metrics: []*v1.TimeDataMetric{{Field: "duration", Value: &v1.TimeDataMetric_AsDouble{AsDouble: value}}},
		}}}
		if err := service.ExportTimeData(context.Background(), request); err == nil {
			t.Fatalf("expected invalid integer %v to be rejected", value)
		}
		if dao.data != nil {
			t.Fatal("DAO called with invalid data")
		}
	}
}
