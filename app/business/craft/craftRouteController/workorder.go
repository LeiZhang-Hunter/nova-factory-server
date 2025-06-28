package craftRouteController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type WorkOrder struct {
	service craftRouteService.ISysProWorkorderService
}

func NewWorkOrder(service craftRouteService.ISysProWorkorderService) *WorkOrder {
	return &WorkOrder{
		service: service,
	}
}

func (w *WorkOrder) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/craft/route/work_order")
	routers.GET("/list", middlewares.HasPermission("craft:route:work_order"), w.List)                    // 生产工单
	routers.POST("/set", middlewares.HasPermission("craft:route:work_order:set"), w.Set)                 // 设置生产工单
	routers.DELETE("/remove/:ids", middlewares.HasPermission("craft:route:work_order:remove"), w.Remove) //移除生产工单
}

// List 生产工单列表
// @Summary 生产工单列表
// @Description 生产工单列表
// @Tags 工艺管理/生产工单管理
// @Param  object query craftRouteModels.SysProWorkorderReq true "生产工单列表参数"
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/product/bom/list [get]
func (w *WorkOrder) List(c *gin.Context) {
	req := new(craftRouteModels.SysProWorkorderReq)
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

// Set 设置生产工单
// @Summary 设置生产工单
// @Description 设置生产工单
// @Tags 工艺管理/生产工单管理
// @Param  object body craftRouteModels.SysSetProWorkorder true "设置生产工单参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/work_order/set [post]
func (w *WorkOrder) Set(c *gin.Context) {
	req := new(craftRouteModels.SysSetProWorkorder)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.WorkorderID == 0 {
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

// Remove 移除生产工单
// @Summary 移除生产工单
// @Description 移除生产工单
// @Tags 工艺管理/生产工单管理
// @Param  ids path string true "ids"
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/work_order/remove/{ids} [delete]
func (w *WorkOrder) Remove(c *gin.Context) {
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
