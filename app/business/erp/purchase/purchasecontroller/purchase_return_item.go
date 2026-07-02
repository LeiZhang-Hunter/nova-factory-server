package purchasecontroller

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/business/erp/purchase/purchaseservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// PurchaseReturnItem 采购退货项控制器
type PurchaseReturnItem struct {
	service purchaseservice.IPurchaseReturnItemService
}

// NewPurchaseReturnItem 创建采购退货项控制器。
func NewPurchaseReturnItem(service purchaseservice.IPurchaseReturnItemService) *PurchaseReturnItem {
	return &PurchaseReturnItem{service: service}
}

// PrivateRoutes 注册采购退货项私有路由。
func (o *PurchaseReturnItem) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/purchase/return-items")
	group.GET("/list", middlewares.HasPermission("erp:purchase:returnItems:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:purchase:returnItems:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:purchase:returnItems:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:purchase:returnItems:remove"), o.Delete)
}

func (o *PurchaseReturnItem) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/erp/purchase/return-items/list", "erp:purchase:returnItems:list")
	router.RegisterPermission("GET", "/erp/purchase/return-items/query/:id", "erp:purchase:returnItems:query")
	router.RegisterPermission("POST", "/erp/purchase/return-items/set", "erp:purchase:returnItems:set")
	router.RegisterPermission("DELETE", "/erp/purchase/return-items/remove/:ids", "erp:purchase:returnItems:remove")
}

// List 查询采购退货项列表。
// @Summary 查询采购退货项列表
// @Description 按条件分页查询采购退货项列表
// @Tags ERP/采购管理
// @Security BearerAuth
// @Param object query purchasemodels.PurchaseReturnItemQuery true "ERP 采购退货项查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/purchase/return-items/list [get]
func (o *PurchaseReturnItem) List(c *gin.Context) {
	req := new(purchasemodels.PurchaseReturnItemQuery)
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

// GetByID 查询采购退货项详情。
// @Summary 查询采购退货项详情
// @Description 根据ID查询采购退货项详情
// @Tags ERP/采购管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/purchase/return-items/query/{id} [get]
func (o *PurchaseReturnItem) GetByID(c *gin.Context) {
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

// Set 新增或修改采购退货项。
// @Summary 新增或修改采购退货项
// @Description 新增或修改采购退货项
// @Tags ERP/采购管理
// @Security BearerAuth
// @Accept application/json
// @Param body body purchasemodels.PurchaseReturnItemUpsert true "ERP 采购退货项参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/purchase/return-items/set [post]
func (o *PurchaseReturnItem) Set(c *gin.Context) {
	req := new(purchasemodels.PurchaseReturnItemUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *purchasemodels.PurchaseReturnItem
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

// Delete 删除采购退货项。
// @Summary 删除采购退货项
// @Description 根据ID删除采购退货项，多个ID用逗号分隔
// @Tags ERP/采购管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/purchase/return-items/remove/{ids} [delete]
func (o *PurchaseReturnItem) Delete(c *gin.Context) {
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
