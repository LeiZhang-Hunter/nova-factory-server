package buildingController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/building/buildingModels"
	"nova-factory-server/app/business/asset/building/buildingService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Building struct {
	service buildingService.BuildingService
}

func NewBuilding(service buildingService.BuildingService) *Building {
	return &Building{
		service: service,
	}
}

func (b *Building) PrivateRoutes(router *gin.RouterGroup) {
	building := router.Group("/asset/building")
	building.GET("/list", middlewares.HasPermission("asset:building:list"), b.List)               // 设备列表
	building.POST("/set", middlewares.HasPermission("asset:building:set"), b.Set)                 // 设置设备信息
	building.DELETE("/remove/:ids", middlewares.HasPermission("asset:building:remove"), b.Remove) //删除设备分组列表
}

func (b *Building) PublicRoutes(router *gin.RouterGroup) {
	apiV1 := router.Group("/api/v1")
	apiV1.GET("/building/list", b.BuildDetailList)
}

// Set 保存建筑物
// @Summary 保存建筑物
// @Description 保存建筑物
// @Tags 资产管理/建筑物管理
// @Param  object body buildingModels.SetSysBuilding true "建筑物参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/building/set [post]
func (b *Building) Set(c *gin.Context) {
	info := new(buildingModels.SetSysBuilding)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	value, err := b.service.Set(c, info)
	baizeContext.SuccessData(c, value)

}

// Remove 删除建筑物
// @Summary 删除建筑物
// @Description 删除建筑物
// @Tags 资产管理/建筑物管理
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/building/remove/{ids}  [delete]
func (b *Building) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := b.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}

// List 读取建筑物列表
// @Summary 读取建筑物列表
// @Description 读取建筑物列表
// @Tags 资产管理/建筑物管理
// @Param  object query buildingModels.SetSysBuildingListReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/building/list [get]
func (b *Building) List(c *gin.Context) {
	req := new(buildingModels.SetSysBuildingListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := b.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// BuildDetailList 读取建筑物以及楼层详情列表
// @Summary 读取建筑物以及楼层详情列表
// @Description 读取建筑物以及楼层详情列表
// @Tags 资产管理/建筑物管理
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /api/v1/building/list [get]
func (b *Building) BuildDetailList(c *gin.Context) {
	list, err := b.service.AllDetail(c)
	if err != nil {
		zap.L().Error("get build list error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}
