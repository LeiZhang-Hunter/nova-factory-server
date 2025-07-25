package craftRouteController

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	routers.GET("/list", middlewares.HasPermission("craft:route"), craft.GetRouteList)                            // 工艺路线列表
	routers.POST("/set", middlewares.HasPermission("craft:route:set"), craft.SetRoute)                            // 设置工艺路线
	routers.DELETE("/remove/:craft_route_id", middlewares.HasPermission("craft:route:remove"), craft.RemoveRoute) //移除工艺路线
	routers.GET("/detail", middlewares.HasPermission("craft:route:detail"), craft.Detail)                         // 工艺路线详情
	routers.POST("/config/save", middlewares.HasPermission("craft:route:config:save"), craft.Save)                // 工艺路线详情
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
// @Router /craft/route/remove/{craft_route_id} [delete]
func (craft *Craft) RemoveRoute(c *gin.Context) {
	craftRouteId := baizeContext.ParamInt64Array(c, "craft_route_id")
	if len(craftRouteId) == 0 {
		baizeContext.Waring(c, "请输入工艺id")
		return
	}

	err := craft.craftService.RemoveCraftRoute(c, craftRouteId)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}

// Detail 工艺详情
// @Summary 工艺详情
// @Description 工艺详情
// @Tags 工艺管理
// @Param  object query craftRouteModels.SysCraftRouteDetailRequest true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/detail [get]
func (craft *Craft) Detail(c *gin.Context) {
	req := new(craftRouteModels.SysCraftRouteDetailRequest)
	err := c.ShouldBindQuery(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	route, err := craft.craftService.DetailCraftRoute(c, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			baizeContext.SuccessData(c, &craftRouteModels.SysCraftRouteConfig{
				RouteID: req.RouteID,
			})
			return
		}
		zap.L().Error("读取工艺失败", zap.Error(err))
		baizeContext.Waring(c, "读取工艺失败")
		return
	}
	baizeContext.SuccessData(c, route)
}

// Save 保存工艺制图
// @Summary 保存工艺制图
// @Description 保存工艺制图
// @Tags 工艺管理
// @Param  object body craftRouteModels.ProcessTopo true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /craft/route/config/save [post]
func (craft *Craft) Save(c *gin.Context) {
	req := new(craftRouteModels.ProcessTopo)
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	route, err := craft.craftService.SaveCraftRoute(c, req)
	if err != nil {
		zap.L().Error("保存工艺制图失败", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, route)
}
