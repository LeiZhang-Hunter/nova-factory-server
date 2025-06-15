package deviceMonitorController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type DeviceMonitor struct {
	service deviceMonitorService.DeviceMonitorService
}

func NewDeviceMonitor(service deviceMonitorService.DeviceMonitorService) *DeviceMonitor {
	return &DeviceMonitor{
		service: service,
	}
}

func (d *DeviceMonitor) PrivateRoutes(router *gin.RouterGroup) {
	monitor := router.Group("/device/monitor")
	monitor.GET("/list", middlewares.HasPermission("device:monitor:list"), d.List)
	monitor.POST("/metric", middlewares.HasPermission("device:monitor:metric"), d.Metric)
}

// List 设备监控
// @Summary 设备监控
// @Description 设备监控
// @Tags 设备监控/设备监控
// @Param  object query daemonizeModels.SysIotAgentListReq true "设备监控"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设备监控"
// @Router /device/monitor/list [get]
func (d *DeviceMonitor) List(c *gin.Context) {
	list, err := d.service.List(c)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Metric 设备指标
// @Summary 设备指标
// @Description 设备指标
// @Tags 设备监控/设备监控
// @Param  object body metricModels.MetricQueryReq true "设备监控"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设备监控"
// @Router /device/monitor/metric [post]
func (d *DeviceMonitor) Metric(c *gin.Context) {
	req := new(metricModels.MetricQueryReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := d.service.Metric(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}
