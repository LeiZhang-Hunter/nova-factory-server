package deviceMonitorController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/middlewares"
)

// DeviceUtilization 计算设备稼动率
type DeviceUtilization struct {
	service deviceMonitorService.DeviceUtilizationService
}

func NewDeviceUtilization(service deviceMonitorService.DeviceUtilizationService) *DeviceUtilization {
	return &DeviceUtilization{
		service: service,
	}
}

func (d *DeviceUtilization) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/device/monitor/utilization")
	group.POST("/stat", middlewares.HasPermission("device:monitor:utilization:stat"), d.Stat)
}

// Stat 稼动率统计
// @Summary 稼动率统计
// @Description 稼动率统计
// @Tags 设备监控/设备监控
// @Param  object body metricModels.MetricQueryReq true "稼动率统计"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设备监控"
// @Router /device/monitor/utilization/stat [post]
func (d *DeviceUtilization) Stat(c *gin.Context) {
	d.service.Stat(c, &deviceMonitorModel.DeviceUtilizationReq{})
}
