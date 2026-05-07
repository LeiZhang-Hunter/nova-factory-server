package shopcontroller

import (
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Cart 商城用户购物车控制器
type Cart struct {
	service service.IShopCartService
}

// NewCart 创建商城用户购物车控制器。
func NewCart(service service.IShopCartService) *Cart {
	return &Cart{service: service}
}

// PrivateRoutes 注册商城用户购物车路由
func (s *Cart) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/user/cart")
	group.GET("/list", middlewares.HasPermission("shop:user:cart:list"), s.List)
	group.GET("/info/:id", middlewares.HasPermission("shop:user:cart:info"), s.GetByID)
	group.POST("/set", middlewares.HasPermission("shop:user:cart:set"), s.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:user:cart:remove"), s.Remove)
}

// List 获取商城用户购物车列表
// @Summary 获取商城用户购物车列表
// @Description 获取商城用户购物车列表
// @Tags 商城/用户购物车
// @Param object query models.CartQuery true "商城用户购物车查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/cart/list [get]
func (s *Cart) List(c *gin.Context) {
	req := new(models.CartQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.List(c, req)
	if err != nil {
		zap.L().Error("get list failed", zap.Any("err", err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 获取商城用户购物车详情
// @Summary 获取商城用户购物车详情
// @Description 根据ID获取商城用户购物车详情
// @Tags 商城/用户购物车
// @Param id path int true "商城用户购物车ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/cart/info/{id} [get]
func (s *Cart) GetByID(c *gin.Context) {
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

// Set 新增或修改商城用户购物车
// @Summary 新增或修改商城用户购物车
// @Description 新增或修改商城用户购物车
// @Tags 商城/用户购物车
// @Param object body models.CartSetReq true "商城用户购物车设置参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /shop/user/cart/set [post]
func (s *Cart) Set(c *gin.Context) {
	req := new(models.CartSetReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Set(c, req)
	if err != nil {
		zap.L().Error("set cart error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Remove 删除商城用户购物车
// @Summary 删除商城用户购物车
// @Description 根据ID删除商城用户购物车
// @Tags 商城/用户购物车
// @Param ids path string true "商城用户购物车ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/user/cart/remove/{ids} [delete]
func (s *Cart) Remove(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.Remove(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
