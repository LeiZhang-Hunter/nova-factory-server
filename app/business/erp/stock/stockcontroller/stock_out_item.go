package stockcontroller

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// StockOutItem 其它出库单项控制器
type StockOutItem struct {
	service stockservice.IStockOutItemService
}

// NewStockOutItem 创建其它出库单项控制器。
func NewStockOutItem(service stockservice.IStockOutItemService) *StockOutItem {
	return &StockOutItem{service: service}
}

// PrivateRoutes 注册其它出库单项私有路由。
func (o *StockOutItem) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/stock/out-item")
	group.GET("/list", middlewares.HasPermission("erp:stock:outItem:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:stock:outItem:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:stock:outItem:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:stock:outItem:remove"), o.Delete)
}

// List 查询其它出库单项列表。
// @Summary 查询其它出库单项列表
// @Description 按条件分页查询其它出库单项列表
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param object query stockmodels.StockOutItemQuery true "ERP 其它出库单项查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/out-item/list [get]
func (o *StockOutItem) List(c *gin.Context) {
	req := new(stockmodels.StockOutItemQuery)
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

// GetByID 查询其它出库单项详情。
// @Summary 查询其它出库单项详情
// @Description 根据ID查询其它出库单项详情
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/out-item/query/{id} [get]
func (o *StockOutItem) GetByID(c *gin.Context) {
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

// Set 新增或修改其它出库单项。
// @Summary 新增或修改其它出库单项
// @Description 新增或修改其它出库单项
// @Tags ERP/库存管理
// @Security BearerAuth
// @Accept application/json
// @Param body body stockmodels.StockOutItemUpsert true "ERP 其它出库单项参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/stock/out-item/set [post]
func (o *StockOutItem) Set(c *gin.Context) {
	req := new(stockmodels.StockOutItemUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *stockmodels.StockOutItem
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

// Delete 删除其它出库单项。
// @Summary 删除其它出库单项
// @Description 根据ID删除其它出库单项，多个ID用逗号分隔
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/stock/out-item/remove/{ids} [delete]
func (o *StockOutItem) Delete(c *gin.Context) {
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
