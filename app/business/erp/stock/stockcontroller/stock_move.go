package stockcontroller

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// StockMove ERP 库存调拨单控制器
type StockMove struct {
	service stockservice.IStockMoveService
}

// NewStockMove 创建 ERP 库存调拨单控制器。
func NewStockMove(service stockservice.IStockMoveService) *StockMove {
	return &StockMove{service: service}
}

// PrivateRoutes 注册 ERP 库存调拨单私有路由。
func (o *StockMove) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/stock/move")
	group.GET("/list", middlewares.HasPermission("erp:stock:move:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:stock:move:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:stock:move:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:stock:move:remove"), o.Delete)
}

// List 查询 ERP 库存调拨单列表。
// @Summary 查询 ERP 库存调拨单列表
// @Description 按条件分页查询 ERP 库存调拨单列表
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param object query stockmodels.StockMoveQuery true "ERP 库存调拨单查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/move/list [get]
func (o *StockMove) List(c *gin.Context) {
	req := new(stockmodels.StockMoveQuery)
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

// GetByID 查询 ERP 库存调拨单详情。
// @Summary 查询 ERP 库存调拨单详情
// @Description 根据ID查询 ERP 库存调拨单详情
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/move/query/{id} [get]
func (o *StockMove) GetByID(c *gin.Context) {
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

// Set 新增或修改 ERP 库存调拨单。
// @Summary 新增或修改 ERP 库存调拨单
// @Description 新增或修改 ERP 库存调拨单
// @Tags ERP/库存管理
// @Security BearerAuth
// @Accept application/json
// @Param body body stockmodels.StockMoveUpsert true "ERP 库存调拨单参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/stock/move/set [post]
func (o *StockMove) Set(c *gin.Context) {
	req := new(stockmodels.StockMoveUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *stockmodels.StockMove
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

// Delete 删除 ERP 库存调拨单。
// @Summary 删除 ERP 库存调拨单
// @Description 根据ID删除 ERP 库存调拨单，多个ID用逗号分隔
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/stock/move/remove/{ids} [delete]
func (o *StockMove) Delete(c *gin.Context) {
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
