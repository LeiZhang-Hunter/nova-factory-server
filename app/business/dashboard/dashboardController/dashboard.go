package dashboardController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/dashboard/dashboardModels"
	"nova-factory-server/app/business/dashboard/dashboardService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Dashboard struct {
	service dashboardService.DashboardService
}

func NewDashboard(service dashboardService.DashboardService) *Dashboard {
	return &Dashboard{
		service: service,
	}
}

func (d *Dashboard) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/dashboard/manager")
	group.GET("/list", middlewares.HasPermission("dashboard:manager:list"), d.List)
	group.POST("/set", middlewares.HasPermission("dashboard:manager:set"), d.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("dashboard:manager:remove"), d.Remove)
}

// List 仪表盘列表
// @Summary 仪表盘列表
// @Description 仪表盘列表
// @Tags 仪表盘/仪表盘管理
// @Param  object query dashboardModels.SysDashboardReq true "组成工序列表参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /dashboard/manager/list [get]
func (d *Dashboard) List(c *gin.Context) {
	req := new(dashboardModels.SysDashboardReq)
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
	return
}

// Set 保存仪表盘
// @Summary 保存仪表盘
// @Description 保存仪表盘
// @Tags 仪表盘/仪表盘管理
// @Param  object body dashboardModels.SetSysDashboard true "设置仪表盘参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /dashboard/manager/set [post]
func (d *Dashboard) Set(c *gin.Context) {
	req := new(dashboardModels.SetSysDashboard)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	value, err := d.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, value)
}

// Remove 删除仪表盘
// @Summary 删除仪表盘
// @Description 删除仪表盘
// @Tags 仪表盘/仪表盘管理
// @Param  ids path string true "ids"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "删除仪表盘参数"
// @Router /dashboard/manager/remove [delete]
func (d *Dashboard) Remove(c *gin.Context) {
	recordIds := baizeContext.ParamStringArray(c, "ids")
	if len(recordIds) == 0 {
		baizeContext.Waring(c, "请选择调度任务")
		return
	}

	err := d.service.Remove(c, recordIds)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}
