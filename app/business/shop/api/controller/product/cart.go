package product

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"
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
	group.GET("/add", c.Add)
}

func (c *Cart) PrivateRoutes(router *gin.RouterGroup) {

}

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
		return
	}

	baizeContext.SuccessData(ctx, info)
}
