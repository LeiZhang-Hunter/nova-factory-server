package stockcontroller

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// StockRecord 产品库存明细控制器
type StockRecord struct {
	service stockservice.IStockRecordService
}

// NewStockRecord 创建产品库存明细控制器。
func NewStockRecord(service stockservice.IStockRecordService) *StockRecord {
	return &StockRecord{service: service}
}

// PrivateRoutes 注册产品库存明细私有路由。
func (o *StockRecord) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/stock/record")
	group.GET("/list", middlewares.HasPermission("erp:stock:record:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:stock:record:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:stock:record:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:stock:record:remove"), o.Delete)
}

func (o *StockRecord) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/erp/stock/record/list", "erp:stock:record:list")
	router.RegisterPermission("GET", "/erp/stock/record/query/:id", "erp:stock:record:query")
	router.RegisterPermission("POST", "/erp/stock/record/set", "erp:stock:record:set")
	router.RegisterPermission("DELETE", "/erp/stock/record/remove/:ids", "erp:stock:record:remove")
}

// List 查询产品库存明细列表。
// @Summary 查询产品库存明细列表
// @Description 按条件分页查询产品库存明细列表
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param object query stockmodels.StockRecordQuery true "ERP 产品库存明细查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/record/list [get]
func (o *StockRecord) List(c *gin.Context) {
	req := new(stockmodels.StockRecordQuery)
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

// GetByID 查询产品库存明细详情。
// @Summary 查询产品库存明细详情
// @Description 根据ID查询产品库存明细详情
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/stock/record/query/{id} [get]
func (o *StockRecord) GetByID(c *gin.Context) {
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

// Set 新增或修改产品库存明细。
// @Summary 新增或修改产品库存明细
// @Description 新增或修改产品库存明细
// @Tags ERP/库存管理
// @Security BearerAuth
// @Accept application/json
// @Param body body stockmodels.StockRecordUpsert true "ERP 产品库存明细参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/stock/record/set [post]
func (o *StockRecord) Set(c *gin.Context) {
	req := new(stockmodels.StockRecordUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *stockmodels.StockRecord
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

// Delete 删除产品库存明细。
// @Summary 删除产品库存明细
// @Description 根据ID删除产品库存明细，多个ID用逗号分隔
// @Tags ERP/库存管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/stock/record/remove/{ids} [delete]
func (o *StockRecord) Delete(c *gin.Context) {
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
