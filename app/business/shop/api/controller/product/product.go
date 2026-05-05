package product

import (
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Product struct {
	service shopservice.IShopGoodsService
}

func NewProduct(service shopservice.IShopGoodsService) *Product {
	return &Product{service: service}
}

func (p *Product) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/product")
	group.GET("/info/:id", p.Info)
}

func (p *Product) PrivateRoutes(router *gin.RouterGroup) {

}

// Info 读取商品详情
// @Summary 读取商品详情
// @Description 根据ID读取商品详情
// @Tags app接口/商城/App商品
// @Produce application/json
// @Param id path int true "商品ID"
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/product/info/{id} [get]
func (p *Product) Info(c *gin.Context) {
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
