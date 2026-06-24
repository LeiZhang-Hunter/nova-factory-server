package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopFavoriteDao 用户商品收藏数据访问接口
type IApiShopFavoriteDao interface {
	// Add 添加收藏
	Add(c *gin.Context, userId int64, goodsId int64) error
	// Remove 移除收藏
	Remove(c *gin.Context, userId int64, goodsId int64) error
	// ListByUserID 获取用户收藏列表
	ListByUserID(c *gin.Context, userId int64, page int64, size int64, goodsName string) ([]*models.ShopUserGoodsFavorite, int64, error)
	// GetByUserAndGoods 获取指定用户的指定商品收藏记录
	GetByUserAndGoods(c *gin.Context, userId int64, goodsId string) (*models.ShopUserGoodsFavorite, error)
	// CountByUserID 统计用户收藏数量
	CountByUserID(c *gin.Context, userId int64) (int64, error)
}
