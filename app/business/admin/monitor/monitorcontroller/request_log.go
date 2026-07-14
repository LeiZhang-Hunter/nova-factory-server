package monitorcontroller

import (
	"fmt"
	"time"

	"nova-factory-server/app/business/admin/monitor/monitormodels"
	"nova-factory-server/app/business/admin/monitor/monitorservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

type RequestLog struct {
	service monitorservice.IRequestLogService
}

func NewRequestLog(service monitorservice.IRequestLogService) *RequestLog {
	return &RequestLog{service: service}
}

func (c *RequestLog) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/monitor/request-log")
	group.GET("/list", middlewares.HasPermission("system:monitor:requestlog"), c.List)
	group.GET("/detail/:id", middlewares.HasPermission("system:monitor:requestlog"), c.Detail)
	group.DELETE("/clean", middlewares.HasPermission("system:monitor:requestlog"), c.Clean)
}

func (c *RequestLog) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/monitor/request-log/list", "system:monitor:requestlog")
	router.RegisterPermission("GET", "/monitor/request-log/detail/:id", "system:monitor:requestlog")
	router.RegisterPermission("DELETE", "/monitor/request-log/clean", "system:monitor:requestlog")
}

func (c *RequestLog) List(ctx *gin.Context) {
	req := new(monitormodels.RequestLogQuery)
	if err := ctx.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := c.service.List(ctx, req)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

func (c *RequestLog) Detail(ctx *gin.Context) {
	id := baizeContext.ParamInt64(ctx, "id")
	if id == 0 {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := c.service.Detail(ctx, id)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

func (c *RequestLog) Clean(ctx *gin.Context) {
	beforeTime := ctx.Query("beforeTime")
	if beforeTime == "" {
		beforeTime = time.Now().AddDate(0, 0, -30).Format("2006-01-02 15:04:05")
	}
	if err := c.service.Clean(ctx, beforeTime); err != nil {
		baizeContext.Waring(ctx, fmt.Sprintf("clean failed: %v", err))
		return
	}
	baizeContext.Success(ctx)
}
