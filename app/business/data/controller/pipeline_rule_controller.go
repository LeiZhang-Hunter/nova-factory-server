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

type PipelineRuleController struct{ service service.IPipelineRuleService }

func NewPipelineRuleController(service service.IPipelineRuleService) *PipelineRuleController {
	return &PipelineRuleController{service: service}
}

func (ctrl *PipelineRuleController) PrivateRoutes(router *gin.RouterGroup) {
	rules := router.Group("/data/pipeline-rules")
	rules.POST("", middlewares.SetLog("创建 Pipeline 规则", middlewares.Insert), rulePermission("add"), ctrl.Create)
	rules.GET("", rulePermission("query"), ctrl.List)
	rules.GET("/:id", rulePermission("query"), ctrl.Get)
	rules.PUT("/:id", middlewares.SetLog("更新 Pipeline 规则", middlewares.Update), rulePermission("edit"), ctrl.Update)
	rules.DELETE("/:id", middlewares.SetLog("删除 Pipeline 规则", middlewares.Delete), rulePermission("remove"), ctrl.Delete)
}

func (ctrl *PipelineRuleController) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("POST", "/data/pipeline-rules", "data:service:rule:add")
	router.RegisterPermission("GET", "/data/pipeline-rules", "data:service:rule:query")
	router.RegisterPermission("GET", "/data/pipeline-rules/:id", "data:service:rule:query")
	router.RegisterPermission("PUT", "/data/pipeline-rules/:id", "data:service:rule:edit")
	router.RegisterPermission("DELETE", "/data/pipeline-rules/:id", "data:service:rule:remove")
}

func (ctrl *PipelineRuleController) Create(c *gin.Context) {
	var req dto.CreatePipelineRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	result, err := ctrl.service.Create(c, &req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, result)
}

func (ctrl *PipelineRuleController) Get(c *gin.Context) {
	result, err := ctrl.service.Get(c, c.Param("id"))
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, result)
}

func (ctrl *PipelineRuleController) List(c *gin.Context) {
	pageNum := cast.ToInt(c.DefaultQuery("pageNum", "1"))
	pageSize := cast.ToInt(c.DefaultQuery("pageSize", "10"))
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 10
	}
	rows, total, err := ctrl.service.List(c, (pageNum-1)*pageSize, pageSize, c.Query("sourceType"), c.Query("name"), c.Query("status"))
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessListData(c, rows, total)
}

func (ctrl *PipelineRuleController) Update(c *gin.Context) {
	var req dto.UpdatePipelineRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	result, err := ctrl.service.Update(c, c.Param("id"), &req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, result)
}

func (ctrl *PipelineRuleController) Delete(c *gin.Context) {
	if err := ctrl.service.Delete(c, c.Param("id")); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

type ServiceConnectionController struct {
	service service.IServiceConnectionService
}

func NewServiceConnectionController(service service.IServiceConnectionService) *ServiceConnectionController {
	return &ServiceConnectionController{service: service}
}

func (ctrl *ServiceConnectionController) PrivateRoutes(router *gin.RouterGroup) {
	connections := router.Group("/data/service-connections")
	connections.POST("", middlewares.SetLog("创建服务连接", middlewares.Insert), connectionPermission("add"), ctrl.Create)
	connections.GET("", connectionPermission("query"), ctrl.List)
	connections.GET("/:id", connectionPermission("query"), ctrl.Get)
	connections.PUT("/:id", middlewares.SetLog("更新服务连接", middlewares.Update), connectionPermission("edit"), ctrl.Update)
	connections.DELETE("/:id", middlewares.SetLog("删除服务连接", middlewares.Delete), connectionPermission("remove"), ctrl.Delete)
}

func rulePermission(action string) gin.HandlerFunc {
	return middlewares.HasPermissions([]string{"data:service:rule:" + action, "data:service:rule"})
}

func connectionPermission(action string) gin.HandlerFunc {
	return middlewares.HasPermissions([]string{"data:service:connection:" + action, "data:service:rule"})
}

func (ctrl *ServiceConnectionController) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("POST", "/data/service-connections", "data:service:connection:add")
	router.RegisterPermission("GET", "/data/service-connections", "data:service:connection:query")
	router.RegisterPermission("GET", "/data/service-connections/:id", "data:service:connection:query")
	router.RegisterPermission("PUT", "/data/service-connections/:id", "data:service:connection:edit")
	router.RegisterPermission("DELETE", "/data/service-connections/:id", "data:service:connection:remove")
}

func (ctrl *ServiceConnectionController) Create(c *gin.Context) {
	var req dto.CreateServiceConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	result, err := ctrl.service.Create(c, &req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, result)
}

func (ctrl *ServiceConnectionController) Get(c *gin.Context) {
	result, err := ctrl.service.Get(c, c.Param("id"))
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, result)
}

func (ctrl *ServiceConnectionController) List(c *gin.Context) {
	pageNum := cast.ToInt(c.DefaultQuery("pageNum", "1"))
	pageSize := cast.ToInt(c.DefaultQuery("pageSize", "10"))
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 10
	}
	rows, total, err := ctrl.service.List(c, (pageNum-1)*pageSize, pageSize, c.Query("sourceType"), c.Query("name"), c.Query("status"))
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessListData(c, rows, total)
}

func (ctrl *ServiceConnectionController) Update(c *gin.Context) {
	var req dto.UpdateServiceConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	result, err := ctrl.service.Update(c, c.Param("id"), &req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, result)
}

func (ctrl *ServiceConnectionController) Delete(c *gin.Context) {
	if err := ctrl.service.Delete(c, c.Param("id")); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
