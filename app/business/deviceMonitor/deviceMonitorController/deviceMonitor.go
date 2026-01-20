package deviceMonitorController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"
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
	monitor.POST("/predict", middlewares.HasPermission("device:monitor:metric"), d.Predict)

}

func (d *DeviceMonitor) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1")
	group.GET("/metric/predict", d.MetricPredict)
	group.POST("/metric/predict/query", d.MetricPredictQuery)
	group.GET("/device/layout", d.DeviceLayout)
	group.GET("/device/building", d.DeviceByBuilding)
}

func (d *DeviceMonitor) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterSchema("POST", "/api/v1/metric/predict/query", nil, metricModels.GatewayMetricDataQueryReq{})
}

// List 设备监控
// @Summary 设备监控
// @Description 设备监控
// @Tags 设备监控/设备监控
// @Param  object query deviceModels.DeviceListReq true "设备监控"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设备监控"
// @Router /device/monitor/list [get]
func (d *DeviceMonitor) List(c *gin.Context) {
	req := new(deviceModels.DeviceListReq)
	err := c.ShouldBindQuery(req)
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

// Predict 预测指标
// @Summary 预测指标
// @Description 预测指标
// @Tags 设备监控/设备监控
// @Param  object body metricModels.MetricQueryReq true "设备监控"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设备监控"
// @Router /device/monitor/predict [post]
func (d *DeviceMonitor) Predict(c *gin.Context) {
	req := new(metricModels.MetricQueryReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := d.service.Predict(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// MetricPredict 指标预测
// @Summary 导入告警数据
// @Description 导入告警数据
// @Tags 设备监控/设备监控
// @Param object body alertModels.AlertLogData true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /api/v1/metric/predict [post]
func (d *DeviceMonitor) MetricPredict(c *gin.Context) {
	req := new(metricModels.MetricQueryReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	data, err := d.service.Predict(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// MetricPredictQuery 指标预测查询
// @Summary 指标预测查询
// @Description 指标预测查询
// @Tags 设备监控/设备监控
// @Param object body metricModels.GatewayMetricDataQueryReq true "指标预测查询参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /api/v1/metric/predict/query [post]
func (d *DeviceMonitor) MetricPredictQuery(c *gin.Context) {
	req := new(metricModels.GatewayMetricDataQueryReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("param error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}

	data, err := d.service.PredictQuery(c, &metricModels.MetricDataQueryReq{
		Type:        req.Type,
		Name:        req.Name,
		Start:       req.Start,
		End:         req.End,
		Step:        req.Step,
		Interval:    req.Interval,
		Level:       req.Level,
		Expression:  req.Expression,
		Field:       req.Field,
		Predict:     req.Predict,
		Having:      req.Having,
		QueryMetric: req.QueryMetric,
	})
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// DeviceLayout 设备布局
// @Summary 设备布局
// @Description 设备布局
// @Tags 设备监控/设备监控
// @Param object body metricModels.GatewayMetricDataQueryReq true "指标预测查询参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /api/v1/device/layout [get]
func (d *DeviceMonitor) DeviceLayout(c *gin.Context) {
	req := new(deviceMonitorModel.DeviceLayoutRequest)
	err := c.ShouldBindQuery(req)
	if err != nil {
		zap.L().Error("param error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}

	layout, err := d.service.DeviceLayout(c, req.FloorId)
	if err != nil {
		zap.L().Error("device layout error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}

	baizeContext.SuccessData(c, layout)
}

// DeviceByBuilding 通过建筑物读取设备
// @Summary 通过建筑物读取设备
// @Description 通过建筑物读取设备
// @Tags 设备监控/通过建筑物读取设备
// @Param object body metricModels.GatewayMetricDataQueryReq true "指标预测查询参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /api/v1/device/building [get]
func (d *DeviceMonitor) DeviceByBuilding(c *gin.Context) {
	req := new(deviceMonitorModel.DeviceLayoutRequest)
	err := c.ShouldBindQuery(req)
	if err != nil {
		zap.L().Error("param error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}

	layout, err := d.service.DeviceLayout(c, req.FloorId)
	if err != nil {
		zap.L().Error("device layout error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}

	baizeContext.SuccessData(c, layout)
}
