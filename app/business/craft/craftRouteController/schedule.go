package craftRouteController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/business/daemonize/daemonizeService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Schedule struct {
	service      craftRouteService.IScheduleService
	agentService daemonizeService.IotAgentService
}

func NewSchedule(service craftRouteService.IScheduleService, agentService daemonizeService.IotAgentService) *Schedule {
	return &Schedule{
		service:      service,
		agentService: agentService,
	}
}

func (schedule *Schedule) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/craft/route/schedule")
	routers.GET("/list", middlewares.HasPermission("craft:route:schedule:list"), schedule.List) // 调度列表
	routers.POST("/set", middlewares.HasPermission("craft:route:schedule:set"), schedule.Set)   // 设置调度
	routers.DELETE("/remove/:ids", middlewares.HasPermission("craft:schedule:product:remove"),
		schedule.Remove) //移除调度
	routers.GET("/month/list", middlewares.HasPermission("craft:route:product:month:list"), schedule.MonthList) // 调度列表
	routers.GET("/detail", middlewares.HasPermission("craft:route:schedule:detail"), schedule.Detail)           // 调度列表
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

// Set 设置调度日程
// @Summary 设置调度日程
// @Description 设置调度日程
// @Tags 工艺管理/调度管理
// @Param  object body craftRouteModels.SetSysProductSchedule true "设置调度日程参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置调度日程参数"
// @Router /craft/route/schedule/set [post]
func (schedule *Schedule) Set(c *gin.Context) {
	data := new(craftRouteModels.SetSysProductSchedule)
	err := c.ShouldBindJSON(data)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	if len(data.TimeManager) == 0 {
		baizeContext.Waring(c, "日程安排错误")
		return
	}

	// 检查type
	if data.Type != craftRouteModels.SPECIAL && data.Type != craftRouteModels.DAILY {
		baizeContext.Waring(c, "调度类型错误")
		return
	}

	// 检查网关是否存在
	info, err := schedule.agentService.GetByObjectId(c, uint64(data.GatewayID))
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if info == nil {
		baizeContext.Waring(c, "网关不存在")
		return
	}

	err = schedule.service.Set(c, data)
	if err != nil {
		zap.L().Error("schedule Set failed", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// Remove 删除调度日程
// @Summary 删除调度日程
// @Description 删除调度日程
// @Tags 工艺管理/调度管理
// @Param  ids path string true "ids"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置调度日程参数"
// @Router /craft/route/schedule/remove [delete]
func (schedule *Schedule) Remove(c *gin.Context) {
	recordIds := baizeContext.ParamStringArray(c, "ids")
	if len(recordIds) == 0 {
		baizeContext.Waring(c, "请选择调度任务")
		return
	}

	err := schedule.service.Remove(c, recordIds)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}

// Detail 调度详情
// @Summary 调度详情
// @Description 调度详情
// @Tags 工艺管理/调度管理
// @Param  object query craftRouteModels.DetailSysProductSchedule true "组成工序列表参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/schedule/detail [get]
func (schedule *Schedule) Detail(c *gin.Context) {
	req := new(craftRouteModels.DetailSysProductSchedule)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	schedule.service.Detail(c, req.Id)
	return
}
