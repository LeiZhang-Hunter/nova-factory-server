package purchasecontroller

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/business/erp/purchase/purchaseservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// PurchaseInItem采购入库项控制器
type PurchaseInItem struct {
	service purchaseservice.IPurchaseInItemService
}

// NewPurchaseInItem 创建采购入库项控制器。
func NewPurchaseInItem(service purchaseservice.IPurchaseInItemService) *PurchaseInItem {
	return &PurchaseInItem{service: service}
}

// PrivateRoutes 注册采购入库项私有路由。
func (o *PurchaseInItem) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/purchase/in-items")
	group.GET("/list", middlewares.HasPermission("erp:purchase:inItems:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:purchase:inItems:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:purchase:inItems:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:purchase:inItems:remove"), o.Delete)
}

// List 查询采购入库项列表。
// @Summary 查询采购入库项列表
// @Description 按条件分页查询采购入库项列表
// @Tags ERP/采购管理
// @Security BearerAuth
// @Param object query purchasemodels.PurchaseInItemQuery true "ERP 采购入库项查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/purchase/in-items/list [get]
func (o *PurchaseInItem) List(c *gin.Context) {
	req := new(purchasemodels.PurchaseInItemQuery)
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

// GetByID 查询采购入库项详情。
// @Summary 查询采购入库项详情
// @Description 根据ID查询采购入库项详情
// @Tags ERP/采购管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/purchase/in-items/query/{id} [get]
func (o *PurchaseInItem) GetByID(c *gin.Context) {
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

// Set 新增或修改采购入库项。
// @Summary 新增或修改采购入库项
// @Description 新增或修改采购入库项
// @Tags ERP/采购管理
// @Security BearerAuth
// @Accept application/json
// @Param body body purchasemodels.PurchaseInItemUpsert true "ERP 采购入库项参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/purchase/in-items/set [post]
func (o *PurchaseInItem) Set(c *gin.Context) {
	req := new(purchasemodels.PurchaseInItemUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *purchasemodels.PurchaseInItem
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

// Delete 删除采购入库项。
// @Summary 删除采购入库项
// @Description 根据ID删除采购入库项，多个ID用逗号分隔
// @Tags ERP/采购管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/purchase/in-items/remove/{ids} [delete]
func (o *PurchaseInItem) Delete(c *gin.Context) {
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
