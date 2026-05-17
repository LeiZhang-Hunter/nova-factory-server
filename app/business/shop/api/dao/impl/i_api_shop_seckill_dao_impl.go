package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/fileUtils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApiShopSeckillDaoImpl 秒杀商品数据访问实现
type IApiShopSeckillDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopSeckillDaoImpl 创建秒杀商品数据访问对象
func NewIApiShopSeckillDaoImpl(ms *gorm.DB) dao.IApiShopSeckillDao {
	return &IApiShopSeckillDaoImpl{
		db:        ms,
		tableName: "shop_store_seckill",
	}
}

// List 查询秒杀商品列表
func (s *IApiShopSeckillDaoImpl) List(c *gin.Context, query *models.SeckillQuery) (*models.SeckillListData, error) {
	db := s.db.WithContext(c).Table(s.tableName).
		Select("shop_store_seckill.*, COALESCE(NULLIF(shop_store_seckill.image, ''), shop_goods.image_url) AS image, COALESCE(NULLIF(shop_store_seckill.title, ''), shop_goods.goods_name) AS title").
		Joins("LEFT JOIN shop_goods ON shop_store_seckill.product_id = shop_goods.id").
		Where("shop_store_seckill.state = ?", commonStatus.NORMAL).
		Where("shop_store_seckill.is_del = ?", 0)

	if title := strings.TrimSpace(query.Title); title != "" {
		db = db.Where("shop_store_seckill.title LIKE ?", "%"+title+"%")
	}
	if query.ActivityID > 0 {
		db = db.Where("shop_store_seckill.activity_id = ?", query.ActivityID)
	}
	if query.ProductID > 0 {
		db = db.Where("shop_store_seckill.product_id = ?", query.ProductID)
	}
	if query.Status != nil {
		db = db.Where("shop_store_seckill.status = ?", *query.Status)
	}
	if query.IsShow != nil {
		db = db.Where("shop_store_seckill.is_show = ?", *query.IsShow)
	}
	if query.IsHot != nil {
		db = db.Where("shop_store_seckill.is_hot = ?", *query.IsHot)
	}
	if query.TimeID > 0 {
		db = db.Where("FIND_IN_SET(?, shop_store_seckill.time_id) > 0", query.TimeID)
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 20
	}
	if query.Size > 200 {
		query.Size = 200
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	rows := make([]*models.Seckill, 0)
	if err := db.Offset(int((query.Page - 1) * query.Size)).
		Limit(int(query.Size)).
		Order("sort ASC").
		Order("id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		if row != nil {
			row.Image = fileUtils.BuildAbsoluteURL(c, row.Image)
		}
	}

	return &models.SeckillListData{Rows: rows, Total: total}, nil
}

// GetByID 根据主键获取秒杀商品
func (s *IApiShopSeckillDaoImpl) GetByID(c *gin.Context, id int64) (*models.Seckill, error) {
	var item models.Seckill
	if err := s.db.WithContext(c).Table(s.tableName).
		Select("shop_store_seckill.*, COALESCE(NULLIF(shop_store_seckill.image, ''), shop_goods.image_url) AS image, COALESCE(NULLIF(shop_store_seckill.title, ''), shop_goods.goods_name) AS title").
		Joins("LEFT JOIN shop_goods ON shop_store_seckill.product_id = shop_goods.id").
		Where("shop_store_seckill.id = ?", id).
		Where("shop_store_seckill.state = ?", commonStatus.NORMAL).
		Where("shop_store_seckill.is_del = ?", 0).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	item.Image = fileUtils.BuildAbsoluteURL(c, item.Image)
	return &item, nil
}
