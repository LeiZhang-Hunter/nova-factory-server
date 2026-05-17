package financecontroller

import (
	"nova-factory-server/app/business/erp/finance/financemodels"
	"nova-factory-server/app/business/erp/finance/financeservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// FinanceReceipt ERP 收款单控制器
type FinanceReceipt struct {
	service financeservice.IFinanceReceiptService
}

// NewFinanceReceipt 创建 ERP 收款单控制器。
func NewFinanceReceipt(service financeservice.IFinanceReceiptService) *FinanceReceipt {
	return &FinanceReceipt{service: service}
}

// PrivateRoutes 注册 ERP 收款单私有路由。
func (o *FinanceReceipt) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/finance/receipt")
	group.GET("/list", middlewares.HasPermission("erp:finance:receipt:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:finance:receipt:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:finance:receipt:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:finance:receipt:remove"), o.Delete)
}

// List 查询 ERP 收款单列表。
// @Summary 查询 ERP 收款单列表
// @Description 按条件分页查询 ERP 收款单列表
// @Tags ERP/财务管理
// @Security BearerAuth
// @Param object query financemodels.FinanceReceiptQuery true "ERP 收款单查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/finance/receipt/list [get]
func (o *FinanceReceipt) List(c *gin.Context) {
	req := new(financemodels.FinanceReceiptQuery)
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

// GetByID 查询 ERP 收款单详情。
// @Summary 查询 ERP 收款单详情
// @Description 根据ID查询 ERP 收款单详情
// @Tags ERP/财务管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/finance/receipt/query/{id} [get]
func (o *FinanceReceipt) GetByID(c *gin.Context) {
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

// Set 新增或修改 ERP 收款单。
// @Summary 新增或修改 ERP 收款单
// @Description 新增或修改 ERP 收款单
// @Tags ERP/财务管理
// @Security BearerAuth
// @Accept application/json
// @Param body body financemodels.FinanceReceiptUpsert true "ERP 收款单参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/finance/receipt/set [post]
func (o *FinanceReceipt) Set(c *gin.Context) {
	req := new(financemodels.FinanceReceiptUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *financemodels.FinanceReceipt
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

// Delete 删除 ERP 收款单。
// @Summary 删除 ERP 收款单
// @Description 根据ID删除 ERP 收款单，多个ID用逗号分隔
// @Tags ERP/财务管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/finance/receipt/remove/{ids} [delete]
func (o *FinanceReceipt) Delete(c *gin.Context) {
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
