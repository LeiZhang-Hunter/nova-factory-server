package craftRouteController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/business/daemonize/daemonizeService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Task struct {
	service      craftRouteService.ISysProTaskService
	agentService daemonizeService.IotAgentService
}

func NewTask(service craftRouteService.ISysProTaskService, agentService daemonizeService.IotAgentService) *Task {
	return &Task{
		service:      service,
		agentService: agentService,
	}
}

func (w *Task) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/craft/route/task")
	routers.GET("/list", middlewares.HasPermission("craft:route:task:list"), w.List)               // 生产工单
	routers.POST("/set", middlewares.HasPermission("craft:route:task:set"), w.Set)                 // 设置生产工单
	routers.DELETE("/remove/:ids", middlewares.HasPermission("craft:route:task:remove"), w.Remove) //移除生产工单
}

func (w *Task) PublicRoutes(router *gin.RouterGroup) {
	routers := router.Group("/api/product/task/v1")
	routers.POST("/schedule", w.Schedule)
}

// Schedule 读取调度任务
// @Summary 读取调度任务
// @Description 读取调度任务
// @Tags 工艺管理/生产任务管理
// @Param  object body craftRouteModels.ScheduleReq true "读取调度任务参数"
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /api/product/task/v1/schedule [post]
func (w *Task) Schedule(c *gin.Context) {
	req := new(craftRouteModels.ScheduleReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	info, err := w.agentService.Info(c, uint64(req.GatewayId))
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if info == nil {
		baizeContext.Waring(c, "agent is not found")
		return
	}
	if info.Username != req.UserName {
		baizeContext.Waring(c, "username error")
		return
	}
	if info.Password != req.Password {
		baizeContext.Waring(c, "password error")
		return
	}
	list, err := w.service.Schedule(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// List 生产任务列表
// @Summary 生产任务列表
// @Description 生产任务列表
// @Tags 工艺管理/生产任务管理
// @Param  object query craftRouteModels.SysProTaskReq true "生生产任务列表参数"
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/task/list [get]
func (w *Task) List(c *gin.Context) {
	req := new(craftRouteModels.SysProTaskReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := w.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Set 设置生产任务
// @Summary 设置生产任务
// @Description 设置生产任务
// @Tags 工艺管理/生产任务管理
// @Param  object body craftRouteModels.SysSetProTask true "设置生生产任务"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/task/set [post]
func (w *Task) Set(c *gin.Context) {
	req := new(craftRouteModels.SysSetProTask)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.TaskID == 0 {
		ret, err := w.service.Add(c, req)
		if err != nil {
			return
		}
		baizeContext.SuccessData(c, ret)
	} else {
		ret, err := w.service.Update(c, req)
		if err != nil {
			return
		}
		baizeContext.SuccessData(c, ret)
	}
}

// Remove 移除生产任务
// @Summary 移除生产任务
// @Description 移除生产任务
// @Tags 工艺管理/生产工单管理
// @Param  ids path string true "ids"
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/task/remove/{ids} [delete]
func (w *Task) Remove(c *gin.Context) {
	recordIds := baizeContext.ParamStringArray(c, "ids")
	if len(recordIds) == 0 {
		baizeContext.Waring(c, "请选择供需id")
		return
	}
	err := w.service.Remove(c, recordIds)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}
