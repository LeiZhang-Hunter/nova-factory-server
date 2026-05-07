package order

import (
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Order 订单控制器
type Order struct {
	service service.IApiShopOrderService
}

// NewOrder 创建订单控制器
func NewOrder(service service.IApiShopOrderService) *Order {
	return &Order{service: service}
}

// PrivateRoutes 注册商城订单路由（商城模块只检查登录，不检查权限）
func (s *Order) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/order")
	group.GET("/list", s.List)
	group.GET("/info/:id", s.GetByID)
	group.GET("/count", s.GetStatistics)
	group.POST("/create", s.Create)
	group.POST("/status", s.UpdateStatus)
	group.POST("/cancel/:id", s.Cancel)
	group.POST("/confirm/:id", s.ConfirmReceive)
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

// Create 创建订单
// @Summary 创建订单
// @Description 创建新订单
// @Tags 商城/用户订单
// @Param object body models.OrderSetReq true "订单创建参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "创建成功"
// @Router /shop/user/order/create [post]
func (s *Order) Create(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	req := new(models.OrderSetReq)
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

// UpdateStatus 更新订单状态
// @Summary 更新订单状态
// @Description 更新订单状态
// @Tags 商城/用户订单
// @Param object body models.OrderStatusReq true "订单状态更新参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "更新成功"
// @Router /shop/user/order/status [post]
func (s *Order) UpdateStatus(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	req := new(models.OrderStatusReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.UpdateStatus(c, userID, req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
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
