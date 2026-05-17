package purchasecontroller

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/business/erp/purchase/purchaseservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// PurchaseOrder ERP 采购订单控制器
type PurchaseOrder struct {
	service purchaseservice.IPurchaseOrderService
}

// NewPurchaseOrder 创建 ERP 采购订单控制器。
func NewPurchaseOrder(service purchaseservice.IPurchaseOrderService) *PurchaseOrder {
	return &PurchaseOrder{service: service}
}

// PrivateRoutes 注册 ERP 采购订单私有路由。
func (o *PurchaseOrder) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/purchase/order")
	group.GET("/list", middlewares.HasPermission("erp:purchase:order:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:purchase:order:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:purchase:order:set"), o.Set)
	group.POST("/update-status", middlewares.HasPermission("erp:purchase:order:set"), o.UpdateStatus)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:purchase:order:remove"), o.Delete)
}

// List 查询 ERP 采购订单列表。
// @Summary 查询 ERP 采购订单列表
// @Description 按条件分页查询 ERP 采购订单列表
// @Tags ERP/采购管理
// @Security BearerAuth
// @Param object query purchasemodels.PurchaseOrderQuery true "ERP 采购订单查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/purchase/order/list [get]
func (o *PurchaseOrder) List(c *gin.Context) {
	req := new(purchasemodels.PurchaseOrderQuery)
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

// GetByID 查询 ERP 采购订单详情。
// @Summary 查询 ERP 采购订单详情
// @Description 根据ID查询 ERP 采购订单详情
// @Tags ERP/采购管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/purchase/order/query/{id} [get]
func (o *PurchaseOrder) GetByID(c *gin.Context) {
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

// Set 新增或修改 ERP 采购订单。
// @Summary 新增或修改 ERP 采购订单
// @Description 新增或修改 ERP 采购订单
// @Tags ERP/采购管理
// @Security BearerAuth
// @Accept application/json
// @Param body body purchasemodels.PurchaseOrderUpsert true "ERP 采购订单及明细参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/purchase/order/set [post]
func (o *PurchaseOrder) Set(c *gin.Context) {
	req := new(purchasemodels.PurchaseOrderUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *purchasemodels.PurchaseOrder
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

// UpdateStatus 更新 ERP 采购订单状态。
// @Summary 更新 ERP 采购订单状态
// @Description 审核或反审核 ERP 采购订单
// @Tags ERP/采购管理
// @Security BearerAuth
// @Accept application/json
// @Param body body purchasemodels.PurchaseOrderStatusReq true "ERP 采购订单状态参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "更新成功"
// @Router /erp/purchase/order/update-status [post]
func (o *PurchaseOrder) UpdateStatus(c *gin.Context) {
	req := new(purchasemodels.PurchaseOrderStatusReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if err := o.service.UpdateStatus(c, req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// Delete 删除 ERP 采购订单。
// @Summary 删除 ERP 采购订单
// @Description 根据ID删除 ERP 采购订单，多个ID用逗号分隔
// @Tags ERP/采购管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/purchase/order/remove/{ids} [delete]
func (o *PurchaseOrder) Delete(c *gin.Context) {
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
