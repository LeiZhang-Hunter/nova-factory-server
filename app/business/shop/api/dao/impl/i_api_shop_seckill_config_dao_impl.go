package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/fileUtils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApiShopSeckillConfigDaoImpl 秒杀配置数据访问实现
type IApiShopSeckillConfigDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopSeckillConfigDaoImpl 创建秒杀配置数据访问对象
func NewIApiShopSeckillConfigDaoImpl(ms *gorm.DB) dao.IApiShopSeckillConfigDao {
	return &IApiShopSeckillConfigDaoImpl{
		db:        ms,
		tableName: "shop_store_seckill_config",
	}
}

// List 查询秒杀配置列表
func (s *IApiShopSeckillConfigDaoImpl) List(c *gin.Context, query *models.SeckillConfigQuery) (*models.SeckillConfigListData, error) {
	db := s.db.WithContext(c).Table(s.tableName).
		Where("state = ?", commonStatus.NORMAL)

	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
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

	rows := make([]*models.SeckillConfig, 0)
	if err := db.Offset(int((query.Page - 1) * query.Size)).
		Limit(int(query.Size)).
		Order("sort DESC").
		Order("id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		if row != nil {
			row.Images = fileUtils.BuildAbsoluteURL(c, row.Images)
		}
	}

	return &models.SeckillConfigListData{Rows: rows, Total: total}, nil
}

// GetByID 根据主键获取秒杀配置
func (s *IApiShopSeckillConfigDaoImpl) GetByID(c *gin.Context, id int64) (*models.SeckillConfig, error) {
	var item models.SeckillConfig
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	item.Images = fileUtils.BuildAbsoluteURL(c, item.Images)
	return &item, nil
}
