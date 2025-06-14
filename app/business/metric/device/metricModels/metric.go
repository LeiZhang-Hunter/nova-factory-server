package metricModels

import "github.com/gogf/gf/os/gtime"

// NovaMetricsDevice is the golang structure for table nova_metrics_device.
type NovaMetricsDevice struct {
	// `gorm:"column:dataset_id;primaryKey;comment:出库id" json:"datasetId,string"`                                                // 出库id
	DataId        uint64            `json:"dataId"`
	DeviceId      uint64            `json:"device_id"       gorm:"column:device_id"      description:"设备id"`            //
	TemplateId    uint64            `json:"template_id"     gorm:"column:template_id"    description:"设备模板id"`          //
	SeriesId      uint64            `json:"__series_id__"          gorm:"column:series_id"          description:"序列id"` //
	Attributes    map[string]string `json:"attributes"      gorm:"column:attributes"    description:"属性"`               //
	StartTimeUnix *gtime.Time       `json:"start_time_unix" gorm:"column:start_time_unix" description:"开始时间"`           //
	TimeUnix      *gtime.Time       `json:"time_unix"       gorm:"column:time_unix"      description:"当前时间"`            //
	Value         float64           `json:"value"           gorm:"column:value"         description:"统计值"`              //
}

type DeviceMetricData struct {
	SeriesId      uint64            `json:"__series_id__"          gorm:"column:series_id"          description:"序列id"` //
	Attributes    map[string]string `json:"attributes"      gorm:"column:attributes"    description:"属性"`               //
	StartTimeUnix *gtime.Time       `json:"start_time_unix" gorm:"column:start_time_unix" description:"开始时间"`           //
	TimeUnix      *gtime.Time       `json:"time_unix"       gorm:"column:time_unix"      description:"当前时间"`            //
	Value         float64           `json:"value"           gorm:"column:value"         description:"统计值"`              //
}

type MetricMap struct {
	Data map[uint64]map[uint64]map[uint64]*DeviceMetricData // device_id =>? template_id => data_id
}

func NewMetricMap() *MetricMap {
	return &MetricMap{
		Data: make(map[uint64]map[uint64]map[uint64]*DeviceMetricData),
	}
}
