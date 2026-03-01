package deviceMonitorController

import (
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
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
	group.GET("/info", middlewares.HasPermission("device:monitor:utilization:info"), d.Info)
}

func (d *DeviceUtilization) PublicRoutes(router *gin.RouterGroup) {
	ele := router.Group("/api/v1/system/monitor")
	ele.GET("/utilization/search", d.Search)
	ele.GET("/utilization/searchv2", d.SearchV2)
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
	req := new(deviceMonitorModel.DeviceUtilizationReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	stat, err := d.service.Stat(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, stat)
}

// Search 能源大屏
// @Summary 能源大屏
// @Description 能源大屏
// @Tags 设备监控/设备监控
// @Param  object body metricModels.MetricQueryReq true "能源大屏"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设备监控"
// @Router /api/v1/system/electric/search [post]
func (d *DeviceUtilization) Search(c *gin.Context) {
	req := new(deviceMonitorModel.DeviceUtilizationReq)
	stat, err := d.service.Search(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, stat)
}

// SearchV2 能源大屏V2
// @Summary 能源大屏V2
// @Description 能源大屏V2
// @Tags 设备监控/设备监控
// @Param  object body metricModels.MetricQueryReq true "能源大屏"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设备监控"
// @Router /api/v1/system/electric/searchv2 [post]
func (d *DeviceUtilization) SearchV2(c *gin.Context) {
	req := new(deviceMonitorModel.DeviceUtilizationReq)
	stat, err := d.service.SearchV2(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, stat)
}

// Info 设备稼动率详情
// @Summary 设备稼动率详情
// @Description 设备稼动率详情
// @Tags 设备监控/设备监控
// @Param DeviceId query int true "设备id"
// @Param Start query int false "开始时间"
// @Param End query int false "结束时间"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设备监控"
// @Router /device/monitor/utilization/info [get]
func (d *DeviceUtilization) Info(c *gin.Context) {
	req := new(deviceMonitorModel.DeviceUtilizationReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	stat, err := d.service.GetDeviceUtilization(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, stat)
}
