package stockcontroller

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// StockInItem 其它入库单项控制器
type StockInItem struct {
	service stockservice.IStockInItemService
}

// NewStockInItem 创建其它入库单项控制器。
func NewStockInItem(service stockservice.IStockInItemService) *StockInItem {
	return &StockInItem{service: service}
}

// PrivateRoutes 注册其它入库单项私有路由。
func (o *StockInItem) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/stock/in-item")
	group.GET("/list", middlewares.HasPermission("erp:stock:inItem:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:stock:inItem:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:stock:inItem:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:stock:inItem:remove"), o.Delete)
}

func (o *StockInItem) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/erp/stock/in-item/list", "erp:stock:inItem:list")
	router.RegisterPermission("GET", "/erp/stock/in-item/query/:id", "erp:stock:inItem:query")
	router.RegisterPermission("POST", "/erp/stock/in-item/set", "erp:stock:inItem:set")
	router.RegisterPermission("DELETE", "/erp/stock/in-item/remove/:ids", "erp:stock:inItem:remove")
}

// List 查询其它入库单项列表。
// @Summary 查询其它入库单项列表
// @Description 按条件分页查询其它入库单项列表
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param object query stockmodels.StockInItemQuery true "ERP 其它入库单项查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/in-item/list [get]
func (o *StockInItem) List(c *gin.Context) {
	req := new(stockmodels.StockInItemQuery)
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

// GetByID 查询其它入库单项详情。
// @Summary 查询其它入库单项详情
// @Description 根据ID查询其它入库单项详情
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/in-item/query/{id} [get]
func (o *StockInItem) GetByID(c *gin.Context) {
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

// Set 新增或修改其它入库单项。
// @Summary 新增或修改其它入库单项
// @Description 新增或修改其它入库单项
// @Tags ERP/库存管理
// @Security BearerAuth
// @Accept application/json
// @Param body body stockmodels.StockInItemUpsert true "ERP 其它入库单项参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/stock/in-item/set [post]
func (o *StockInItem) Set(c *gin.Context) {
	req := new(stockmodels.StockInItemUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *stockmodels.StockInItem
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

// Delete 删除其它入库单项。
// @Summary 删除其它入库单项
// @Description 根据ID删除其它入库单项，多个ID用逗号分隔
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/stock/in-item/remove/{ids} [delete]
func (o *StockInItem) Delete(c *gin.Context) {
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
