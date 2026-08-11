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

const (
	// HeaderConfigMD5 采集器拉取时携带的当前配置 MD5
	HeaderConfigMD5 = "X-Config-MD5"
)

// CollectorRuleController 采集器规则下发控制器
type CollectorRuleController struct {
	service service.ICollectorRuleService
}

// NewCollectorRuleController 创建采集器规则下发控制器
func NewCollectorRuleController(service service.ICollectorRuleService) *CollectorRuleController {
	return &CollectorRuleController{service: service}
}

// PrivateRoutes 注册私有路由
func (ctrl *CollectorRuleController) PrivateRoutes(router *gin.RouterGroup) {
	router.POST("/data/collectors/:id/rules", middlewares.SetLog("下发规则", middlewares.Insert), middlewares.HasPermission("data:dispatch:add"), ctrl.Dispatch)
	router.POST("/data/collectors/:id/rules/preview", middlewares.HasPermission("data:dispatch:query"), ctrl.PreviewProposed)
	router.GET("/data/collectors/:id/rules/preview", middlewares.HasPermission("data:dispatch:query"), ctrl.PreviewCurrent)
	router.GET("/data/dispatch-logs", middlewares.HasPermission("data:dispatch:query"), ctrl.ListLogs)
}

// PrivateMcpRoutes MCP 权限注册
func (ctrl *CollectorRuleController) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("POST", "/data/collectors/:id/rules", "data:dispatch:add")
	router.RegisterPermission("POST", "/data/collectors/:id/rules/preview", "data:dispatch:query")
	router.RegisterPermission("GET", "/data/collectors/:id/rules/preview", "data:dispatch:query")
	router.RegisterPermission("GET", "/data/dispatch-logs", "data:dispatch:query")
}

// Dispatch 给采集器下发规则
// @Summary 给采集器下发规则
// @Description 将多选规则的 pipelines 合并为一条配置快照并写入下发记录
// @Tags 数据平台-采集器下发
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "采集器ID"
// @Param body body dto.BindCollectorRulesRequest true "下发规则请求"
// @Success 200 {object} response.ResponseData{data=dto.DispatchLogResponse}
// @Router /data/collectors/{id}/rules [post]
func (ctrl *CollectorRuleController) Dispatch(c *gin.Context) {
	var req dto.BindCollectorRulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	result, err := ctrl.service.Dispatch(c, c.Param("id"), baizeContext.GetUserName(c), &req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, result)
}

// PreviewProposed 拟下发集合预览
// @Summary 拟下发集合预览
// @Description 提交拟下发规则集合，预览合并后的配置与校验和（不落库不写日志）
// @Tags 数据平台-采集器下发
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "采集器ID"
// @Param body body dto.BindCollectorRulesRequest true "拟下发规则请求"
// @Success 200 {object} response.ResponseData{data=dto.DispatchPreviewResponse}
// @Router /data/collectors/{id}/rules/preview [post]
func (ctrl *CollectorRuleController) PreviewProposed(c *gin.Context) {
	var req dto.BindCollectorRulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	result, err := ctrl.service.PreviewProposed(c, c.Param("id"), &req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, result)
}

// PreviewCurrent 当前已下发规则预览
// @Summary 当前已下发规则预览
// @Description 预览采集器最近一次下发的合并配置（等价于设备下次拉取内容）
// @Tags 数据平台-采集器下发
// @Security BearerAuth
// @Produce json
// @Param id path string true "采集器ID"
// @Success 200 {object} response.ResponseData{data=dto.DispatchPreviewResponse}
// @Router /data/collectors/{id}/rules/preview [get]
func (ctrl *CollectorRuleController) PreviewCurrent(c *gin.Context) {
	result, err := ctrl.service.PreviewCurrent(c, c.Param("id"))
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, result)
}

// PullRules 采集器拉取最新下发规则
// @Summary 采集器拉取最新下发规则
// @Description 采集器定时调用，携带当前配置 MD5 获取最新下发规则，并刷新在线状态
// @Tags 数据平台-采集器下发
// @Produce json
// @Param X-Device-Id header string true "设备ID"
// @Param X-Device-Token header string true "设备Token"
// @Param X-Config-MD5 header string false "当前配置 MD5"
// @Success 200 {object} response.ResponseData{data=dto.PullRulesResponse}
// @Router /data/collectors/rules [get]
func (ctrl *CollectorRuleController) PullRules(c *gin.Context) {
	deviceID := c.GetHeader(middlewares.HeaderDeviceID)
	token := c.GetHeader(middlewares.HeaderDeviceToken)
	clientMD5 := c.GetHeader(HeaderConfigMD5)
	result, err := ctrl.service.PullRules(c, deviceID, token, clientMD5)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, result)
}

// ListLogs 下发版本记录
// @Summary 下发版本记录
// @Description 分页查询规则下发记录
// @Tags 数据平台-采集器下发
// @Security BearerAuth
// @Produce json
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param collectorId query string false "采集器ID"
// @Success 200 {object} response.ResponseData{data=[]dto.DispatchLogResponse}
// @Router /data/dispatch-logs [get]
func (ctrl *CollectorRuleController) ListLogs(c *gin.Context) {
	pageNum := cast.ToInt(c.DefaultQuery("pageNum", "1"))
	pageSize := cast.ToInt(c.DefaultQuery("pageSize", "10"))
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 10
	}
	rows, total, err := ctrl.service.ListLogs(c, (pageNum-1)*pageSize, pageSize, c.Query("collectorId"))
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessListData(c, rows, total)
}
