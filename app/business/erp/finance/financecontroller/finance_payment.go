package financecontroller

import (
	"nova-factory-server/app/business/erp/finance/financemodels"
	"nova-factory-server/app/business/erp/finance/financeservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// FinancePayment ERP 付款单控制器
type FinancePayment struct {
	service financeservice.IFinancePaymentService
}

// NewFinancePayment 创建 ERP 付款单控制器。
func NewFinancePayment(service financeservice.IFinancePaymentService) *FinancePayment {
	return &FinancePayment{service: service}
}

// PrivateRoutes 注册 ERP 付款单私有路由。
func (o *FinancePayment) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/finance/payment")
	group.GET("/list", middlewares.HasPermission("erp:finance:payment:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:finance:payment:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:finance:payment:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:finance:payment:remove"), o.Delete)
}

func (o *FinancePayment) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/erp/finance/payment/list", "erp:finance:payment:list")
	router.RegisterPermission("GET", "/erp/finance/payment/query/:id", "erp:finance:payment:query")
	router.RegisterPermission("POST", "/erp/finance/payment/set", "erp:finance:payment:set")
	router.RegisterPermission("DELETE", "/erp/finance/payment/remove/:ids", "erp:finance:payment:remove")
}

// List 查询 ERP 付款单列表。
// @Summary 查询 ERP 付款单列表
// @Description 按条件分页查询 ERP 付款单列表
// @Tags ERP/财务管理
// @Security BearerAuth
// @Param object query financemodels.FinancePaymentQuery true "ERP 付款单查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/finance/payment/list [get]
func (o *FinancePayment) List(c *gin.Context) {
	req := new(financemodels.FinancePaymentQuery)
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

// GetByID 查询 ERP 付款单详情。
// @Summary 查询 ERP 付款单详情
// @Description 根据ID查询 ERP 付款单详情
// @Tags ERP/财务管理
// @Security BearerAuth
// @Param id path int true "主键ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/finance/payment/query/{id} [get]
func (o *FinancePayment) GetByID(c *gin.Context) {
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

// Set 新增或修改 ERP 付款单。
// @Summary 新增或修改 ERP 付款单
// @Description 新增或修改 ERP 付款单
// @Tags ERP/财务管理
// @Security BearerAuth
// @Accept application/json
// @Param body body financemodels.FinancePaymentUpsert true "ERP 付款单参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/finance/payment/set [post]
func (o *FinancePayment) Set(c *gin.Context) {
	req := new(financemodels.FinancePaymentUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *financemodels.FinancePayment
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

// Delete 删除 ERP 付款单。
// @Summary 删除 ERP 付款单
// @Description 根据ID删除 ERP 付款单，多个ID用逗号分隔
// @Tags ERP/财务管理
// @Security BearerAuth
// @Param ids path string true "主键ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/finance/payment/remove/{ids} [delete]
func (o *FinancePayment) Delete(c *gin.Context) {
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
