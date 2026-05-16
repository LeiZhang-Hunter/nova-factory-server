package stockcontroller

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Stock ERP 产品库存控制器
type Stock struct {
	service stockservice.IStockService
}

// NewStock 创建 ERP 产品库存控制器。
func NewStock(service stockservice.IStockService) *Stock {
	return &Stock{service: service}
}

// PrivateRoutes 注册 ERP 产品库存私有路由。
func (o *Stock) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/stock/stock")
	group.GET("/list", middlewares.HasPermission("erp:stock:stock:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:stock:stock:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:stock:stock:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:stock:stock:remove"), o.Delete)
}

// List 查询 ERP 产品库存列表。
// @Summary 查询 ERP 产品库存列表
// @Description 按条件分页查询 ERP 产品库存列表
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param object query stockmodels.StockQuery true "ERP 产品库存查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/stock/list [get]
func (o *Stock) List(c *gin.Context) {
	req := new(stockmodels.StockQuery)
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

// GetByID 查询 ERP 产品库存详情。
// @Summary 查询 ERP 产品库存详情
// @Description 根据ID查询 ERP 产品库存详情
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/stock/query/{id} [get]
func (o *Stock) GetByID(c *gin.Context) {
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

// Set 新增或修改 ERP 产品库存。
// @Summary 新增或修改 ERP 产品库存
// @Description 新增或修改 ERP 产品库存
// @Tags ERP/库存管理
// @Security BearerAuth
// @Accept application/json
// @Param body body stockmodels.StockUpsert true "ERP 产品库存参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/stock/stock/set [post]
func (o *Stock) Set(c *gin.Context) {
	req := new(stockmodels.StockUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *stockmodels.Stock
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

// Delete 删除 ERP 产品库存。
// @Summary 删除 ERP 产品库存
// @Description 根据ID删除 ERP 产品库存，多个ID用逗号分隔
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/stock/stock/remove/{ids} [delete]
func (o *Stock) Delete(c *gin.Context) {
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
