package daemonizeController

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/daemonize/grpc/v1"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
	"nova-factory-server/app/business/daemonize/daemonizeService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type IotAgent struct {
	service   daemonizeService.IotAgentService
	daemonize daemonizeService.DaemonizeService
}

func NewIotAgentController(service daemonizeService.IotAgentService, daemonize daemonizeService.DaemonizeService) *IotAgent {
	return &IotAgent{
		service:   service,
		daemonize: daemonize,
	}
}

func (i *IotAgent) PrivateRoutes(router *gin.RouterGroup) {
	agent := router.Group("/daemonize/agent")
	agent.GET("/list", middlewares.HasPermission("daemonize:agent:list"), i.List)                     // agent列表
	agent.POST("/set", middlewares.HasPermission("daemonize:agent:set"), i.Set)                       // 设置agent
	agent.DELETE("/remove/:ids", middlewares.HasPermission("daemonize:agent:remove"), i.Remove)       //移除agent
	agent.POST("/process/start", middlewares.HasPermission("daemonize:agent:process:start"), i.Start) //启动进程
	agent.POST("/process/stop", middlewares.HasPermission("daemonize:agent:process:stop"), i.Stop)    //停止进程
	agent.GET("/info", middlewares.HasPermission("daemonize:agent:info"), i.Info)                     // agent列表
}

// List Agent列表
// @Summary Agent列表
// @Description Agent列表
// @Tags 网关管理/Agent管理
// @Param  object query daemonizeModels.SysIotAgentListReq true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /daemonize/agent/list [get]
func (i *IotAgent) List(c *gin.Context) {
	req := new(daemonizeModels.SysIotAgentListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := i.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Set 设置Agent
// @Summary 设置Agent
// @Description 设置Agent
// @Tags 网关管理/Agent管理
// @Param  object body daemonizeModels.SysIotAgentSetReq true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /daemonize/agent/set [post]
func (i *IotAgent) Set(c *gin.Context) {
	req := new(daemonizeModels.SysIotAgentSetReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.ObjectID == 0 {
		data, err := i.service.Add(c, req)
		if err != nil {
			baizeContext.Waring(c, "添加Agent失败")
			return
		}
		baizeContext.SuccessData(c, data)
	} else {
		data, err := i.service.Update(c, req)
		if err != nil {
			baizeContext.Waring(c, "添加Agent失败")
			return
		}
		baizeContext.SuccessData(c, data)
	}
}

// Remove 移除Agent
// @Summary 移除Agent
// @Description 移除Agent
// @Tags 网关管理/Agent管理
// @Param  ids path string true "ids"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /daemonize/agent/remove/{ids} [delete]
func (i *IotAgent) Remove(c *gin.Context) {
	contextIds := baizeContext.ParamStringArray(c, "ids")
	if len(contextIds) == 0 {
		baizeContext.Waring(c, "请选择供需id")
		return
	}
	err := i.service.Remove(c, contextIds)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}

// Stop 停止Agent进程
// @Summary 停止Agent进程
// @Description 停止Agent进程
// @Tags 网关管理/Agent管理
// @Param  object body daemonizeModels.StartProcessReq true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /daemonize/agent/process/stop [post]
func (i *IotAgent) Stop(c *gin.Context) {
	req := new(daemonizeModels.StartProcessReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if len(req.ProcessOperateInfoList) == 0 {
		baizeContext.Waring(c, "请选择选项")
		return
	}
	err = i.daemonize.BroadcastAgentOperateProcess(c, v1.AgentCmd_Stop, daemonizeModels.ToPbProcessList(req.ProcessOperateInfoList))
	if err != nil {
		baizeContext.Waring(c, "操作失败")
		return
	}
	baizeContext.Success(c)
}

// Start 启动Agent进程
// @Summary 启动Agent进程
// @Description 启动Agent进程
// @Tags 网关管理/Agent管理
// @Param  object body daemonizeModels.StartProcessReq true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /daemonize/agent/process/start [post]
func (i *IotAgent) Start(c *gin.Context) {
	req := new(daemonizeModels.StartProcessReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if len(req.ProcessOperateInfoList) == 0 {
		baizeContext.Waring(c, "请选择选项")
		return
	}
	err = i.daemonize.BroadcastAgentOperateProcess(c, v1.AgentCmd_Start, daemonizeModels.ToPbProcessList(req.ProcessOperateInfoList))
	if err != nil {
		baizeContext.Waring(c, "操作失败")
		return
	}
	baizeContext.Success(c)
}

// Info agent详情
// @Summary agent详情
// @Description agent详情
// @Tags 网关管理/Agent管理
// @Param  object query daemonizeModels.SysIotAgentQueryReq true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /daemonize/agent/info [get]
func (i *IotAgent) Info(c *gin.Context) {
	req := new(daemonizeModels.SysIotAgentQueryReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	info, err := i.service.Info(c, req.ObjectID)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, info)
}
