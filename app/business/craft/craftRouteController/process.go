package craftRouteController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Process struct {
	processService craftRouteService.ICraftProcessService
}

func NewProcess(processService craftRouteService.ICraftProcessService) *Process {
	return &Process{
		processService: processService,
	}
}

func (p *Process) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/craft/process")
	routers.GET("/list", middlewares.HasPermission("craft:process"), p.GetProcessList)                         // 工序列表
	routers.POST("/set", middlewares.HasPermission("craft:process:set"), p.SetProcess)                         // 设置工序
	routers.DELETE("/remove/:process_ids", middlewares.HasPermission("craft:process:remove"), p.RemoveProcess) //移除工序
}

// GetProcessList 工序列表
// @Summary 工序列表
// @Description 工序列表
// @Tags 工艺管理/工序管理
// @Param  object query craftRouteModels.SysProProcessListReq true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/process/list [get]
func (p *Process) GetProcessList(c *gin.Context) {
	req := new(craftRouteModels.SysProProcessListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := p.processService.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// SetProcess 设置工序
// @Summary 设置工序
// @Description 设置工序
// @Tags 工艺管理/工序管理
// @Param  object body craftRouteModels.SysProProcess true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/process/set [post]
func (p *Process) SetProcess(c *gin.Context) {
	req := new(craftRouteModels.SysProProcess)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.ProcessID == 0 {
		process, err := p.processService.Add(c, req)
		if err != nil {
			baizeContext.Waring(c, "添加供需失败")
			return
		}
		baizeContext.SuccessData(c, process)
	} else {
		process, err := p.processService.Update(c, req)
		if err != nil {
			baizeContext.Waring(c, "添加供需失败")
			return
		}
		baizeContext.SuccessData(c, process)
	}
}

// RemoveProcess 删除工序
// @Summary 删除工序
// @Description 删除工序
// @Tags 工艺管理/工序管理
// @Param  deviceIds path string true "deviceIds"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/process/remove/{process_ids} [delete]
func (p *Process) RemoveProcess(c *gin.Context) {
	processIds := baizeContext.ParamInt64Array(c, "process_ids")
	if len(processIds) == 0 {
		baizeContext.Waring(c, "请选择供需id")
		return
	}
	err := p.processService.Remove(c, processIds)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}
