package controller

import (
	"nova-factory-server/app/business/data/models/dto"
	"nova-factory-server/app/business/data/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// ModelConfigController 数据平台模型配置控制器。
type ModelConfigController struct {
	service service.IModelConfigService
}

// NewModelConfigController 创建模型配置控制器。
func NewModelConfigController(modelConfigService service.IModelConfigService) *ModelConfigController {
	return &ModelConfigController{service: modelConfigService}
}

// PrivateRoutes 注册私有路由。
func (ctrl *ModelConfigController) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/data/config/model")
	group.GET("", middlewares.HasPermission("data:config:model:query"), ctrl.Get)
	group.PUT("", middlewares.SetLog("更新模型配置", middlewares.Update), middlewares.HasPermission("data:config:model:edit"), ctrl.Update)
}

// PrivateMcpRoutes MCP 权限注册。
func (ctrl *ModelConfigController) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/data/config/model", "data:config:model:query")
	router.RegisterPermission("PUT", "/data/config/model", "data:config:model:edit")
}

// Get 获取模型配置。
// @Summary 获取模型配置
// @Description 获取数据平台模型配置，未保存时返回默认值
// @Tags 数据平台-配置管理
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.ResponseData{data=dto.ModelConfigResponse}
// @Router /data/config/model [get]
func (ctrl *ModelConfigController) Get(c *gin.Context) {
	result, err := ctrl.service.Get(c)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, result)
}

// Update 更新模型配置。
// @Summary 更新模型配置
// @Description 更新数据平台模型配置
// @Tags 数据平台-配置管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.UpdateModelConfigRequest true "更新模型配置请求"
// @Success 200 {object} response.ResponseData{data=dto.ModelConfigResponse}
// @Router /data/config/model [put]
func (ctrl *ModelConfigController) Update(c *gin.Context) {
	var req dto.UpdateModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	result, err := ctrl.service.Update(c, &req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, result)
}
