package stockcontroller

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// StockCheck ERP 库存盘点单控制器
type StockCheck struct {
	service stockservice.IStockCheckService
}

// NewStockCheck 创建 ERP 库存盘点单控制器。
func NewStockCheck(service stockservice.IStockCheckService) *StockCheck {
	return &StockCheck{service: service}
}

// PrivateRoutes 注册 ERP 库存盘点单私有路由。
func (o *StockCheck) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/stock/check")
	group.GET("/list", middlewares.HasPermission("erp:stock:check:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:stock:check:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:stock:check:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:stock:check:remove"), o.Delete)
}

// List 查询 ERP 库存盘点单列表。
// @Summary 查询 ERP 库存盘点单列表
// @Description 按条件分页查询 ERP 库存盘点单列表
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param object query stockmodels.StockCheckQuery true "ERP 库存盘点单查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/check/list [get]
func (o *StockCheck) List(c *gin.Context) {
	req := new(stockmodels.StockCheckQuery)
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

// GetByID 查询 ERP 库存盘点单详情。
// @Summary 查询 ERP 库存盘点单详情
// @Description 根据ID查询 ERP 库存盘点单详情
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/check/query/{id} [get]
func (o *StockCheck) GetByID(c *gin.Context) {
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

// Set 新增或修改 ERP 库存盘点单。
// @Summary 新增或修改 ERP 库存盘点单
// @Description 新增或修改 ERP 库存盘点单
// @Tags ERP/库存管理
// @Security BearerAuth
// @Accept application/json
// @Param body body stockmodels.StockCheckUpsert true "ERP 库存盘点单参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/stock/check/set [post]
func (o *StockCheck) Set(c *gin.Context) {
	req := new(stockmodels.StockCheckUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *stockmodels.StockCheck
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

// Delete 删除 ERP 库存盘点单。
// @Summary 删除 ERP 库存盘点单
// @Description 根据ID删除 ERP 库存盘点单，多个ID用逗号分隔
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/stock/check/remove/{ids} [delete]
func (o *StockCheck) Delete(c *gin.Context) {
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
