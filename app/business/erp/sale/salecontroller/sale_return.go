package salecontroller

import (
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// SaleReturn 销售退货控制器
type SaleReturn struct {
	service saleservice.ISaleReturnService
}

// NewSaleReturn 创建销售退货控制器。
func NewSaleReturn(service saleservice.ISaleReturnService) *SaleReturn {
	return &SaleReturn{service: service}
}

// PrivateRoutes 注册销售退货私有路由。
func (o *SaleReturn) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/sale/return")
	group.GET("/list", middlewares.HasPermission("erp:sale:return:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:sale:return:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:sale:return:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:sale:return:remove"), o.Delete)
}

func (o *SaleReturn) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/erp/sale/return/list", "erp:sale:return:list")
	router.RegisterPermission("GET", "/erp/sale/return/query/:id", "erp:sale:return:query")
	router.RegisterPermission("POST", "/erp/sale/return/set", "erp:sale:return:set")
	router.RegisterPermission("DELETE", "/erp/sale/return/remove/:ids", "erp:sale:return:remove")
}

// List 查询销售退货列表。
// @Summary 查询销售退货列表
// @Description 按条件分页查询销售退货列表
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param object query salemodels.SaleReturnQuery true "ERP 销售退货查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/sale/return/list [get]
func (o *SaleReturn) List(c *gin.Context) {
	req := new(salemodels.SaleReturnQuery)
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

// GetByID 查询销售退货详情。
// @Summary 查询销售退货详情
// @Description 根据ID查询销售退货详情
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/sale/return/query/{id} [get]
func (o *SaleReturn) GetByID(c *gin.Context) {
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

// Set 新增或修改销售退货。
// @Summary 新增或修改销售退货
// @Description 新增或修改销售退货
// @Tags ERP/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body salemodels.SaleReturnUpsert true "ERP 销售退货参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/sale/return/set [post]
func (o *SaleReturn) Set(c *gin.Context) {
	req := new(salemodels.SaleReturnUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *salemodels.SaleReturn
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

// Delete 删除销售退货。
// @Summary 删除销售退货
// @Description 根据ID删除销售退货，多个ID用逗号分隔
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/sale/return/remove/{ids} [delete]
func (o *SaleReturn) Delete(c *gin.Context) {
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
