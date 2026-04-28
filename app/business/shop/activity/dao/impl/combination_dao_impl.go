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

type ShopCombinationDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewShopCombinationDao(ms *gorm.DB) dao.IShopCombinationDao {
	return &ShopCombinationDaoImpl{db: ms, tableName: "shop_store_combination"}
}

func (s *ShopCombinationDaoImpl) Set(c *gin.Context, req *models.CombinationSet) (*models.Combination, error) {
	if req.ID > 0 {
		return s.update(c, req)
	}
	return s.create(c, req)
}

func (s *ShopCombinationDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
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

func (s *ShopCombinationDaoImpl) GetByID(c *gin.Context, id int64) (*models.Combination, error) {
	var item models.Combination
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (s *ShopCombinationDaoImpl) List(c *gin.Context, req *models.CombinationQuery) (*models.CombinationListData, error) {
	db := s.db.WithContext(c).Table(s.tableName)
	if title := strings.TrimSpace(req.Title); title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	if req.ProductID > 0 {
		db = db.Where("product_id = ?", req.ProductID)
	}
	if req.IsShow != nil {
		db = db.Where("is_show = ?", *req.IsShow)
	}
	if req.IsHost != nil {
		db = db.Where("is_host = ?", *req.IsHost)
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
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*models.Combination, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order("sort ASC").Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.CombinationListData{Rows: rows, Total: total}, nil
}

func (s *ShopCombinationDaoImpl) create(c *gin.Context, req *models.CombinationSet) (*models.Combination, error) {
	model := buildCombinationModel(req)
	model.ID = snowflake.GenID()
	model.SetCreateBy(baizeContext.GetUserId(c))
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ShopCombinationDaoImpl) update(c *gin.Context, req *models.CombinationSet) (*models.Combination, error) {
	now := time.Now()
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", req.ID).
		Where("state = ?", commonStatus.NORMAL).
		Updates(buildCombinationUpdates(req, baizeContext.GetUserId(c), &now)).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, req.ID)
}

func buildCombinationModel(req *models.CombinationSet) *models.Combination {
	return &models.Combination{
		ProductID:     req.ProductID,
		MerID:         req.MerID,
		Image:         req.Image,
		Images:        req.Images,
		Title:         req.Title,
		Attr:          req.Attr,
		People:        req.People,
		Info:          req.Info,
		Price:         req.Price,
		Sort:          req.Sort,
		Sales:         req.Sales,
		Stock:         req.Stock,
		IsHost:        req.IsHost,
		IsShow:        req.IsShow,
		IsPostage:     req.IsPostage,
		Postage:       req.Postage,
		StartTime:     req.StartTime,
		StopTime:      req.StopTime,
		EffectiveTime: req.EffectiveTime,
		Browse:        req.Browse,
		UnitName:      req.UnitName,
		Weight:        req.Weight,
		Volume:        req.Volume,
		Num:           req.Num,
		OnceNum:       req.OnceNum,
		Quota:         req.Quota,
		QuotaShow:     req.QuotaShow,
		Virtual:       req.Virtual,
	}
}

func buildCombinationUpdates(req *models.CombinationSet, userID int64, now *time.Time) map[string]any {
	updates := map[string]any{
		"product_id":     req.ProductID,
		"mer_id":         req.MerID,
		"image":          req.Image,
		"images":         req.Images,
		"title":          req.Title,
		"attr":           req.Attr,
		"people":         req.People,
		"info":           req.Info,
		"price":          req.Price,
		"sort":           req.Sort,
		"sales":          req.Sales,
		"stock":          req.Stock,
		"is_host":        req.IsHost,
		"is_show":        req.IsShow,
		"is_postage":     req.IsPostage,
		"postage":        req.Postage,
		"start_time":     req.StartTime,
		"stop_time":      req.StopTime,
		"effective_time": req.EffectiveTime,
		"browse":         req.Browse,
		"unit_name":      req.UnitName,
		"weight":         req.Weight,
		"volume":         req.Volume,
		"num":            req.Num,
		"once_num":       req.OnceNum,
		"quota":          req.Quota,
		"quota_show":     req.QuotaShow,
		"virtual":        req.Virtual,
		"update_by":      userID,
		"update_time":    now,
	}
	return updates
}
