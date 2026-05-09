package product

import (
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Product struct {
	service service.IApiShopGoodsService
}

func NewProduct(service service.IApiShopGoodsService) *Product {
	return &Product{service: service}
}

func (p *Product) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/product")
	group.GET("/info/:id", p.Info)
	group.GET("/list", p.List)
}

func (p *Product) PrivateRoutes(router *gin.RouterGroup) {
	goods := router.Group("/api/v1/app/shop/goods")
	goods.GET("/repurchase", p.Repurchase)
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

// List 获取商品列表
// @Summary 获取商品列表
// @Description 根据分类ID获取商品列表
// @Tags app接口/商城/App商品
// @Produce application/json
// @Param categoryId query int false "商品分类ID"
// @Param goodsName query string false "商品名称"
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/product/list [get]
func (p *Product) List(c *gin.Context) {
	req := new(models.GoodsQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := p.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Repurchase 获取用户复购商品列表
// @Summary 获取用户复购商品列表
// @Description 获取用户历史购买过的商品列表
// @Tags app接口/商城/App商品
// @Produce application/json
// @Param categoryId query int false "商品分类ID"
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/goods/repurchase [get]
func (p *Product) Repurchase(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	if userID == 0 {
		baizeContext.InvalidToken(c)
		return
	}
	req := new(models.GoodsQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := p.service.ListRepurchase(c, userID, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
