package purchasecontroller

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/business/erp/purchase/purchaseservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// PurchaseReturn 采购退货控制器
type PurchaseReturn struct {
	service purchaseservice.IPurchaseReturnService
}

// NewPurchaseReturn 创建采购退货控制器。
func NewPurchaseReturn(service purchaseservice.IPurchaseReturnService) *PurchaseReturn {
	return &PurchaseReturn{service: service}
}

// PrivateRoutes 注册采购退货私有路由。
func (o *PurchaseReturn) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/purchase/return")
	group.GET("/list", middlewares.HasPermission("erp:purchase:return:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:purchase:return:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:purchase:return:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:purchase:return:remove"), o.Delete)
}

// List 查询采购退货列表。
// @Summary 查询采购退货列表
// @Description 按条件分页查询采购退货列表
// @Tags ERP/采购管理
// @Security BearerAuth
// @Param object query purchasemodels.PurchaseReturnQuery true "ERP 采购退货查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/purchase/return/list [get]
func (o *PurchaseReturn) List(c *gin.Context) {
	req := new(purchasemodels.PurchaseReturnQuery)
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

// GetByID 查询采购退货详情。
// @Summary 查询采购退货详情
// @Description 根据ID查询采购退货详情
// @Tags ERP/采购管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/purchase/return/query/{id} [get]
func (o *PurchaseReturn) GetByID(c *gin.Context) {
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

// Set 新增或修改采购退货。
// @Summary 新增或修改采购退货
// @Description 新增或修改采购退货
// @Tags ERP/采购管理
// @Security BearerAuth
// @Accept application/json
// @Param body body purchasemodels.PurchaseReturnUpsert true "ERP 采购退货参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/purchase/return/set [post]
func (o *PurchaseReturn) Set(c *gin.Context) {
	req := new(purchasemodels.PurchaseReturnUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *purchasemodels.PurchaseReturn
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

// Delete 删除采购退货。
// @Summary 删除采购退货
// @Description 根据ID删除采购退货，多个ID用逗号分隔
// @Tags ERP/采购管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/purchase/return/remove/{ids} [delete]
func (o *PurchaseReturn) Delete(c *gin.Context) {
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
