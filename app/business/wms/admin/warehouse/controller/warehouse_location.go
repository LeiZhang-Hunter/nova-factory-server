package controller

import (
	"nova-factory-server/app/business/wms/admin/warehouse/models"
	"nova-factory-server/app/business/wms/admin/warehouse/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// WarehouseLocation WMS 库位控制器。
type WarehouseLocation struct {
	service service.IWarehouseLocationService
}

// NewWarehouseLocation 创建 WMS 库位控制器。
func NewWarehouseLocation(service service.IWarehouseLocationService) *WarehouseLocation {
	return &WarehouseLocation{service: service}
}

// PrivateRoutes 注册 WMS 库位私有路由。
func (w *WarehouseLocation) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/wms/admin/warehouse/location")
	group.GET("/list", middlewares.HasPermission("wms:warehouse:location:list"), w.List)
	group.GET("/info/:id", middlewares.HasPermission("wms:warehouse:location:info"), w.Info)
	group.POST("/set", middlewares.HasPermission("wms:warehouse:location:set"), w.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("wms:warehouse:location:remove"), w.Remove)
}

// List 查询库位列表。
// @Summary 查询库位列表
// @Description 按条件分页查询 WMS 库位列表
// @Tags WMS/仓储管理/库位
// @Security BearerAuth
// @Param object query models.WarehouseLocationQuery true "WMS 库位查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /wms/admin/warehouse/location/list [get]
func (w *WarehouseLocation) List(c *gin.Context) {
	req := new(models.WarehouseLocationQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := w.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Info 查询库位详情。
// @Summary 查询库位详情
// @Description 根据ID查询 WMS 库位详情
// @Tags WMS/仓储管理/库位
// @Security BearerAuth
// @Param id path int true "库位ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /wms/admin/warehouse/location/info/{id} [get]
func (w *WarehouseLocation) Info(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := w.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 保存库位。
// @Summary 保存库位
// @Description 新增或修改 WMS 库位
// @Tags WMS/仓储管理/库位
// @Security BearerAuth
// @Accept application/json
// @Param body body models.WarehouseLocationSet true "WMS 库位保存参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /wms/admin/warehouse/location/set [post]
func (w *WarehouseLocation) Set(c *gin.Context) {
	req := new(models.WarehouseLocationSet)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := w.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Remove 删除库位。
// @Summary 删除库位
// @Description 根据ID删除 WMS 库位，多个ID使用逗号分隔
// @Tags WMS/仓储管理/库位
// @Security BearerAuth
// @Param ids path string true "库位ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /wms/admin/warehouse/location/remove/{ids} [delete]
func (w *WarehouseLocation) Remove(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := w.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
