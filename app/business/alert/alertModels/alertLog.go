package alertModels

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/utils/snowflake"
	"strconv"
	"time"
)

type AlertLogData struct {
	Id     uint64 `json:"id"`
	Alerts []struct {
		Labels struct {
			DeviceId   string `json:"device_id"`
			TemplateId string `json:"template_id"`
			DataId     string `json:"data_id"`
			AlertId    string `json:"alert_id"`
			GroupId    string `json:"group_id"`
			Context    string `json:"context"`
			Message    string `json:"message"`
		} `json:"labels"`
		Annotations struct {
			Reason string `json:"reason"`
		} `json:"annotations"`
		StartsAt time.Time `json:"startsAt"`
		EndsAt   time.Time `json:"endsAt"`
	} `json:"alerts"`
	CommonLabels map[string]string `json:"commonLabels"`
	GatewayId    int64             `json:"gatewayId"`
	Username     string            `json:"username"`
	Password     string            `json:"password"`
}

// SysAlertLog 告警日志
type SysAlertLog struct {
	ID               int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id"`              // 自增标识
	GatewayID        int64  `gorm:"column:gateway_id;not null;comment:网关id" json:"gateway_id"`                   // 网关id
	AlertID          int64  `gorm:"column:alert_id;not null;comment:告警id" json:"alert_id"`                       // 告警id
	DeviceID         int64  `gorm:"column:device_id;not null;comment:设备id" json:"device_id"`                     // 设备id
	DeviceTemplateID int64  `gorm:"column:device_template_id;not null;comment:设备模板id" json:"device_template_id"` // 设备模板id
	DeviceDataID     int64  `gorm:"column:device_data_id;not null;comment:设备数据id" json:"device_data_id"`         // 设备数据id
	Context          string `gorm:"column:context;comment:内容" json:"context"`                                    // 内容
	Message          string `gorm:"column:message;comment:消息" json:"message"`                                    // 消息
	Data             string `gorm:"column:data;comment:消息快照" json:"data"`                                        // 消息快照
	DeptID           int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                  // 部门ID
	State            bool   `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                            // 操作状态（0正常 -1删除）
	baize.BaseEntity
}

func FromDataToSysAlertLog(data *AlertLogData, deptId uint64, c *gin.Context) []*SysAlertLog {
	alertJson, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	ret := make([]*SysAlertLog, 0)
	for _, alert := range data.Alerts {
		var alertLog SysAlertLog
		alertLog.ID = snowflake.GenID()
		alertLog.Data = string(alertJson)
		alertLog.GatewayID = data.GatewayId
		alertLog.AlertID = int64(data.Id)
		alertLog.DeviceID, err = strconv.ParseInt(alert.Labels.DeviceId, 10, 64)
		if err != nil {
			zap.L().Error("parse int error", zap.Error(err))
		}
		alertLog.DeviceTemplateID, err = strconv.ParseInt(alert.Labels.TemplateId, 10, 64)
		if err != nil {
			zap.L().Error("parse int error", zap.Error(err))
		}
		alertLog.DeviceDataID, err = strconv.ParseInt(alert.Labels.DataId, 10, 64)
		if err != nil {
			zap.L().Error("parse int error", zap.Error(err))
		}
		alertLog.DeptID = int64(deptId)
		alertLog.Context = alert.Labels.Context
		alertLog.Message = alert.Labels.Message
		alertLog.CreateTime = &alert.StartsAt
		alertLog.UpdateTime = &alert.StartsAt
		ret = append(ret, &alertLog)
	}

	return ret
}

type SysAlertLogListReq struct {
	GatewayID int64  `json:"gateway_id"` // 网关id
	AlertID   int64  `json:"alert_id"`   // 告警id
	Message   string `json:"message"`
	baize.BaseEntityDQL
}

type SysAlertLogList struct {
	Rows  []*SysAlertLog `json:"rows"`
	Total int64          `json:"total"`
}
