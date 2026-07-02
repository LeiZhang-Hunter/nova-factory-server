package salecontroller

import (
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// SaleOutItem 销售出库项控制器
type SaleOutItem struct {
	service saleservice.ISaleOutItemService
}

// NewSaleOutItem 创建销售出库项控制器。
func NewSaleOutItem(service saleservice.ISaleOutItemService) *SaleOutItem {
	return &SaleOutItem{service: service}
}

// PrivateRoutes 注册销售出库项私有路由。
func (o *SaleOutItem) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/sale/out-items")
	group.GET("/list", middlewares.HasPermission("erp:sale:outItems:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:sale:outItems:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:sale:outItems:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:sale:outItems:remove"), o.Delete)
}

func (o *SaleOutItem) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/erp/sale/out-items/list", "erp:sale:outItems:list")
	router.RegisterPermission("GET", "/erp/sale/out-items/query/:id", "erp:sale:outItems:query")
	router.RegisterPermission("POST", "/erp/sale/out-items/set", "erp:sale:outItems:set")
	router.RegisterPermission("DELETE", "/erp/sale/out-items/remove/:ids", "erp:sale:outItems:remove")
}

// List 查询销售出库项列表。
// @Summary 查询销售出库项列表
// @Description 按条件分页查询销售出库项列表
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param object query salemodels.SaleOutItemQuery true "ERP 销售出库项查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/sale/out-items/list [get]
func (o *SaleOutItem) List(c *gin.Context) {
	req := new(salemodels.SaleOutItemQuery)
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

// GetByID 查询销售出库项详情。
// @Summary 查询销售出库项详情
// @Description 根据ID查询销售出库项详情
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/sale/out-items/query/{id} [get]
func (o *SaleOutItem) GetByID(c *gin.Context) {
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

// Set 新增或修改销售出库项。
// @Summary 新增或修改销售出库项
// @Description 新增或修改销售出库项
// @Tags ERP/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body salemodels.SaleOutItemUpsert true "ERP 销售出库项参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/sale/out-items/set [post]
func (o *SaleOutItem) Set(c *gin.Context) {
	req := new(salemodels.SaleOutItemUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *salemodels.SaleOutItem
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

// Delete 删除销售出库项。
// @Summary 删除销售出库项
// @Description 根据ID删除销售出库项，多个ID用逗号分隔
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/sale/out-items/remove/{ids} [delete]
func (o *SaleOutItem) Delete(c *gin.Context) {
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
