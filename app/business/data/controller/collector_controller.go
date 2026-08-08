package controller

import (
	"nova-factory-server/app/business/data/models/dto"
	"nova-factory-server/app/business/data/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// CollectorController 采集器控制器
type CollectorController struct {
	collectorService service.ICollectorService
}

// NewCollectorController 创建控制器
func NewCollectorController(collectorService service.ICollectorService) *CollectorController {
	return &CollectorController{collectorService: collectorService}
}

// PrivateRoutes 注册私有路由
func (ctrl *CollectorController) PrivateRoutes(router *gin.RouterGroup) {
	collectors := router.Group("/data/collectors")
	collectors.POST("", middlewares.SetLog("创建采集器", middlewares.Insert), middlewares.HasPermission("data:collector:add"), ctrl.CreateCollector)
	collectors.GET("", middlewares.HasPermission("data:collector:query"), ctrl.ListCollectors)
	collectors.GET("/online", middlewares.HasPermission("data:collector:query"), ctrl.ListOnlineCollectors)
	collectors.GET("/:id", middlewares.HasPermission("data:collector:query"), ctrl.GetCollector)
	collectors.PUT("/:id", middlewares.SetLog("更新采集器", middlewares.Update), middlewares.HasPermission("data:collector:edit"), ctrl.UpdateCollector)
	collectors.DELETE("/:id", middlewares.SetLog("删除采集器", middlewares.Delete), middlewares.HasPermission("data:collector:remove"), ctrl.DeleteCollector)
}

// PrivateMcpRoutes MCP 权限注册
func (ctrl *CollectorController) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("POST", "/data/collectors", "data:collector:add")
	router.RegisterPermission("GET", "/data/collectors", "data:collector:query")
	router.RegisterPermission("GET", "/data/collectors/online", "data:collector:query")
	router.RegisterPermission("GET", "/data/collectors/:id", "data:collector:query")
	router.RegisterPermission("PUT", "/data/collectors/:id", "data:collector:edit")
	router.RegisterPermission("DELETE", "/data/collectors/:id", "data:collector:remove")
}

// CreateCollector 创建采集器
// @Summary 创建采集器
// @Description 创建采集器设备
// @Tags 数据平台-采集器管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateCollectorRequest true "创建采集器请求"
// @Success 200 {object} response.ResponseData{data=dto.Collector}
// @Router /data/collectors [post]
func (ctrl *CollectorController) CreateCollector(c *gin.Context) {
	var req dto.CreateCollectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	result, err := ctrl.collectorService.CreateCollector(c, &req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	baizeContext.SuccessData(c, result)
}

// GetCollector 获取采集器详情
// @Summary 获取采集器详情
// @Description 获取采集器设备详情
// @Tags 数据平台-采集器管理
// @Security BearerAuth
// @Produce json
// @Param id path string true "采集器ID"
// @Success 200 {object} response.ResponseData{data=dto.Collector}
// @Router /data/collectors/{id} [get]
func (ctrl *CollectorController) GetCollector(c *gin.Context) {
	id := c.Param("id")

	result, err := ctrl.collectorService.GetCollector(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	baizeContext.SuccessData(c, result)
}

// ListCollectors 获取采集器列表
// @Summary 获取采集器列表
// @Description 获取采集器设备列表
// @Tags 数据平台-采集器管理
// @Security BearerAuth
// @Produce json
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param name query string false "设备名称"
// @Param deviceId query string false "设备ID"
// @Param status query string false "状态"
// @Success 200 {object} response.ResponseData{data=[]dto.Collector}
// @Router /data/collectors [get]
func (ctrl *CollectorController) ListCollectors(c *gin.Context) {
	pageNum := cast.ToInt(c.DefaultQuery("pageNum", "1"))
	pageSize := cast.ToInt(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")
	deviceID := c.Query("deviceId")
	status := c.Query("status")

	offset := (pageNum - 1) * pageSize
	list, total, err := ctrl.collectorService.ListCollectors(c, offset, pageSize, name, deviceID, status)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	baizeContext.SuccessListData(c, list, total)
}

// ListOnlineCollectors 获取在线采集器列表（供下发弹窗使用）
// @Summary 获取在线采集器
// @Description 获取全部在线采集器
// @Tags 数据平台-采集器管理
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.ResponseData{data=[]dto.Collector}
// @Router /data/collectors/online [get]
func (ctrl *CollectorController) ListOnlineCollectors(c *gin.Context) {
	list, err := ctrl.collectorService.ListAllOnline(c)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	baizeContext.SuccessData(c, list)
}

// UpdateCollector 更新采集器
// @Summary 更新采集器
// @Description 更新采集器设备信息
// @Tags 数据平台-采集器管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "采集器ID"
// @Param body body dto.UpdateCollectorRequest true "更新采集器请求"
// @Success 200 {object} response.ResponseData{data=dto.Collector}
// @Router /data/collectors/{id} [put]
func (ctrl *CollectorController) UpdateCollector(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateCollectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	result, err := ctrl.collectorService.UpdateCollector(c, id, &req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	baizeContext.SuccessData(c, result)
}

// DeleteCollector 删除采集器
// @Summary 删除采集器
// @Description 删除采集器设备
// @Tags 数据平台-采集器管理
// @Security BearerAuth
// @Produce json
// @Param id path string true "采集器ID"
// @Success 200 {object} response.ResponseData
// @Router /data/collectors/{id} [delete]
func (ctrl *CollectorController) DeleteCollector(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.collectorService.DeleteCollector(c, id); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	baizeContext.Success(c)
}
