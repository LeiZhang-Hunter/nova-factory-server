package alertModels

import (
	"encoding/json"
	"nova-factory-server/app/utils/snowflake"
	"strconv"

	"github.com/gogf/gf/os/gtime"
	"go.uber.org/zap"
)

// NovaAlertLog 控制日志
type NovaAlertLog struct {
	ObjectId         uint64      `json:"object_id,string" gorm:"column:object_id"`
	RuleName         string      `json:"rule_name" gorm:"-"`
	DeviceName       string      `json:"device_name" gorm:"-"`
	GatewayID        int64       `gorm:"column:gateway_id;not null;comment:网关id" json:"gateway_id"`                     // 网关id
	DeviceId         uint64      `json:"device_id,string"       gorm:"column:device_id"      description:"设备id"`        //
	DeviceTemplateID uint64      `gorm:"column:device_template_id;not null;comment:设备模板id" json:"device_template_id"`   // 设备模板id
	DeviceDataID     uint64      `gorm:"column:device_data_id;not null;comment:设备数据id" json:"device_data_id"`           // 设备数据id
	AlertID          uint64      `gorm:"column:alert_id;not null;comment:告警id" json:"alert_id"`                         // 告警id
	SeriesId         uint64      `json:"series_id,string"          gorm:"column:series_id"          description:"序列id"` //
	Context          string      `json:"context"        orm:"context"        description:"日志内容"`                        // 公司uuid
	Reason           string      `json:"reason"        orm:"reason"        description:"日志理由"`                          // 公司uuid
	Message          string      `json:"message"        orm:"message"        description:"日志消息"`                        // 公司uuid
	Data             string      `json:"data"      gorm:"column:data"    description:"属性"`                              //
	StartTimeUnix    *gtime.Time `json:"start_time_unix" gorm:"column:start_time_unix" description:"开始时间"`              //
	TimeUnix         *gtime.Time `json:"time_unix"       gorm:"column:time_unix"      description:"当前时间"`               //
}

func FromDataToNovaAlertLog(data *AlertLogData) ([]*NovaAlertLog, []*AlertLogInfo) {
	ret := make([]*NovaAlertLog, 0)
	infos := make([]*AlertLogInfo, 0)
	for _, alert := range data.Alerts {
		info := ToAlertLogInfo(data, &alert)
		alertJson, err := json.Marshal(info)
		if err != nil {
			zap.L().Error("json marshal fail", zap.Error(err))
			continue
		}

		var alertLog NovaAlertLog
		alertLog.ObjectId = uint64(snowflake.GenID())
		info.Id = uint64(alertLog.ObjectId)
		alertLog.Data = string(alertJson)
		alertLog.GatewayID = data.GatewayId
		alertLog.AlertID = uint64(data.Id)
		if alert.Labels.DeviceId != "" {
			alertLog.DeviceId, err = strconv.ParseUint(alert.Labels.DeviceId, 10, 64)
			if err != nil {
				zap.L().Error("parse int error", zap.Error(err))
			}
		} else {
			alertLog.DeviceId = 0
		}

		alertLog.DeviceTemplateID, err = strconv.ParseUint(alert.Labels.TemplateId, 10, 64)
		if err != nil {
			zap.L().Error("parse int error", zap.Error(err))
		}
		alertLog.DeviceDataID, err = strconv.ParseUint(alert.Labels.DataId, 10, 64)
		if err != nil {
			zap.L().Error("parse int error", zap.Error(err))
		}
		alertLog.Context = alert.Labels.Context
		alertLog.Message = alert.Labels.Message
		alertLog.StartTimeUnix = gtime.New(alert.StartsAt)
		alertLog.TimeUnix = gtime.New(alert.StartsAt)
		ret = append(ret, &alertLog)
		infos = append(infos, info)
	}

	return ret, infos
}

type NovaAlertLogList struct {
	Rows  []*NovaAlertLog `json:"rows"`
	Total int64           `json:"total"`
}
