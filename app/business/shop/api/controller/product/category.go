package product

import (
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Category struct {
	service shopservice.IShopCategoryService
}

func NewCategory(service shopservice.IShopCategoryService) *Category {
	return &Category{service: service}
}

func (c *Category) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/category")
	group.GET("/all", c.All)
}

func (c *Category) PrivateRoutes(router *gin.RouterGroup) {

}

// All 读取全部分类
// @Summary 读取全部商品分类
// @Description 读取商城全部商品分类，并按父子关系返回树形结构
// @Tags app接口/商城/App商品分类
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/category/all [get]
func (c *Category) All(ctx *gin.Context) {
	data, err := c.service.All(ctx)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}
