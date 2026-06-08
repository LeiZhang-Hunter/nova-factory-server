package product

import (
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/constant/shop"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/store"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const categoryCacheRefreshInterval = 5 * time.Minute

type Category struct {
	service shopservice.IShopCategoryService
}

func NewCategory(service shopservice.IShopCategoryService) *Category {
	categoryCache := store.NewShopCategoryStore()
	controller := &Category{service: service, cache: categoryCache}
	store.RegisterStore(shop.ShopCategoryStoreName, categoryCache)
	controller.startCacheRefresh()
	return controller
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
	if c.cache != nil {
		c.cache.Set(toCategoryStoreData(data))
	}
	baizeContext.SuccessData(ctx, data)
}

// startCacheRefresh 开启缓存刷新
func (c *Category) startCacheRefresh() {
	c.refreshCache()
	go func() {
		ticker := time.NewTicker(categoryCacheRefreshInterval)
		defer ticker.Stop()
		for range ticker.C {
			c.refreshCache()
		}
	}()
}

func (c *Category) refreshCache() {
	if c == nil || c.service == nil || c.cache == nil {
		return
	}

	data, err := c.service.All(&gin.Context{})
	if err != nil {
		zap.L().Warn("refresh shop category cache failed", zap.Error(err))
		return
	}
	c.cache.Set(toCategoryStoreData(data))
}

func toCategoryStoreData(data []*shopmodels.CategoryInfo) []store.ShopCategoryData {
	categoryDataList := make([]store.ShopCategoryData, 0, len(data))
	for _, row := range data {
		categoryDataList = append(categoryDataList, row)
	}
	return categoryDataList
}
