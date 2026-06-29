package order

import (
	"net/http"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Order 订单控制器
type Order struct {
	service            service.IApiShopOrderService
	orderRefundService service.IApiShopOrderRefundService
}

// NewOrder 创建订单控制器
func NewOrder(service service.IApiShopOrderService, orderRefundService service.IApiShopOrderRefundService) *Order {
	return &Order{service: service, orderRefundService: orderRefundService}
}

// PrivateRoutes 注册商城订单路由（商城模块只检查登录，不检查权限）
func (s *Order) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/order")
	group.GET("/list", s.List)
	group.GET("/info/:id", s.GetByID)
	group.GET("/count", s.GetStatistics)
	group.POST("/confirm", s.Confirm)
	group.POST("/create", s.Create)
	group.POST("/pay/:id", s.Pay)
	group.POST("/cancel/:id", s.Cancel)
	group.POST("/confirm/:id", s.ConfirmReceive)
	// 售后/退款申请
	group.POST("/refund/apply", s.Apply)
}

// List 获取订单列表
// @Summary 获取订单列表
// @Description 获取当前用户的订单列表
// @Tags 商城/用户订单
// @Param object query models.OrderQuery true "订单查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/order/list [get]
func (s *Order) List(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	req := new(models.OrderQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.List(c, userID, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 获取订单详情
// @Summary 获取订单详情
// @Description 根据ID获取订单详情
// @Tags 商城/用户订单
// @Param id path int true "订单ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/order/info/{id} [get]
func (s *Order) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Confirm 获取确认单数据
// @Summary 获取确认单数据
// @Description 根据当前地址、配送方式等信息实时试算确认单
// @Tags app接口/商城/App订单
// @Param object body models.OrderConfirmReq true "确认单参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/order/confirm [post]
func (s *Order) Confirm(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	req := new(models.OrderConfirmReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Confirm(c, userID, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Create 创建订单
// @Summary 创建订单
// @Description 根据预订单 orderKey 正式落订单
// @Tags app接口/商城/App订单
// @Param object body models.OrderCreateReq true "订单创建参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "创建成功"
// @Router /api/v1/app/shop/order/create [post]
func (s *Order) Create(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	req := new(models.OrderCreateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Create(c, userID, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Pay 支付订单
// @Summary 支付订单
// @Description 支付待付款订单，可通过payChannel指定支付通道（0=使用订单默认通道 1=微信 2=支付宝）
// @Tags 商城/用户订单
// @Param id path int true "订单ID"
// @Param object body models.OrderPayReq false "支付通道参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "支付成功"
// @Router /api/v1/app/shop/order/pay/{id} [post]
func (s *Order) Pay(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	var req models.OrderPayReq
	_ = c.ShouldBindJSON(&req) // 空 body/格式错误 → PayChannel=0 → 走订单默认通道
	data, err := s.service.Pay(c, userID, id, req.PayChannel)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Cancel 取消订单
// @Summary 取消订单
// @Description 取消待支付订单
// @Tags 商城/用户订单
// @Param id path int true "订单ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "取消成功"
// @Router /shop/user/order/cancel/{id} [post]
func (s *Order) Cancel(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	reason := c.Query("reason")
	if err := s.service.Cancel(c, userID, id, reason); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// ConfirmReceive 确认收货
// @Summary 确认收货
// @Description 确认已发货订单收货
// @Tags 商城/用户订单
// @Param id path int true "订单ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "确认成功"
// @Router /shop/user/order/confirm/{id} [post]
func (s *Order) ConfirmReceive(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.ConfirmReceive(c, userID, id); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// GetStatistics 获取订单统计
// @Summary 获取订单统计
// @Description 获取用户各状态订单数量统计
// @Tags 商城/用户订单
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/order/count [get]
func (s *Order) GetStatistics(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	data, err := s.service.GetStatistics(c, userID)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Apply 申请售后
// @Summary 申请售后
// @Description 申请售后
// @Tags 商城/售后
// @Param object body models.RefundApplyReq true "售后申请参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "申请成功"
// @Router /shop/user/refund/apply [post]
func (s *Order) Apply(c *gin.Context) {
	var req models.RefundApplyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	userID := baizeContext.GetUserId(c)
	if userID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "请先登录"})
		return
	}
	resp, err := s.orderRefundService.Apply(c, userID, &req)
	if err != nil {
		zap.L().Error("申请售后失败", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": resp})
}
