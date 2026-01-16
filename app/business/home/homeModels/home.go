package homeModels

import (
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/business/monitor/monitorModels"
)

type DeviceStats struct {
	TotalDevices int64 `json:"total_devices"`
	// OnlineDevices 上线设备数
	OnlineDevices int64 `json:"onlineDevices"`
	// OfflineDevices 下线设备数
	OfflineDevices int64 `json:"offlineDevices"`
	// ExceptionCount 异常统计
	ExceptionCount int64 `json:"exceptionCount"`
	// MaintenanceCount 维护设备数
	MaintenanceCount int64 `json:"maintenanceCount"`
}

type HomeStats struct {
	DeviceStats DeviceStats `json:"deviceStats"`
	// BuildingCount 建筑统计
	BuildingCount int64 `json:"buildingCount"`
	// ScheduleCount 调度计划数
	ScheduleCount int64 `json:"scheduleCount"`
	// Alerts 告警策略
	Alerts []*alertModels.NovaAlertLog `json:"alerts"`
	// Alerts 告警数量
	AlertsCount int64 `json:"alertsCount"`
	// Monitor 监控服务
	Monitor *monitorModels.Server `json:"monitor"`
	// DeviceCounter 设备统计
	DeviceCounter *metricModels.MetricQueryData `json:"deviceCounter"`
	// DeviceCounterRank 根据设备写入总数排行
	DeviceCounterRank *deviceMonitorModel.TypeDeviceCounterRank `json:"deviceCounterRank"`
}
