package stockcontroller

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// StockIn 其它入库单控制器
type StockIn struct {
	service stockservice.IStockInService
}

// NewStockIn 创建其它入库单控制器。
func NewStockIn(service stockservice.IStockInService) *StockIn {
	return &StockIn{service: service}
}

// PrivateRoutes 注册其它入库单私有路由。
func (o *StockIn) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/stock/in")
	group.GET("/list", middlewares.HasPermission("erp:stock:in:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:stock:in:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:stock:in:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:stock:in:remove"), o.Delete)
}

func (o *StockIn) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/erp/stock/in/list", "erp:stock:in:list")
	router.RegisterPermission("GET", "/erp/stock/in/query/:id", "erp:stock:in:query")
	router.RegisterPermission("POST", "/erp/stock/in/set", "erp:stock:in:set")
	router.RegisterPermission("DELETE", "/erp/stock/in/remove/:ids", "erp:stock:in:remove")
}

// List 查询其它入库单列表。
// @Summary 查询其它入库单列表
// @Description 按条件分页查询其它入库单列表
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param object query stockmodels.StockInQuery true "ERP 其它入库单查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/in/list [get]
func (o *StockIn) List(c *gin.Context) {
	req := new(stockmodels.StockInQuery)
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

// GetByID 查询其它入库单详情。
// @Summary 查询其它入库单详情
// @Description 根据ID查询其它入库单详情
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/in/query/{id} [get]
func (o *StockIn) GetByID(c *gin.Context) {
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

// Set 新增或修改其它入库单。
// @Summary 新增或修改其它入库单
// @Description 新增或修改其它入库单
// @Tags ERP/库存管理
// @Security BearerAuth
// @Accept application/json
// @Param body body stockmodels.StockInUpsert true "ERP 其它入库单参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/stock/in/set [post]
func (o *StockIn) Set(c *gin.Context) {
	req := new(stockmodels.StockInUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *stockmodels.StockIn
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

// Delete 删除其它入库单。
// @Summary 删除其它入库单
// @Description 根据ID删除其它入库单，多个ID用逗号分隔
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/stock/in/remove/{ids} [delete]
func (o *StockIn) Delete(c *gin.Context) {
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
