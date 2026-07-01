package salecontroller

import (
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// SaleReturnItem 销售退货项控制器
type SaleReturnItem struct {
	service saleservice.ISaleReturnItemService
}

// NewSaleReturnItem 创建销售退货项控制器。
func NewSaleReturnItem(service saleservice.ISaleReturnItemService) *SaleReturnItem {
	return &SaleReturnItem{service: service}
}

// PrivateRoutes 注册销售退货项私有路由。
func (o *SaleReturnItem) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/sale/return-items")
	group.GET("/list", middlewares.HasPermission("erp:sale:returnItems:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:sale:returnItems:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:sale:returnItems:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:sale:returnItems:remove"), o.Delete)
}

func (o *SaleReturnItem) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/erp/sale/return-items/list", "erp:sale:returnItems:list")
	router.RegisterPermission("GET", "/erp/sale/return-items/query/:id", "erp:sale:returnItems:query")
	router.RegisterPermission("POST", "/erp/sale/return-items/set", "erp:sale:returnItems:set")
	router.RegisterPermission("DELETE", "/erp/sale/return-items/remove/:ids", "erp:sale:returnItems:remove")
}

// List 查询销售退货项列表。
// @Summary 查询销售退货项列表
// @Description 按条件分页查询销售退货项列表
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param object query salemodels.SaleReturnItemQuery true "ERP 销售退货项查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/sale/return-items/list [get]
func (o *SaleReturnItem) List(c *gin.Context) {
	req := new(salemodels.SaleReturnItemQuery)
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

// GetByID 查询销售退货项详情。
// @Summary 查询销售退货项详情
// @Description 根据ID查询销售退货项详情
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/sale/return-items/query/{id} [get]
func (o *SaleReturnItem) GetByID(c *gin.Context) {
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

// Set 新增或修改销售退货项。
// @Summary 新增或修改销售退货项
// @Description 新增或修改销售退货项
// @Tags ERP/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body salemodels.SaleReturnItemUpsert true "ERP 销售退货项参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/sale/return-items/set [post]
func (o *SaleReturnItem) Set(c *gin.Context) {
	req := new(salemodels.SaleReturnItemUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *salemodels.SaleReturnItem
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

// Delete 删除销售退货项。
// @Summary 删除销售退货项
// @Description 根据ID删除销售退货项，多个ID用逗号分隔
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/sale/return-items/remove/{ids} [delete]
func (o *SaleReturnItem) Delete(c *gin.Context) {
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
