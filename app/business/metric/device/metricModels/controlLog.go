package metricModels

import "github.com/gogf/gf/os/gtime"

// NovaControlLog 控制日志
type NovaControlLog struct {
	DeviceId      uint64            `json:"device_id,string"       gorm:"column:device_id"      description:"设备id"` //
	DataId        uint64            `json:"data_id,string" gorm:"column:data_id"`
	Message       string            `json:"message"        orm:"message"        description:"日志消息"`                        // 公司uuid
	Type          string            `json:"type"        orm:"message"        description:"类型"`                             // 公司uuid
	SeriesId      uint64            `json:"series_id,string"          gorm:"column:series_id"          description:"序列id"` //
	Attributes    map[string]string `json:"attributes"      gorm:"column:attributes;type:JSONB"    description:"属性"`       //
	StartTimeUnix *gtime.Time       `json:"start_time_unix" gorm:"column:start_time_unix" description:"开始时间"`              //
	TimeUnix      *gtime.Time       `json:"time_unix"       gorm:"column:time_unix"      description:"当前时间"`               //
}
