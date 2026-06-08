package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApiShopFavoriteDaoImpl 用户商品收藏数据访问实现
type IApiShopFavoriteDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopFavoriteDaoImpl 创建用户商品收藏数据访问对象
func NewIApiShopFavoriteDaoImpl(ms *gorm.DB) dao.IApiShopFavoriteDao {
	return &IApiShopFavoriteDaoImpl{
		db:        ms,
		tableName: "shop_user_goods_favorite",
	}
}

// Add 添加收藏
func (d *IApiShopFavoriteDaoImpl) Add(c *gin.Context, userId int64, goodsId string) error {
	favorite := &models.ShopUserGoodsFavorite{
		UserID:  userId,
		GoodsID: goodsId,
	}
	return d.db.WithContext(c).Table(d.tableName).Create(favorite).Error
}

// Remove 移除收藏
func (d *IApiShopFavoriteDaoImpl) Remove(c *gin.Context, userId int64, goodsId string) error {
	return d.db.WithContext(c).Table(d.tableName).
		Where("user_id = ? AND goods_id = ?", userId, goodsId).
		Delete(nil).Error
}

// ListByUserID 获取用户收藏列表
func (d *IApiShopFavoriteDaoImpl) ListByUserID(c *gin.Context, userId int64, page int64, size int64, goodsName string) ([]*models.ShopUserGoodsFavorite, int64, error) {
	var rows []*models.ShopUserGoodsFavorite
	var total int64

	baseQuery := func() *gorm.DB {
		db := d.db.WithContext(c).Table(d.tableName+" AS f").Where("f.user_id = ?", userId)
		if goodsName != "" {
			db = db.Joins("JOIN shop_goods ON shop_goods.goods_id = f.goods_id").
				Where("shop_goods.goods_name LIKE ?", "%"+goodsName+"%")
		}
		return db
	}

	if err := baseQuery().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	if err := baseQuery().Select("f.*").Offset(int((page - 1) * size)).Limit(int(size)).Order("f.id DESC").Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// GetByUserAndGoods 获取指定用户的指定商品收藏记录
func (d *IApiShopFavoriteDaoImpl) GetByUserAndGoods(c *gin.Context, userId int64, goodsId string) (*models.ShopUserGoodsFavorite, error) {
	var item models.ShopUserGoodsFavorite
	if err := d.db.WithContext(c).Table(d.tableName).
		Where("user_id = ? AND goods_id = ?", userId, goodsId).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// CountByUserID 统计用户收藏数量
func (d *IApiShopFavoriteDaoImpl) CountByUserID(c *gin.Context, userId int64) (int64, error) {
	var count int64
	if err := d.db.WithContext(c).Table(d.tableName).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
