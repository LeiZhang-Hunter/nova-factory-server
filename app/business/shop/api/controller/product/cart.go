package product

import (
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Cart struct {
	shopCartService service.IApiShopCartService
}

func NewCart(shopCartService service.IApiShopCartService) *Cart {
	return &Cart{
		shopCartService: shopCartService,
	}
}

func (c *Cart) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/cart")
	group.POST("/add", c.Add)
	group.GET("/list", c.List)
}

func (c *Cart) PrivateRoutes(router *gin.RouterGroup) {
}

// Add 加入购物车
// @Summary 加入购物车
// @Description 加入购物车（同一 SKU 重复加入会累加数量）
// @Tags app接口/商城/App购物车
// @Produce application/json
// @Param object body models.CartSetDataReq true "购物车新增修改参数"
// @Success 200 {object} response.ResponseData "加入成功"
// @Router /api/v1/app/shop/cart/add [post]
func (c *Cart) Add(ctx *gin.Context) {
	req := new(models.CartSetDataReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		zap.L().Error("bind json error", zap.Error(err))
		baizeContext.ParameterError(ctx)
		return
	}

	info, err := c.shopCartService.GenCart(ctx, req)
	if err != nil {
		zap.L().Error("gen cart error", zap.Error(err))
		baizeContext.Waring(ctx, err.Error())
		return
	}

	baizeContext.SuccessData(ctx, info)
}

// List 查询用户购物车列表
// @Summary 查询用户购物车列表
// @Description 查询当前登录用户购物车列表
// @Tags app接口/商城/App购物车
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/cart/list [get]
func (c *Cart) List(ctx *gin.Context) {
	data, err := c.shopCartService.List(ctx)
	if err != nil {
		zap.L().Error("list cart error", zap.Error(err))
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}
