package product

import (
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Product 前台商品控制器。
type Product struct {
	service shopservice.IShopGoodsService
}

// NewProduct 创建前台商品控制器。
func NewProduct(service shopservice.IShopGoodsService) *Product {
	return &Product{service: service}
}

func (p *Product) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/product")
	group.GET("/info/:id", p.GetByID)
}

func (p *Product) PrivateRoutes(router *gin.RouterGroup) {

}

// GetByID 获取商品详情
// @Summary 获取商品详情
// @Description 根据ID获取商品详情
// @Tags app接口/商城/App商品
// @Param id path int true "商品ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/product/info/{id} [get]
func (p *Product) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := p.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
