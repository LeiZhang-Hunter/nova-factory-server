package metricserviceimpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/gtime"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	metricmodels "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	"nova-factory-server/app/business/iot/metric/device/metricservice"
	"nova-factory-server/app/constant/device"
	iotdbconstant "nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/cache"
)

type IMetricServiceImpl struct {
	dao   metricdao.IMetricDao
	cache cache.Cache
}

func NewIMetricServiceImpl(dao metricdao.IMetricDao, cache cache.Cache) metricservice.IMetricService {
	return &IMetricServiceImpl{
		dao:   dao,
		cache: cache,
	}
}

func (m *IMetricServiceImpl) Export(c context.Context, request *v1.ExportMetricsServiceRequest) error {
	if request == nil {
		return nil
	}
	if len(request.ResourceMetrics) == 0 {
		return nil
	}

	for _, resourceMetric := range request.ResourceMetrics {
		if resourceMetric == nil {
			continue
		}

		if len(resourceMetric.ScopeMetrics) == 0 {
			continue
		}

		for _, metric := range resourceMetric.ScopeMetrics {
			if metric == nil {
				continue
			}

			if len(metric.Metrics) == 0 {
				continue
			}

			var values []*metricmodels.NovaMetricsDevice = make([]*metricmodels.NovaMetricsDevice, 0)
			for _, mVale := range metric.Metrics {
				if mVale == nil {
					return errors.New("nil device metric")
				}
				var v float64
				if mVale.GetValue() == nil {
					v = 0
				} else {
					ptr, ok := mVale.GetValue().(*v1.DeviceMetric_AsDouble)
					if ok {
						v = ptr.AsDouble
					} else {
						v = float64(mVale.GetAsInt())
					}
				}
				values = append(values, &metricmodels.NovaMetricsDevice{
					DeviceId:   mVale.DeviceId,
					TemplateId: mVale.TemplateId,
					Attributes: map[string]string{
						"test": "1111",
					},
					DataId:        mVale.DataId,
					SeriesId:      request.SeriesId,
					StartTimeUnix: gtime.New(time.UnixMicro(int64(mVale.StartTimeUnixNano))),
					TimeUnix:      gtime.New(time.UnixMicro(int64(mVale.TimeUnixNano))),
					Value:         v,
				})
			}
			err := m.export(c, values, request.SeriesId)
			if err != nil {
				zap.L().Error("export error", zap.Error(err))
				return err
			}
		}

	}
	return nil
}

func (m *IMetricServiceImpl) export(c context.Context, values []*metricmodels.NovaMetricsDevice, seriesId uint64) error {
	data := metricmodels.NewMetricMap()

	// v.Value
	// 格式化数据结构
	for _, v := range values {
		_, ok := data.Data[uint64(v.DeviceId)]
		if !ok {
			data.Data[uint64(v.DeviceId)] = make(map[uint64]map[uint64]*metricmodels.DeviceMetricData)
		}

		_, ok = data.Data[uint64(v.DeviceId)][uint64(v.TemplateId)]
		if !ok {
			data.Data[uint64(v.DeviceId)][v.TemplateId] = make(map[uint64]*metricmodels.DeviceMetricData)
		}

		mapValue := make(map[string]string)
		data.Data[uint64(v.DeviceId)][uint64(v.TemplateId)][uint64(v.DataId)] = &metricmodels.DeviceMetricData{
			SeriesId:      seriesId,
			Attributes:    mapValue,
			StartTimeUnix: v.StartTimeUnix,
			TimeUnix:      v.TimeUnix,
			Value:         v.Value,
		}

		bytes, err := json.Marshal(data.Data[uint64(v.DeviceId)])
		if err != nil {
			zap.L().Error("json marshal error", zap.Error(err))
			continue
		}
		m.cache.Set(context.Background(), device.MakeDeviceKey(uint64(v.DeviceId)), string(bytes), 600*time.Second)
	}

	return m.dao.Export(c, values)
}

func (m *IMetricServiceImpl) List(c *gin.Context, req *devicemonitormodel.DevDataReq) (*devicemonitormodel.DevDataResp, error) {
	list, err := m.dao.List(c, req)
	if err != nil {
		zap.L().Error("get dev data list error", zap.Error(err))
		return nil, err
	}
	if list == nil {
		return &devicemonitormodel.DevDataResp{
			Rows: make([]devicemonitormodel.DevData, 0),
		}, nil
	}
	total, err := m.dao.Count(c, req)
	if err != nil {
		return nil, err
	}
	list.Total = total
	var devStrMap map[string]bool = make(map[string]bool)
	for _, v := range list.Rows {
		devStrMap[v.Dev] = true
	}
	var devStrList []string = make([]string, len(devStrMap))
	count := 0
	for k, _ := range devStrMap {
		devStrList[count] = k
		count++
	}

	return list, nil
}

// ExportTimeData 将运行状态字段规范为整数后按表分组写入。
func (m *IMetricServiceImpl) ExportTimeData(c context.Context, request *v1.ExportTimeDataRequest) error {
	if request == nil {
		return nil
	}

	if len(request.ResourceMetrics) == 0 {
		return nil
	}

	var list map[string][]*v1.ResourceTimeMetrics = make(map[string][]*v1.ResourceTimeMetrics)
	for _, v := range request.ResourceMetrics {
		if v == nil {
			continue
		}
		if strings.HasPrefix(v.Table, iotdbconstant.ROOT_RUN_STATUS_DEVICE_TEMPLATE_NAME) {
			v = proto.Clone(v).(*v1.ResourceTimeMetrics)
			if err := normalizeRunStatusMetrics(v); err != nil {
				return err
			}
		}
		_, ok := list[v.Table]
		if !ok {
			list[v.Table] = make([]*v1.ResourceTimeMetrics, 0)
		}
		list[v.Table] = append(list[v.Table], v)
	}

	if len(list) == 0 {
		return nil
	}
	err := m.dao.ExportTimeData(c, list)
	if err != nil {
		zap.L().Error("export time data error", zap.Error(err))
		return err
	}
	return nil
}

// normalizeRunStatusMetrics 将模板要求的整数值转换为 AsInt，拒绝有损转换。
func normalizeRunStatusMetrics(resource *v1.ResourceTimeMetrics) error {
	for _, metric := range resource.Metrics {
		if metric == nil || (metric.Field != "duration" && metric.Field != "status") {
			continue
		}
		switch value := metric.GetValue().(type) {
		case *v1.TimeDataMetric_AsInt:
			continue
		case *v1.TimeDataMetric_AsDouble:
			if math.IsNaN(value.AsDouble) || math.IsInf(value.AsDouble, 0) ||
				value.AsDouble < -0x1p63 || value.AsDouble >= 0x1p63 || math.Trunc(value.AsDouble) != value.AsDouble {
				return fmt.Errorf("%s.%s requires an INT64 value, got %v", resource.Table, metric.Field, value.AsDouble)
			}
			metric.Value = &v1.TimeDataMetric_AsInt{AsInt: int64(value.AsDouble)}
		case nil:
			metric.Value = &v1.TimeDataMetric_AsInt{AsInt: 0}
		}
	}
	return nil
}

func (m *IMetricServiceImpl) ExportScheduleLog(c context.Context, request *v1.ExportTimeDataRequest) error {
	return nil
}
