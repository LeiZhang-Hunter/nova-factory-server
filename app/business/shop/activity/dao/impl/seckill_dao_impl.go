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

type ShopSeckillDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewShopSeckillDao(ms *gorm.DB) dao.IShopSeckillDao {
	return &ShopSeckillDaoImpl{db: ms, tableName: "shop_store_seckill"}
}

func (s *ShopSeckillDaoImpl) Set(c *gin.Context, req *models.SeckillSet) (*models.Seckill, error) {
	if req.ID > 0 {
		return s.update(c, req)
	}
	return s.create(c, req)
}

// BatchCreate 批量新增秒杀商品。
func (s *ShopSeckillDaoImpl) BatchCreate(c *gin.Context, reqs []*models.SeckillSet, batchSize int) error {
	if len(reqs) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = len(reqs)
	}
	models := make([]*models.Seckill, 0, len(reqs))
	now := time.Now()
	for _, req := range reqs {
		if req == nil {
			continue
		}
		model := buildSeckillModel(req)
		model.ID = snowflake.GenID()
		model.IsDel = 0
		model.State = commonStatus.NORMAL
		model.DeptID = baizeContext.GetDeptId(c)
		model.CreateBy = baizeContext.GetUserId(c)
		model.UpdateBy = baizeContext.GetUserId(c)
		model.CreateTime = &now
		model.UpdateTime = &now
		models = append(models, model)
	}
	if len(models) == 0 {
		return nil
	}
	return s.db.WithContext(c).Table(s.tableName).CreateInBatches(models, batchSize).Error
}

// BatchUpdate 批量更新秒杀商品。
func (s *ShopSeckillDaoImpl) BatchUpdate(c *gin.Context, reqs []*models.SeckillSet, batchSize int) error {
	if len(reqs) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = len(reqs)
	}
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		for start := 0; start < len(reqs); start += batchSize {
			end := start + batchSize
			if end > len(reqs) {
				end = len(reqs)
			}
			for _, req := range reqs[start:end] {
				if req == nil {
					continue
				}
				now := time.Now()
				if err := tx.Table(s.tableName).
					Where("id = ?", req.ID).
					Where("state = ?", commonStatus.NORMAL).
					Where("is_del = ?", 0).
					Updates(buildSeckillUpdates(req, baizeContext.GetUserId(c), &now)).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *ShopSeckillDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	now := time.Now()
	return s.db.WithContext(c).Table(s.tableName).
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Where("is_del = ?", 0).
		Updates(map[string]any{
			"is_del":      1,
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": &now,
		}).Error
}

func (s *ShopSeckillDaoImpl) GetByID(c *gin.Context, id int64) (*models.Seckill, error) {
	var item models.Seckill
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

func (s *ShopSeckillDaoImpl) List(c *gin.Context, req *models.SeckillQuery) (*models.SeckillListData, error) {
	db := s.baseQuery(c)
	if title := strings.TrimSpace(req.Title); title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	if req.ActivityID > 0 {
		db = db.Where("activity_id = ?", req.ActivityID)
	}
	if req.ProductID > 0 {
		db = db.Where("product_id = ?", req.ProductID)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.IsShow != nil {
		db = db.Where("is_show = ?", *req.IsShow)
	}
	if req.IsHot != nil {
		db = db.Where("is_hot = ?", *req.IsHot)
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
	rows := make([]*models.Seckill, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order("sort ASC").Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.SeckillListData{Rows: rows, Total: total}, nil
}

// DeductStock 原子扣减秒杀活动库存。
func (s *ShopSeckillDaoImpl) DeductStock(c *gin.Context, id int64, quantity int64) error {
	if quantity <= 0 {
		return errors.New("扣减库存数量必须大于0")
	}
	result := activityCurrentDB(c, s.db).WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		Where("is_del = ?", 0).
		Where("stock >= ?", quantity).
		Updates(map[string]any{
			"stock":       gorm.Expr("stock - ?", quantity),
			"sales":       gorm.Expr("sales + ?", quantity),
			"update_time": gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("秒杀活动库存不足")
	}
	return nil
}

// RestoreStock 原子回补秒杀活动库存。
func (s *ShopSeckillDaoImpl) RestoreStock(c *gin.Context, id int64, quantity int64) error {
	if quantity <= 0 {
		return errors.New("回补库存数量必须大于0")
	}
	result := activityCurrentDB(c, s.db).WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		Where("is_del = ?", 0).
		Updates(map[string]any{
			"stock":       gorm.Expr("stock + ?", quantity),
			"sales":       gorm.Expr("GREATEST(sales - ?, 0)", quantity),
			"update_time": gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("秒杀活动不存在")
	}
	return nil
}

func (s *ShopSeckillDaoImpl) baseQuery(c *gin.Context) *gorm.DB {
	return s.db.WithContext(c).Table(s.tableName).
		Where("state = ?", commonStatus.NORMAL).
		Where("is_del = ?", 0)
}

func activityCurrentDB(c *gin.Context, db *gorm.DB) *gorm.DB {
	if c == nil {
		return db
	}
	if value, ok := c.Get("db"); ok {
		if tx, ok := value.(*gorm.DB); ok && tx != nil {
			return tx
		}
	}
	return db
}

func (s *ShopSeckillDaoImpl) create(c *gin.Context, req *models.SeckillSet) (*models.Seckill, error) {
	model := buildSeckillModel(req)
	model.ID = snowflake.GenID()
	model.IsDel = 0
	model.State = commonStatus.NORMAL
	model.DeptID = baizeContext.GetDeptId(c)
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ShopSeckillDaoImpl) update(c *gin.Context, req *models.SeckillSet) (*models.Seckill, error) {
	now := time.Now()
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", req.ID).
		Where("state = ?", commonStatus.NORMAL).
		Where("is_del = ?", 0).
		Updates(buildSeckillUpdates(req, baizeContext.GetUserId(c), &now)).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, req.ID)
}

func buildSeckillModel(req *models.SeckillSet) *models.Seckill {
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	return &models.Seckill{
		ActivityID:   req.ActivityID,
		ProductID:    req.ProductID,
		Image:        req.Image,
		Images:       req.Images,
		Title:        req.Title,
		Info:         req.Info,
		Price:        req.Price,
		Cost:         req.Cost,
		OtPrice:      req.OtPrice,
		GiveIntegral: req.GiveIntegral,
		Sort:         req.Sort,
		Stock:        req.Stock,
		Sales:        req.Sales,
		UnitName:     req.UnitName,
		Postage:      req.Postage,
		StartTime:    req.StartTime,
		StopTime:     req.StopTime,
		AddTime:      nowStr,
		Status:       req.Status,
		IsPostage:    req.IsPostage,
		IsHot:        req.IsHot,
		Num:          req.Num,
		IsShow:       req.IsShow,
		TimeID:       req.TimeID,
		TempID:       req.TempID,
		Weight:       req.Weight,
		Volume:       req.Volume,
		Quota:        req.Quota,
		QuotaShow:    req.QuotaShow,
		OnceNum:      req.OnceNum,
		Logistics:    req.Logistics,
		Freight:      req.Freight,
		CustomForm:   req.CustomForm,
		VirtualType:  req.VirtualType,
		IsCommission: req.IsCommission,
	}
}

func buildSeckillUpdates(req *models.SeckillSet, userID int64, now *time.Time) map[string]any {
	return map[string]any{
		"activity_id":   req.ActivityID,
		"product_id":    req.ProductID,
		"image":         req.Image,
		"images":        req.Images,
		"title":         req.Title,
		"info":          req.Info,
		"price":         req.Price,
		"cost":          req.Cost,
		"ot_price":      req.OtPrice,
		"give_integral": req.GiveIntegral,
		"sort":          req.Sort,
		"stock":         req.Stock,
		"sales":         req.Sales,
		"unit_name":     req.UnitName,
		"postage":       req.Postage,
		"start_time":    req.StartTime,
		"stop_time":     req.StopTime,
		"status":        req.Status,
		"is_postage":    req.IsPostage,
		"is_hot":        req.IsHot,
		"num":           req.Num,
		"is_show":       req.IsShow,
		"time_id":       req.TimeID,
		"temp_id":       req.TempID,
		"weight":        req.Weight,
		"volume":        req.Volume,
		"quota":         req.Quota,
		"quota_show":    req.QuotaShow,
		"once_num":      req.OnceNum,
		"logistics":     req.Logistics,
		"freight":       req.Freight,
		"custom_form":   req.CustomForm,
		"virtual_type":  req.VirtualType,
		"is_commission": req.IsCommission,
		"update_by":     userID,
		"update_time":   now,
	}
}
