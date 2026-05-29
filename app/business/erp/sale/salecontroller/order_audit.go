package salecontroller

import (
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// OrderAudit ERP订单审核控制器。
type OrderAudit struct {
	service saleservice.IOrderAuditService
}

// NewOrderAudit 创建 ERP订单审核控制器。
func NewOrderAudit(service saleservice.IOrderAuditService) *OrderAudit {
	return &OrderAudit{service: service}
}

// PrivateRoutes 注册 ERP订单审核私有路由。
func (o *OrderAudit) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/sale/order-audit")
	group.GET("/list", middlewares.HasPermission("erp:sale:orderAudit:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:sale:orderAudit:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:sale:orderAudit:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:sale:orderAudit:remove"), o.Delete)
	group.POST("/approve", middlewares.HasPermission("erp:sale:orderAudit:approve"), o.Approve)
	group.POST("/reject", middlewares.HasPermission("erp:sale:orderAudit:reject"), o.Reject)
	group.POST("/import", middlewares.HasPermission("erp:sale:orderAudit:import"), o.Import)
}

func (o *OrderAudit) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterSchema("POST", "/erp/order-audit/import", nil, salemodels.OrderAuditImportReq{})
}

// List ERP订单审核列表
// @Summary ERP订单审核列表
// @Description 按条件分页查询ERP订单审核记录
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param tid query string false "网店订单编号"
// @Param auditStatus query int false "审核状态"
// @Param transferStatus query int false "转正式订单状态"
// @Param receiverName query string false "收货人名称"
// @Param page query int false "页码"
// @Param size query int false "每页条数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/order-audit/list [get]
func (o *OrderAudit) List(c *gin.Context) {
	req := new(salemodels.OrderAuditQuery)
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

// GetByID ERP订单审核详情
// @Summary ERP订单审核详情
// @Description 根据ID查询ERP订单审核详情
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param id path int true "订单审核ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/order-audit/query/{id} [get]
func (o *OrderAudit) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.GetByID(c, uint64(id))
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set ERP订单审核保存
// @Summary ERP订单审核保存
// @Description 新增或修改ERP订单审核记录
// @Tags ERP/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body salemodels.OrderAuditSet true "ERP订单审核保存参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /erp/sale/order-audit/set [post]
func (o *OrderAudit) Set(c *gin.Context) {
	req := new(salemodels.OrderAuditSet)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Import ERP订单审核批量导入
// @Summary ERP订单审核批量导入
// @Description 批量导入订单审核记录
// @Tags ERP/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body salemodels.OrderAuditImportReq true "ERP订单审核批量导入参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "导入成功"
// @Router /erp/order-audit/import [post]
func (o *OrderAudit) Import(c *gin.Context) {
	req := new(salemodels.OrderAuditImportReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.Import(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete ERP订单审核删除
// @Summary ERP订单审核删除
// @Description 根据ID删除ERP订单审核记录
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param ids path string true "订单审核ID，多个以逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/sale/order-audit/remove/{ids} [delete]
func (o *OrderAudit) Delete(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	orderIDs := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		orderIDs = append(orderIDs, uint64(id))
	}
	if len(orderIDs) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := o.service.DeleteByIDs(c, orderIDs); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// Approve ERP订单审核通过
// @Summary ERP订单审核通过
// @Description 审核通过订单审核记录并转入正式ERP订单
// @Tags ERP/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body salemodels.OrderAuditAction true "订单审核动作参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "审核通过成功"
// @Router /erp/sale/order-audit/approve [post]
func (o *OrderAudit) Approve(c *gin.Context) {
	req := new(salemodels.OrderAuditAction)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.Approve(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Reject ERP订单审核驳回
// @Summary ERP订单审核驳回
// @Description 驳回订单审核记录
// @Tags ERP/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body salemodels.OrderAuditAction true "订单审核动作参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "驳回成功"
// @Router /erp/sale/order-audit/reject [post]
func (o *OrderAudit) Reject(c *gin.Context) {
	req := new(salemodels.OrderAuditAction)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.Reject(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
