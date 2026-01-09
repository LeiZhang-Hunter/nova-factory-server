package buildingController

import (
	"nova-factory-server/app/business/asset/building/buildingModels"
	"nova-factory-server/app/business/asset/building/buildingService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Floor struct {
	service buildingService.FloorService
}

func NewFloor(service buildingService.FloorService) *Floor {
	return &Floor{
		service: service,
	}
}

func (b *Floor) PrivateRoutes(router *gin.RouterGroup) {
	building := router.Group("/asset/floor")
	building.GET("/list", middlewares.HasPermission("asset:floor:list"), b.List)
	building.POST("/set", middlewares.HasPermission("asset:floor:set"), b.Set)
	building.DELETE("/remove/:ids", middlewares.HasPermission("asset:floor:remove"), b.Remove)
}

// Set 保存楼层
// @Summary 保存楼层
// @Description 保存楼层
// @Tags 资产管理/楼层管理
// @Param  object body buildingModels.SetSysFloor true "楼层参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/floor/set [post]
func (b *Floor) Set(c *gin.Context) {
	info := new(buildingModels.SetSysFloor)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	count, err := b.service.CheckUniqueFloor(c, info.ID, info.BuildingID, info.Level)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	if count > 0 {
		baizeContext.Waring(c, "楼层已经存在，不能重复添加")
		return
	}

	value, err := b.service.Set(c, info)
	baizeContext.SuccessData(c, value)

}

// Remove 删除楼层
// @Summary 删除楼层
// @Description 删除楼层
// @Tags 资产管理/楼层管理
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /asset/floor/remove/{ids}  [delete]
func (b *Floor) Remove(c *gin.Context) {
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

// List 读取楼层列表
// @Summary 读取楼层列表
// @Description 读取楼层列表
// @Tags 资产管理/楼层管理
// @Param  object query buildingModels.SetSysFloorListReq true "楼层列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /asset/floor/list [get]
func (b *Floor) List(c *gin.Context) {
	req := new(buildingModels.SetSysFloorListReq)
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
