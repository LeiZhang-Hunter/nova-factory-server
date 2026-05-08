package favorite

import (
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(NewFavorite)

// Favorite 收藏控制器
type Favorite struct {
	service service.IApiShopFavoriteService
}

// NewFavorite 创建收藏控制器
func NewFavorite(service service.IApiShopFavoriteService) *Favorite {
	return &Favorite{service: service}
}

// PrivateRoutes 注册商城收藏路由（商城模块只检查登录，不检查权限）
func (f *Favorite) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/favorite")
	group.POST("", f.Add)
	group.DELETE("/:goodsId", f.Remove)
	group.GET("/favorites", f.List)
	group.GET("/:goodsId/status", f.Status)
}

// Add 添加收藏
// @Summary 添加收藏
// @Description 添加商品到收藏列表
// @Tags 商城/用户收藏
// @Param object body models.FavoriteAddReq true "收藏参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "添加成功"
// @Router /api/v1/app/shop/favorite [post]
func (f *Favorite) Add(c *gin.Context) {
	userId := baizeContext.GetUserId(c)
	req := new(models.FavoriteAddReq)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("参数错误", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	if req.GoodsId == "" {
		baizeContext.ParameterError(c)
		return
	}
	err := f.service.AddFavorite(c, userId, req.GoodsId)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// Remove 移除收藏
// @Summary 移除收藏
// @Description 从收藏列表移除商品
// @Tags 商城/用户收藏
// @Param goodsId path string true "商品ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "移除成功"
// @Router /api/v1/app/shop/favorite/{goodsId} [delete]
func (f *Favorite) Remove(c *gin.Context) {
	userId := baizeContext.GetUserId(c)
	goodsId := c.Param("goodsId")
	if goodsId == "" {
		baizeContext.ParameterError(c)
		return
	}
	err := f.service.RemoveFavorite(c, userId, goodsId)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// List 获取收藏列表
// @Summary 获取收藏列表
// @Description 获取当前用户的商品收藏列表
// @Tags 商城/用户收藏
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/favorite/favorites [get]
func (f *Favorite) List(c *gin.Context) {
	userId := baizeContext.GetUserId(c)
	page := baizeContext.QueryInt64(c, "page")
	size := baizeContext.QueryInt64(c, "size")

	data, err := f.service.ListFavorites(c, userId, page, size)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Status 查询收藏状态
// @Summary 查询收藏状态
// @Description 查询指定商品是否被收藏
// @Tags 商城/用户收藏
// @Param goodsId path string true "商品ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/favorite/{goodsId}/status [get]
func (f *Favorite) Status(c *gin.Context) {
	userId := baizeContext.GetUserId(c)
	goodsId := c.Param("goodsId")
	if goodsId == "" {
		baizeContext.ParameterError(c)
		return
	}
	isFavorite, err := f.service.CheckFavoriteStatus(c, userId, goodsId)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, &models.FavoriteStatusResp{IsFavorite: isFavorite})
}
