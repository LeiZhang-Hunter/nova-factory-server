package craftRouteController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type RouteProduct struct {
	service craftRouteService.ISysProRouteProductService
}

func NewRouteProduct(service craftRouteService.ISysProRouteProductService) *RouteProduct {
	return &RouteProduct{
		service: service,
	}
}

func (p *RouteProduct) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/craft/route/product")
	routers.GET("/list", middlewares.HasPermission("craft:route:product"), p.GetRouteProductList)                        // 工序列表
	routers.POST("/set", middlewares.HasPermission("craft:route:product:set"), p.SetRouteProduct)                        // 设置工序
	routers.DELETE("/remove/:record_ids", middlewares.HasPermission("craft:route:product:remove"), p.RemoveRouteProduct) //移除工序
}

// GetRouteProductList 产品制程列表
// @Summary 产品制程列表
// @Description 产品制程列表
// @Tags 工艺管理/产品制程管理
// @Param  object query craftRouteModels.SysProRouteProductReq true "组成工序列表参数"
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/product/list [get]
func (p *RouteProduct) GetRouteProductList(c *gin.Context) {
	req := new(craftRouteModels.SysProRouteProductReq)
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

// SetRouteProduct 设置产品制程
// @Summary 设置产品制程
// @Description 设置产品制程
// @Tags 工艺管理/产品制程管理
// @Param  object body craftRouteModels.SysProRouteSetProduct true "设置产品制程参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/product/set [post]
func (p *RouteProduct) SetRouteProduct(c *gin.Context) {
	req := new(craftRouteModels.SysProRouteSetProduct)
	err := c.ShouldBindJSON(req)
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

// RemoveRouteProduct 删除产品制程
// @Summary 删除产品制程
// @Description 删除产品制程
// @Tags 工艺管理/产品制程管理
// @Param  record_ids path string true "record_ids"
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/product/remove/{record_ids} [delete]
func (p *RouteProduct) RemoveRouteProduct(c *gin.Context) {
	recordIds := baizeContext.ParamStringArray(c, "record_ids")
	if len(recordIds) == 0 {
		baizeContext.Waring(c, "请选择供需id")
		return
	}
	err := p.service.Remove(c, recordIds)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}
