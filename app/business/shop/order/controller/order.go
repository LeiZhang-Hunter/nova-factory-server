package controller

import (
	"nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/business/shop/order/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// Order Shop订单控制器。
type Order struct {
	service service.IOrderService
}

// NewOrder 创建 Shop订单控制器。
func NewOrder(service service.IOrderService) *Order {
	return &Order{service: service}
}

// PrivateRoutes 注册 Shop订单相关私有路由。
func (o *Order) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/order")
	group.GET("/list", middlewares.HasPermission("shop:order:list"), o.List)
	group.GET("/refund/list", middlewares.HasPermission("shop:order:refund:list"), o.RefundList)
	group.POST("/refund/review", middlewares.HasPermission("shop:order:refund:review"), o.ReviewRefund)
	group.POST("/refund/refund", middlewares.HasPermission("shop:order:refund:refund"), o.Refund)
	group.GET("/query/:id", middlewares.HasPermission("shop:order:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("shop:order:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:order:remove"), o.Delete)
	group.POST("/synchronize-sales-orders", middlewares.HasPermission("shop:order:synchronizeSalesOrders"), o.SynchronizeSalesOrders)
	group.POST("/review-payment-voucher", middlewares.HasPermission("shop:order:reviewPaymentVoucher"), o.ReviewPaymentVoucher)
}

func (o *Order) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/shop/order/list", "shop:order:list")
	router.RegisterPermission("GET", "/shop/order/refund/list", "shop:order:refund:list")
	router.RegisterPermission("POST", "/shop/order/refund/review", "shop:order:refund:review")
	router.RegisterPermission("POST", "/shop/order/refund/refund", "shop:order:refund:refund")
	router.RegisterPermission("GET", "/shop/order/query/:id", "shop:order:query")
	router.RegisterPermission("POST", "/shop/order/set", "shop:order:set")
	router.RegisterPermission("DELETE", "/shop/order/remove/:ids", "shop:order:remove")
	router.RegisterPermission("POST", "/shop/order/synchronize-sales-orders", "shop:order:synchronizeSalesOrders")
}

// List Shop订单列表
// @Summary Shop订单列表
// @Description 按条件分页查询Shop订单
// @Tags Shop/销售管理
// @Security BearerAuth
// @Param tid query string false "网店订单编号"
// @Param status query string false "订单状态"
// @Param syncStatus query int false "同步状态"
// @Param receiverName query string false "收货人名称"
// @Param page query int false "页码"
// @Param size query int false "每页条数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /Shop/order/list [get]
func (o *Order) List(c *gin.Context) {
	req := new(models.OrderQuery)
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

// RefundList Shop售后单列表
// @Summary Shop售后单列表
// @Description 按条件分页查询Shop售后单
// @Tags Shop/售后管理
// @Security BearerAuth
// @Param tid query string false "网店订单编号"
// @Param status query int false "售后状态：0待审核 1已审核 2退款中 3退款成功 4退款失败 5退款关闭 6已拒绝"
// @Param syncStatus query int false "ERP同步状态：0待同步 1成功 2失败"
// @Param page query int false "页码"
// @Param size query int false "每页条数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /shop/order/refund/list [get]
func (o *Order) RefundList(c *gin.Context) {
	req := new(models.RefundQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.ListRefunds(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID Shop订单详情
// @Summary Shop订单详情
// @Description 根据ID查询Shop订单详情
// @Tags Shop/销售管理
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /Shop/order/query/{id} [get]
func (o *Order) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.GetByID(c, uint64(id))
	if err != nil {
		zap.L().Error("get order by id failed", zap.Error(err), zap.Int64("id", id))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set Shop订单保存
// @Summary Shop订单保存
// @Description 新增或修改Shop订单
// @Tags Shop/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body models.OrderSet true "Shop订单保存参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /Shop/order/set [post]
func (o *Order) Set(c *gin.Context) {
	req := new(models.OrderSet)
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

// Delete Shop订单删除
// @Summary Shop订单删除
// @Description 根据ID删除Shop订单
// @Tags Shop/销售管理
// @Security BearerAuth
// @Param ids path string true "订单ID，多个以逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /Shop/order/remove/{ids} [delete]
func (o *Order) Delete(c *gin.Context) {
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

// SynchronizeSalesOrders 同步销售订单
// @Summary 同步销售订单
// @Description 调用管家婆订单同步接口推送销售订单
// @Tags Shop/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body grasp.OrderSyncRequest true "销售订单同步参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "同步成功"
// @Router /Shop/order/synchronize-sales-orders [post]
func (o *Order) SynchronizeSalesOrders(c *gin.Context) {
	req := new(models.OrderSyncRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if len(req.GetOrders()) == 0 {

	}
	data, err := o.service.SynchronizeSalesOrders(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// ReviewPaymentVoucher 审核支付凭证
// @Summary 审核支付凭证
// @Description 审核支付凭证
// @Tags Shop/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body models.PaymentVoucherReviewReq true "支付凭证审核参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "审核成功"
// @Router /Shop/order/review-payment-voucher [post]
func (o *Order) ReviewPaymentVoucher(c *gin.Context) {
	req := new(models.PaymentVoucherReviewReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if err := o.service.ReviewPaymentVoucher(c, req); err != nil {
		zap.L().Error("review payment voucher failed", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// ReviewRefund Shop售后单审核
// @Summary Shop售后单审核
// @Description 管理后台审核售后单；退货退款审核通过后同步ERP，仅退款不同步ERP
// @Tags Shop/售后管理
// @Security BearerAuth
// @Accept application/json
// @Param body body models.RefundReviewReq true "售后单审核参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "审核成功"
// @Router /shop/order/refund/review [post]
func (o *Order) ReviewRefund(c *gin.Context) {
	req := new(models.RefundReviewReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if err := o.service.ReviewRefund(c, req); err != nil {
		zap.L().Error("review shop order refund failed", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// Refund Shop售后单退款
// @Summary Shop售后单退款
// @Description 管理后台手动触发售后退款；仅待退款或退款失败状态允许操作
// @Tags Shop/售后管理
// @Security BearerAuth
// @Accept application/json
// @Param body body models.RefundManualReq true "售后单退款参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "退款发起成功"
// @Router /shop/order/refund/refund [post]
func (o *Order) Refund(c *gin.Context) {
	req := new(models.RefundManualReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if err := o.service.Refund(c, req); err != nil {
		zap.L().Error("refund shop order refund failed", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
