package salecontroller

import (
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Order ERP订单控制器。
type Order struct {
	service saleservice.IOrderService
}

// NewOrder 创建 ERP订单控制器。
func NewOrder(service saleservice.IOrderService) *Order {
	return &Order{service: service}
}

// PrivateRoutes 注册 ERP订单相关私有路由。
func (o *Order) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/sale/order")
	group.GET("/list", middlewares.HasPermission("erp:sale:order:list"), o.List)
	group.GET("/query/:id", middlewares.HasPermission("erp:sale:order:query"), o.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:sale:order:set"), o.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("erp:sale:order:remove"), o.Delete)
	group.GET("/check-login-state", middlewares.HasPermission("erp:sale:order:checkLoginState"), o.CheckLoginState)
	group.POST("/synchronize-sales-orders", middlewares.HasPermission("erp:sale:order:synchronizeSalesOrders"), o.SynchronizeSalesOrders)
}

// List ERP订单列表
// @Summary ERP订单列表
// @Description 按条件分页查询ERP订单
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param tid query string false "网店订单编号"
// @Param status query string false "订单状态"
// @Param syncStatus query int false "同步状态"
// @Param receiverName query string false "收货人名称"
// @Param page query int false "页码"
// @Param size query int false "每页条数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/order/list [get]
func (o *Order) List(c *gin.Context) {
	req := new(salemodels.OrderQuery)
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

// GetByID ERP订单详情
// @Summary ERP订单详情
// @Description 根据ID查询ERP订单详情
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/order/query/{id} [get]
func (o *Order) GetByID(c *gin.Context) {
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

// Set ERP订单保存
// @Summary ERP订单保存
// @Description 新增或修改ERP订单
// @Tags ERP/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body salemodels.OrderSet true "ERP订单保存参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /erp/order/set [post]
func (o *Order) Set(c *gin.Context) {
	req := new(salemodels.OrderSet)
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

// Delete ERP订单删除
// @Summary ERP订单删除
// @Description 根据ID删除ERP订单
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param ids path string true "订单ID，多个以逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/order/remove/{ids} [delete]
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

// CheckLoginState 检查管家婆登录状态
// @Summary 检查管家婆登录状态
// @Description 检查当前启用 ERP 集成配置的登录状态
// @Tags ERP/销售管理
// @Security BearerAuth
// @Param checkUrl query string false "登录状态检查地址"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/order/check-login-state [get]
func (o *Order) CheckLoginState(c *gin.Context) {
	req := new(salemodels.CheckLoginStateReq)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.CheckLoginState(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// SynchronizeSalesOrders 同步销售订单
// @Summary 同步销售订单
// @Description 调用管家婆订单同步接口推送销售订单
// @Tags ERP/销售管理
// @Security BearerAuth
// @Accept application/json
// @Param body body grasp.OrderSyncRequest true "销售订单同步参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "同步成功"
// @Router /erp/order/synchronize-sales-orders [post]
func (o *Order) SynchronizeSalesOrders(c *gin.Context) {
	req := new(salemodels.OrderSyncRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.SynchronizeSalesOrders(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
