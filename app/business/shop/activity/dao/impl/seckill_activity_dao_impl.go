package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/activity/dao"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopSeckillActivityDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewShopSeckillActivityDao(ms *gorm.DB) dao.IShopSeckillActivityDao {
	return &ShopSeckillActivityDaoImpl{db: ms, tableName: "eb_store_seckill_activity"}
}

func (s *ShopSeckillActivityDaoImpl) Set(c *gin.Context, req *models.SeckillActivitySet) (*models.SeckillActivity, error) {
	if req.ID > 0 {
		return s.update(c, req)
	}
	return s.create(c, req)
}

func (s *ShopSeckillActivityDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
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

func (s *ShopSeckillActivityDaoImpl) GetByID(c *gin.Context, id int64) (*models.SeckillActivity, error) {
	var item models.SeckillActivity
	if err := s.baseQuery(c).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (s *ShopSeckillActivityDaoImpl) List(c *gin.Context, req *models.SeckillActivityQuery) (*models.SeckillActivityListData, error) {
	db := s.baseQuery(c)
	if title := strings.TrimSpace(req.Title); title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	if req.Type != nil {
		db = db.Where("type = ?", *req.Type)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.LinkID > 0 {
		db = db.Where("link_id = ?", req.LinkID)
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
	rows := make([]*models.SeckillActivity, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.SeckillActivityListData{Rows: rows, Total: total}, nil
}

func (s *ShopSeckillActivityDaoImpl) baseQuery(c *gin.Context) *gorm.DB {
	return s.db.WithContext(c).Table(s.tableName).
		Where("state = ?", commonStatus.NORMAL)
}

func (s *ShopSeckillActivityDaoImpl) create(c *gin.Context, req *models.SeckillActivitySet) (*models.SeckillActivity, error) {
	model := buildSeckillActivityModel(req)
	model.ID = snowflake.GenID()
	model.State = commonStatus.NORMAL
	model.DeptID = baizeContext.GetDeptId(c)
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ShopSeckillActivityDaoImpl) update(c *gin.Context, req *models.SeckillActivitySet) (*models.SeckillActivity, error) {
	now := time.Now()
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", req.ID).
		Where("state = ?", commonStatus.NORMAL).
		Updates(buildSeckillActivityUpdates(req, baizeContext.GetUserId(c), &now)).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, req.ID)
}

func buildSeckillActivityModel(req *models.SeckillActivitySet) *models.SeckillActivity {
	return &models.SeckillActivity{
		Type:         req.Type,
		Title:        req.Title,
		StartDay:     req.StartDay,
		EndDay:       req.EndDay,
		TimeIDs:      strings.Join(req.TimeIDs, ","),
		OnceNum:      req.OnceNum,
		Num:          req.Num,
		IsCommission: req.IsCommission,
		Status:       req.Status,
		LinkID:       req.LinkID,
		AddTime:      time.Now().Unix(),
	}
}

func buildSeckillActivityUpdates(req *models.SeckillActivitySet, userID int64, now *time.Time) map[string]any {
	return map[string]any{
		"type":          req.Type,
		"title":         req.Title,
		"start_day":     req.StartDay,
		"end_day":       req.EndDay,
		"time_ids":      strings.Join(req.TimeIDs, ","),
		"once_num":      req.OnceNum,
		"num":           req.Num,
		"is_commission": req.IsCommission,
		"status":        req.Status,
		"link_id":       req.LinkID,
		"update_by":     userID,
		"update_time":   now,
	}
}
