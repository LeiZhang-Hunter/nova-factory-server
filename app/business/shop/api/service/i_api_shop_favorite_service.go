package service

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopFavoriteService 用户商品收藏服务接口
type IApiShopFavoriteService interface {
	// AddFavorite 添加收藏
	AddFavorite(c *gin.Context, userId int64, goodsId string) error
	// RemoveFavorite 移除收藏
	RemoveFavorite(c *gin.Context, userId int64, goodsId string) error
	// ListFavorites 获取收藏列表
	ListFavorites(c *gin.Context, userId int64, page int64, size int64, goodsName string) (*models.GoodsListData, error)
	// CheckFavoriteStatus 检查收藏状态
	CheckFavoriteStatus(c *gin.Context, userId int64, goodsId string) (bool, error)
}
