package craftRouteController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Schedule struct {
	service craftRouteService.IScheduleService
}

func NewSchedule(service craftRouteService.IScheduleService) *Schedule {
	return &Schedule{
		service: service,
	}
}

func (schedule *Schedule) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/craft/route/schedule")
	routers.GET("/list", middlewares.HasPermission("craft:route:schedule:list"), schedule.List) // 调度列表
	routers.POST("/set", middlewares.HasPermission("craft:route:schedule:set"), schedule.Set)   // 设置调度
	routers.DELETE("/remove/:ids", middlewares.HasPermission("craft:schedule:product:remove"),
		schedule.Remove) //移除调度
	routers.GET("/month/list", middlewares.HasPermission("craft:route:product:month:list"), schedule.MonthList) // 调度列表
}

// List 度列表
// @Summary 度列表
// @Description 度列表
// @Tags 工艺管理/调度管理
// @Param  object query craftRouteModels.SysProductScheduleReq true "组成工序列表参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/schedule/list [get]
func (schedule *Schedule) List(c *gin.Context) {
	req := new(craftRouteModels.SysProductScheduleListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := schedule.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
	return
}

// MonthList 月调度列表
// @Summary 月调度列表
// @Description 月调度列表
// @Tags 工艺管理/调度管理
// @Param  object query craftRouteModels.SysProductScheduleReq true "组成工序列表参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/schedule/month/list [get]
func (schedule *Schedule) MonthList(c *gin.Context) {
	req := new(craftRouteModels.SysProductScheduleReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	monthSchedule, err := schedule.service.GetMonthSchedule(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, monthSchedule)
}

// Set 月调度列表
// @Summary 月调度列表
// @Description 月调度列表
// @Tags 工艺管理/调度管理
// @Param  object query craftRouteModels.SetSysProductSchedule true "组成工序列表参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/schedule/set [post]
func (schedule *Schedule) Set(c *gin.Context) {
	data := new(craftRouteModels.SetSysProductSchedule)
	err := c.ShouldBindJSON(data)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	schedule.service.Set(c, data)
}

func (schedule *Schedule) Remove(c *gin.Context) {

}
