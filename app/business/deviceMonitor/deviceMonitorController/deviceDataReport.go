package deviceMonitorController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/business/metric/device/metricService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type DeviceReport struct {
	service             metricService.IMetricService
	deviceReportService deviceMonitorService.IDeviceDataReportService
}

func NewDeviceReport(service metricService.IMetricService, deviceReportService deviceMonitorService.IDeviceDataReportService) *DeviceReport {
	return &DeviceReport{
		service:             service,
		deviceReportService: deviceReportService,
	}
}

func (d *DeviceReport) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/device/monitor/data")
	group.GET("/dev/list", middlewares.HasPermission("device:monitor:data:dev:list"), d.DevList)
	group.POST("/list", middlewares.HasPermission("device:monitor:data:list"), d.List)
}

// DevList 设备测点列表
// @Summary 设备测点列表
// @Description 设备测点列表
// @Tags 设备监控/设备监控
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /device/monitor/data/dev/list [get]
func (d *DeviceReport) DevList(c *gin.Context) {
	list, err := d.deviceReportService.DevList(c)
	if err != nil {
		return
	}
	baizeContext.SuccessData(c, list)
}

// List 实时数据
// @Summary 实时数据
// @Description 实时数据
// @Tags 设备监控/设备监控
// @Param  object body deviceMonitorModel.DevDataReq true "获取模板数据列表请求参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /device/monitor/data/list [get]
func (d *DeviceReport) List(c *gin.Context) {
	req := new(deviceMonitorModel.DevDataReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := d.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

func (d *DeviceReport) Export(router *gin.Context) {

}
