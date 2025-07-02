package craftRouteController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type ProcessContext struct {
	processContextService craftRouteService.ICraftProcessContextService
}

func NewProcessContext(processContextService craftRouteService.ICraftProcessContextService) *ProcessContext {
	return &ProcessContext{
		processContextService: processContextService,
	}
}

func (p *ProcessContext) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/craft/process/context")
	routers.GET("/list", middlewares.HasPermission("craft:process:context"), p.GetProcessContextList)                             // 工序内容列表
	routers.POST("/set", middlewares.HasPermission("craft:process:context:set"), p.SetProcessContextList)                         // 设置工序内容
	routers.DELETE("/remove/:context_ids", middlewares.HasPermission("craft:process:context:remove"), p.RemoveProcessContextList) //移除工序内容
}

// GetProcessContextList 工序内容列表
// @Summary 工序内容列表
// @Description 工序内容列表
// @Tags 工艺管理/工序内容管理
// @Param  object query craftRouteModels.SysProProcessContextListReq true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/process/context/list [get]
func (p *ProcessContext) GetProcessContextList(c *gin.Context) {
	req := new(craftRouteModels.SysProProcessContextListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := p.processContextService.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// SetProcessContextList 设置工序内容
// @Summary 设置工序内容
// @Description 设置工序内容
// @Tags 工艺管理/工序内容管理
// @Param  object body craftRouteModels.SysProSetProcessContent true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/process/context/set [post]
func (p *ProcessContext) SetProcessContextList(c *gin.Context) {
	req := new(craftRouteModels.SysProSetProcessContent)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.ContentID == 0 {
		ret, err := p.processContextService.Add(c, req)
		if err != nil {
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, ret)
	} else {
		ret, err := p.processContextService.Update(c, req)
		if err != nil {
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, ret)
	}
}

// RemoveProcessContextList 移除工序内容
// @Summary 移除工序内容
// @Description 移除工序内容
// @Tags 工艺管理/工序内容管理
// @Param  context_ids path string true "context_ids"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/process/context/remove/{context_ids} [delete]
func (p *ProcessContext) RemoveProcessContextList(c *gin.Context) {
	contextIds := baizeContext.ParamStringArray(c, "context_ids")
	if len(contextIds) == 0 {
		baizeContext.Waring(c, "请选择供需id")
		return
	}
	err := p.processContextService.Remove(c, contextIds)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}
