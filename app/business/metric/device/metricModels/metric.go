package metricModels

import (
	"github.com/gogf/gf/os/gtime"
)

// NovaMetricsDevice is the golang structure for table nova_metrics_device.
type NovaMetricsDevice struct {
	DeviceId      uint64            `json:"device_id,string"       gorm:"column:device_id"      description:"设备id"`   //
	TemplateId    uint64            `json:"template_id,string"     gorm:"column:template_id"    description:"设备模板id"` //
	DataId        uint64            `json:"data_id,string" gorm:"column:data_id"`
	SeriesId      uint64            `json:"series_id,string"          gorm:"column:series_id"          description:"序列id"` //
	Attributes    map[string]string `json:"attributes"      gorm:"column:attributes;type:JSONB"    description:"属性"`       //
	StartTimeUnix *gtime.Time       `json:"start_time_unix" gorm:"column:start_time_unix" description:"开始时间"`              //
	TimeUnix      *gtime.Time       `json:"time_unix"       gorm:"column:time_unix"      description:"当前时间"`               //
	Value         float64           `json:"value"           gorm:"column:value"         description:"统计值"`                 //
}

type DeviceMetricData struct {
	SeriesId      uint64            `json:"series_id,string"          gorm:"column:series_id"          description:"序列id"` //
	Attributes    map[string]string `json:"attributes"      gorm:"column:attributes"    description:"属性"`                  //
	StartTimeUnix *gtime.Time       `json:"start_time_unix" gorm:"column:start_time_unix" description:"开始时间"`              //
	Name          string            `json:"name" gorm:"-"`
	Unit          string            `json:"unit" gorm:"-"`
	TimeUnix      *gtime.Time       `json:"time_unix"       gorm:"column:time_unix"      description:"当前时间"` //
	Value         float64           `json:"value"           gorm:"column:value"         description:"统计值"`   //
}

type MetricMap struct {
	Data map[uint64]map[uint64]map[uint64]*DeviceMetricData // device_id =>? template_id => data_id
}

func NewMetricMap() *MetricMap {
	return &MetricMap{
		Data: make(map[uint64]map[uint64]map[uint64]*DeviceMetricData),
	}
}

// 指标

type MetricQueryValue struct {
	Time  int64  `json:"time"`
	Value string `json:"value"`
}

type MetricQueryData struct {
	Labels map[string]string `json:"label"`
	Values []MetricQueryValue
	Id     string
}

func NewMetricQueryData() *MetricQueryData {
	return &MetricQueryData{
		Labels: make(map[string]string),
		Values: make([]MetricQueryValue, 0),
	}
}

type MetricQueryReq struct {
	DeviceId    uint64 `json:"device_id,string"       gorm:"column:device_id"      description:"设备id"`   //
	TemplateId  uint64 `json:"template_id,string"     gorm:"column:template_id"    description:"设备模板id"` //
	DataId      uint64 `json:"data_id,string"`
	Query       string `json:"query"`
	Start       uint64 `json:"start"`
	End         uint64 `json:"end"`
	Step        int    `json:"step"`
	ServiceName string `json:"service_name"`
	Expression  string `json:"expression"`
	Limit       int    `json:"limit"`
}
