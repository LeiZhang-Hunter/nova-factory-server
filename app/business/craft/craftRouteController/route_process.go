package craftRouteController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type RouteProcess struct {
	service craftRouteService.IProcessRouteService
}

func NewSysProRouteProcess(service craftRouteService.IProcessRouteService) *RouteProcess {
	return &RouteProcess{
		service: service,
	}
}

func (p *RouteProcess) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/craft/route/process")
	routers.GET("/list", middlewares.HasPermission("craft:route:process"), p.GetRouteProcessList)                               // 工序内容列表
	routers.POST("/set", middlewares.HasPermission("craft:route:process:set"), p.SetRouteProcess)                               // 设置工序内容
	routers.DELETE("/remove/:route_process_ids", middlewares.HasPermission("craft:route:process:remove"), p.RemoveRouteProcess) //移除工序内容
}

// GetRouteProcessList 组成工序列表
// @Summary 组成工序列表
// @Description 组成工序列表
// @Tags 工艺管理/工序组成管理
// @Param  object body craftRouteModels.SysProRouteProcessListReq true "组成工序列表参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/process/list [get]
func (p *RouteProcess) GetRouteProcessList(c *gin.Context) {
	req := new(craftRouteModels.SysProRouteProcessListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := p.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// SetRouteProcess 设置组成工序
// @Summary 设置组成工序
// @Description 设置组成工序
// @Tags 工艺管理/工序组成管理
// @Param  object body craftRouteModels.SysProRouteProcessSetRequest true "设置组成工序参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/process/set [post]
func (p *RouteProcess) SetRouteProcess(c *gin.Context) {
	req := new(craftRouteModels.SysProRouteProcessSetRequest)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.RecordID == 0 {
		data, err := p.service.Add(c, req)
		if err != nil {
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, data)
	} else {
		data, err := p.service.Update(c, req)
		if err != nil {
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, data)
	}

}

// RemoveRouteProcess 删除组成工序
// @Summary 删除组成工序
// @Description 删除组成工序
// @Tags 工艺管理/工序组成管理
// @Param  route_process_ids path string true "route_process_ids"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/process/remove [delete]
func (p *RouteProcess) RemoveRouteProcess(c *gin.Context) {
	contextIds := baizeContext.ParamStringArray(c, "route_process_ids")
	if len(contextIds) == 0 {
		baizeContext.Waring(c, "请选择供需id")
		return
	}
	err := p.service.Remove(c, contextIds)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}
