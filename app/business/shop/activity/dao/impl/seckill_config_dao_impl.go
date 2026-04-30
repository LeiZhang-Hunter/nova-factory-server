package impl

import (
	"errors"
	"time"

	"nova-factory-server/app/business/shop/activity/dao"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ShopSeckillConfigDaoImpl 提供商品秒杀配置的数据访问能力。
type ShopSeckillConfigDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewShopSeckillConfigDao 创建商品秒杀配置 DAO。
func NewShopSeckillConfigDao(ms *gorm.DB) dao.IShopSeckillConfigDao {
	return &ShopSeckillConfigDaoImpl{
		db:        ms,
		tableName: "shop_store_seckill_config",
	}
}

// Set 新增或修改商品秒杀配置。
func (s *ShopSeckillConfigDaoImpl) Set(c *gin.Context, req *models.SeckillConfigSet) (*models.SeckillConfig, error) {
	if req.ID > 0 {
		return s.update(c, req)
	}
	return s.create(c, req)
}

// DeleteByIDs 软删除商品秒杀配置。
func (s *ShopSeckillConfigDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	now := time.Now()
	return s.db.WithContext(c).Table(s.tableName).
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": &now,
		}).Error
}

// GetByID 根据主键获取商品秒杀配置。
func (s *ShopSeckillConfigDaoImpl) GetByID(c *gin.Context, id int64) (*models.SeckillConfig, error) {
	var item models.SeckillConfig
	if err := s.baseQuery(c).Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// List 分页查询商品秒杀配置列表。
func (s *ShopSeckillConfigDaoImpl) List(c *gin.Context, req *models.SeckillConfigQuery) (*models.SeckillConfigListData, error) {
	db := s.baseQuery(c)
	if req.BeginClock != nil {
		db = db.Where("begin_clock = ?", *req.BeginClock)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	if req.Size > 200 {
		req.Size = 200
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*models.SeckillConfig, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).
		Limit(int(req.Size)).
		Order("sort ASC").
		Order("id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.SeckillConfigListData{Rows: rows, Total: total}, nil
}

func (s *ShopSeckillConfigDaoImpl) baseQuery(c *gin.Context) *gorm.DB {
	return s.db.WithContext(c).Table(s.tableName).
		Where("state = ?", commonStatus.NORMAL)
}

func (s *ShopSeckillConfigDaoImpl) create(c *gin.Context, req *models.SeckillConfigSet) (*models.SeckillConfig, error) {
	model := &models.SeckillConfig{
		ID:            snowflake.GenID(),
		BeginClock:    req.BeginClock,
		ContinueClock: req.ContinueClock,
		Images:        req.Images,
		Sort:          req.Sort,
		DeptID:        baizeContext.GetDeptId(c),
		State:         commonStatus.NORMAL,
	}
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ShopSeckillConfigDaoImpl) update(c *gin.Context, req *models.SeckillConfigSet) (*models.SeckillConfig, error) {
	now := time.Now()
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", req.ID).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"begin_clock":    req.BeginClock,
			"continue_clock": req.ContinueClock,
			"images":         req.Images,
			"sort":           req.Sort,
			"status":         req.Status,
			"update_by":      baizeContext.GetUserId(c),
			"update_time":    &now,
		}).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, req.ID)
}
