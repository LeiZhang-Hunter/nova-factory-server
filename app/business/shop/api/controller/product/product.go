package product

import (
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/constant/shop"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"
	"nova-factory-server/app/utils/store"

	"github.com/gin-gonic/gin"
)

type Product struct {
	service service.IApiShopGoodsService
}

func NewProduct(service service.IApiShopGoodsService) *Product {
	return &Product{
		service: service,
	}
}

func (p *Product) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/product")
	group.GET("/info/:id", p.Info)
	group.GET("/list", p.List)
	group.POST("/search", p.Search)
}

func (p *Product) PrivateRoutes(router *gin.RouterGroup) {
	product := router.Group("/api/v1/app/shop/product")
	product.GET("/repurchase", p.Repurchase)
}

func (p *Product) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterSchema("POST", "/api/v1/app/shop/product/search", nil, models.GoodsSearchReq{})
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
	if req.CategoryId != 0 || req.CategoryIds != nil {
		categoryStore := store.GetStore(shop.ShopCategoryStoreName)
		if categoryStore == nil {
			baizeContext.Waring(c, "分类存储初始化失败")
			return
		}

		req.CategoryIds = categoryStore.GetCategoryIDs(req.CategoryId)
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

// Search 商品检索
// @Summary 商品检索
// @Description 传入多个商品名称，基于商品向量检索相似商品并回填数据库中的最新商品数据
// @Tags app接口/商城/App商品
// @Accept application/json
// @Produce application/json
// @Param object body models.GoodsSearchReq true "商品检索参数"
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/product/search [post]
func (p *Product) Search(c *gin.Context) {
	req := new(models.GoodsSearchReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.Limit <= 0 {
		req.Limit = 1
	}
	if req.Limit > 10 {
		req.Limit = 10
	}
	if len(req.GoodsNames) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if len(req.GoodsNames) > 10 {
		baizeContext.Waring(c, "输入商品太多了，最多50条")
		return
	}
	isSale := true
	req.IsSale = &isSale
	data, err := p.service.Search(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
