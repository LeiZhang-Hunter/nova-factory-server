package dashboardcontroller

import (
	"nova-factory-server/app/business/iot/dashboard/dashboardmodels"
	"nova-factory-server/app/business/iot/dashboard/dashboardservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Data struct {
	service dashboardservice.DashboardDataService
}

func NewData(service dashboardservice.DashboardDataService) *Data {
	return &Data{
		service: service,
	}
}

func (d *Data) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/dashboard/data")
	//group.GET("/list", middlewares.HasPermission("dashboard:manager:list"), d.List)
	group.POST("/set", middlewares.HasPermission("dashboard:data:set"), d.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("dashboard:data:remove"), d.Remove)
	group.GET("/info", middlewares.HasPermission("dashboard:data:info"), d.Info)
}

// Set 保存面板
// @Summary 保存面板
// @Description 保存面板
// @Tags 仪表盘/面板管理
// @Param  object body dashboardmodels.SetSysDashboardData true "设置仪表盘参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /dashboard/data/set [post]
func (d *Data) Set(c *gin.Context) {
	req := new(dashboardmodels.SetSysDashboardData)
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

// Remove 删除面板
// @Summary 删除面板
// @Description 删除面板
// @Tags 仪表盘/面板管理
// @Param  ids path string true "ids"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "删除面板"
// @Router /dashboard/data/remove [delete]
func (d *Data) Remove(c *gin.Context) {
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

// Info 读取面板
// @Summary 读取面板
// @Description 读取面板
// @Tags 仪表盘/面板管理
// @Param  object query dashboardmodels.GetSysDashboardData true "设置仪表盘参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /dashboard/data/info [get]
func (d *Data) Info(c *gin.Context) {
	req := new(dashboardmodels.GetSysDashboardData)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	value, err := d.service.Info(c, req.DatashboardID)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if value == nil {
		baizeContext.SuccessData(c, &dashboardmodels.SysDashboardData{
			ID:            0,
			Data:          "[]",
			DatashboardID: req.DatashboardID,
		})
		return
	}
	baizeContext.SuccessData(c, value)
}
