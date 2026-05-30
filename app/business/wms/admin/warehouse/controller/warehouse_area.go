package controller

import (
	"nova-factory-server/app/business/wms/admin/warehouse/models"
	"nova-factory-server/app/business/wms/admin/warehouse/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// WarehouseArea WMS 库区控制器。
type WarehouseArea struct {
	service service.IWarehouseAreaService
}

// NewWarehouseArea 创建 WMS 库区控制器。
func NewWarehouseArea(service service.IWarehouseAreaService) *WarehouseArea {
	return &WarehouseArea{service: service}
}

// PrivateRoutes 注册 WMS 库区私有路由。
func (w *WarehouseArea) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/wms/warehouse/area")
	group.GET("/list", middlewares.HasPermission("wms:warehouse:area:list"), w.List)
	group.GET("/info/:id", middlewares.HasPermission("wms:warehouse:area:info"), w.Info)
	group.POST("/set", middlewares.HasPermission("wms:warehouse:area:set"), w.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("wms:warehouse:area:remove"), w.Remove)
}

// List 查询库区列表。
// @Summary 查询库区列表
// @Description 按条件分页查询 WMS 库区列表
// @Tags WMS/仓储管理/库区
// @Security BearerAuth
// @Param object query models.WarehouseAreaQuery true "WMS 库区查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /wms/admin/warehouse/area/list [get]
func (w *WarehouseArea) List(c *gin.Context) {
	req := new(models.WarehouseAreaQuery)
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

// Info 查询库区详情。
// @Summary 查询库区详情
// @Description 根据ID查询 WMS 库区详情
// @Tags WMS/仓储管理/库区
// @Security BearerAuth
// @Param id path int true "库区ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /wms/admin/warehouse/area/info/{id} [get]
func (w *WarehouseArea) Info(c *gin.Context) {
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

// Set 保存库区。
// @Summary 保存库区
// @Description 新增或修改 WMS 库区
// @Tags WMS/仓储管理/库区
// @Security BearerAuth
// @Accept application/json
// @Param body body models.WarehouseAreaSet true "WMS 库区保存参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /wms/admin/warehouse/area/set [post]
func (w *WarehouseArea) Set(c *gin.Context) {
	req := new(models.WarehouseAreaSet)
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

// Remove 删除库区。
// @Summary 删除库区
// @Description 根据ID删除 WMS 库区，多个ID使用逗号分隔
// @Tags WMS/仓储管理/库区
// @Security BearerAuth
// @Param ids path string true "库区ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /wms/admin/warehouse/area/remove/{ids} [delete]
func (w *WarehouseArea) Remove(c *gin.Context) {
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
