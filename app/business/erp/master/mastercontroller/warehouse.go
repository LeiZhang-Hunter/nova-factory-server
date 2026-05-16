package mastercontroller

import (
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Warehouse ERP 仓库控制器
type Warehouse struct {
	service masterservice.IWarehouseService
}

// NewWarehouse 创建 ERP 仓库控制器。
func NewWarehouse(service masterservice.IWarehouseService) *Warehouse {
	return &Warehouse{service: service}
}

// PrivateRoutes 注册 ERP 仓库私有路由。
func (o *Warehouse) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/master/warehouse")
	group.GET("/list", middlewares.HasPermission("erp:master:warehouse:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:master:warehouse:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:master:warehouse:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:master:warehouse:remove"), o.Delete)
}

// List 查询 ERP 仓库列表。
// @Summary 查询 ERP 仓库列表
// @Description 按条件分页查询 ERP 仓库列表
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param object query mastermodels.WarehouseQuery true "ERP 仓库查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/warehouse/list [get]
func (o *Warehouse) List(c *gin.Context) {
	req := new(mastermodels.WarehouseQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 查询 ERP 仓库详情。
// @Summary 查询 ERP 仓库详情
// @Description 根据ID查询 ERP 仓库详情
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/master/warehouse/query/{id} [get]
func (o *Warehouse) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 新增或修改 ERP 仓库。
// @Summary 新增或修改 ERP 仓库
// @Description 新增或修改 ERP 仓库
// @Tags ERP/基础资料
// @Security BearerAuth
// @Accept application/json
// @Param body body mastermodels.WarehouseUpsert true "ERP 仓库参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/master/warehouse/set [post]
func (o *Warehouse) Set(c *gin.Context) {
	req := new(mastermodels.WarehouseUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *mastermodels.Warehouse
		err  error
	)
	if req.ID > 0 {
		data, err = o.service.Update(c, req)
	} else {
		data, err = o.service.Create(c, req)
	}
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete 删除 ERP 仓库。
// @Summary 删除 ERP 仓库
// @Description 根据ID删除 ERP 仓库，多个ID用逗号分隔
// @Tags ERP/基础资料
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/master/warehouse/remove/{ids} [delete]
func (o *Warehouse) Delete(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := o.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
