package deviceMonitorController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
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
