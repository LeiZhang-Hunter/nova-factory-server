package craftRouteController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type RouteProductBom struct {
	service craftRouteService.ISysProRouteProductBomService
}

func NewRouteProductBom(service craftRouteService.ISysProRouteProductBomService) *RouteProductBom {
	return &RouteProductBom{
		service: service,
	}
}

func (p *RouteProductBom) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/craft/route/product/bom")
	routers.GET("/list", middlewares.HasPermission("craft:route:product:bom"), p.GetRouteProductBomList)                        // 工序列表
	routers.POST("/set", middlewares.HasPermission("craft:route:product:bom:set"), p.SetRouteProductBom)                        // 设置工序
	routers.DELETE("/remove/:record_ids", middlewares.HasPermission("craft:route:product:bom:remove"), p.RemoveRouteProductBom) //移除工序
}

// GetRouteProductBomList 产品制程物料BOM物料BOM列表
// @Summary 产品制程物料BOM物料BOM列表
// @Description 产品制程物料BOM物料BOM列表
// @Tags 工艺管理/产品制程物料BOM物料BOM管理
// @Param  object query craftRouteModels.SysProRouteProductBomReq true "组成工序列表参数"
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/product/bom/list [get]
func (p *RouteProductBom) GetRouteProductBomList(c *gin.Context) {
	req := new(craftRouteModels.SysProRouteProductBomReq)
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

// SetRouteProductBom 设置产品制程物料BOM物料BOM
// @Summary 设置产品制程物料BOM
// @Description 设置产品制程物料BOM
// @Tags 工艺管理/产品制程物料BOM管理
// @Param  object body craftRouteModels.SysSetProRouteProductBom true "设置产品制程物料BOM参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/product/bom/set [post]
func (p *RouteProductBom) SetRouteProductBom(c *gin.Context) {
	req := new(craftRouteModels.SysSetProRouteProductBom)
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

// RemoveRouteProductBom 删除产品制程物料BOM
// @Summary 删除产品制程物料BOM
// @Description 删除产品制程物料BOM
// @Tags 工艺管理/产品制程物料BOM管理
// @Param  record_ids path string true "record_ids"
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/product/bom/remove/{record_ids} [delete]
func (p *RouteProductBom) RemoveRouteProductBom(c *gin.Context) {
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
