package craftRouteController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Craft struct {
	craftService craftRouteService.ICraftRouteService
}

func NewCraft(craft craftRouteService.ICraftRouteService) *Craft {
	return &Craft{
		craftService: craft,
	}
}

func (craft *Craft) PrivateRoutes(router *gin.RouterGroup) {
	routers := router.Group("/craft/route")
	routers.GET("/list", middlewares.HasPermission("craft:route"), craft.GetRouteList)                        // 知识库列表
	routers.POST("/set", middlewares.HasPermission("craft:route:set"), craft.SetRoute)                        // 设置工艺路线
	routers.DELETE("/remove/:dataset_id", middlewares.HasPermission("craft:route:remove"), craft.RemoveRoute) //移除工艺路线
}

// GetRouteList 读取工艺列表
// @Summary 读取工艺列表
// @Description 读取工艺列表
// @Tags 工艺管理
// @Param  object query craftRouteModels.SysCraftRouteListReq true "设备分组列表请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /craft/route/list [get]
func (craft *Craft) GetRouteList(c *gin.Context) {
	req := new(craftRouteModels.SysCraftRouteListReq)
	err := c.ShouldBind(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.SuccessData(c, &craftRouteModels.SysCraftRouteListData{})
		return
	}
	crafts, err := craft.craftService.SelectCraftRoute(c, req)
	if err != nil {
		zap.L().Error("工艺列表错误", zap.Error(err))
		baizeContext.SuccessData(c, &aiDataSetModels.ChunkListResponse{})
		return
	}
	baizeContext.SuccessData(c, crafts)
}

// SetRoute 设置工艺路线
// @Summary 设置工艺路线
// @Description 设置工艺路线
// @Tags 工艺管理
// @Param  object body craftRouteModels.SysCraftRouteRequest true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/set [post]
func (craft *Craft) SetRoute(c *gin.Context) {
	req := new(craftRouteModels.SysCraftRouteRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.RouteID == 0 {
		set, err := craft.craftService.AddCraftRoute(c, req)
		if err != nil {
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, set)
	} else {
		set, err := craft.craftService.UpdateCraftRoute(c, req)
		if err != nil {
			baizeContext.Waring(c, err.Error())
			return
		}
		baizeContext.SuccessData(c, set)
	}

}

// RemoveRoute 移除工艺
// @Summary 移除工艺
// @Description 移除工艺
// @Tags 工艺管理
// @Param  craft_route_id path int64 true "craft_route_id"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/remove [delete]
func (craft *Craft) RemoveRoute(c *gin.Context) {
	craftRouteId := baizeContext.ParamInt64Array(c, "craft_route_id")
	if len(craftRouteId) == 0 {
		baizeContext.Waring(c, "请输入知识库id")
		return
	}

	err := craft.craftService.RemoveCraftRoute(c, craftRouteId)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}
