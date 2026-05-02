package impl

import (
	"errors"
	"time"

	"nova-factory-server/app/business/shop/home/dao"
	"nova-factory-server/app/business/shop/home/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ShopHomeModuleItemDaoImpl 提供首页模块明细的数据访问能力。
type ShopHomeModuleItemDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewShopHomeModuleItemDao 创建首页模块明细 DAO。
func NewShopHomeModuleItemDao(ms *gorm.DB) dao.IShopHomeModuleItemDao {
	return &ShopHomeModuleItemDaoImpl{
		db:        ms,
		tableName: "shop_home_module_item",
	}
}

// Set 新增或修改首页模块明细。
func (s *ShopHomeModuleItemDaoImpl) Set(c *gin.Context, req *models.HomeModuleItemSet) (*models.HomeModuleItem, error) {
	if req.ID > 0 {
		if err := s.update(c, req); err != nil {
			return nil, err
		}
		return s.GetByID(c, req.ID)
	}
	return s.create(c, req)
}

// DeleteByIDs 软删除首页模块明细。
func (s *ShopHomeModuleItemDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
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

// GetByID 根据主键获取首页模块明细详情。
func (s *ShopHomeModuleItemDaoImpl) GetByID(c *gin.Context, id int64) (*models.HomeModuleItem, error) {
	var item models.HomeModuleItem
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

// HasByModuleIDs 判断指定模块下是否存在有效的首页模块明细。
func (s *ShopHomeModuleItemDaoImpl) HasByModuleIDs(c *gin.Context, moduleIDs []int64) (bool, error) {
	if len(moduleIDs) == 0 {
		return false, nil
	}
	var total int64
	if err := s.baseQuery(c).
		Where("module_id IN ?", moduleIDs).
		Count(&total).Error; err != nil {
		return false, err
	}
	return total > 0, nil
}

// List 分页查询首页模块明细列表。
func (s *ShopHomeModuleItemDaoImpl) List(c *gin.Context, req *models.HomeModuleItemQuery) (*models.HomeModuleItemListData, error) {
	db := s.baseQuery(c)
	if req.ModuleID > 0 {
		db = db.Where("module_id = ?", req.ModuleID)
	}
	if req.BusinessType != "" {
		db = db.Where("business_type = ?", req.BusinessType)
	}
	if req.ItemName != "" {
		db = db.Where("item_name LIKE ?", "%"+req.ItemName+"%")
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

	rows := make([]*models.HomeModuleItem, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).
		Limit(int(req.Size)).
		Order("sort ASC").
		Order("id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.HomeModuleItemListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func (s *ShopHomeModuleItemDaoImpl) baseQuery(c *gin.Context) *gorm.DB {
	return s.db.WithContext(c).Table(s.tableName).
		Where("state = ?", commonStatus.NORMAL)
}

func (s *ShopHomeModuleItemDaoImpl) create(c *gin.Context, req *models.HomeModuleItemSet) (*models.HomeModuleItem, error) {
	model := &models.HomeModuleItem{
		ID:           snowflake.GenID(),
		ModuleID:     req.ModuleID,
		BusinessType: req.BusinessType,
		LinkID:       req.LinkID,
		ItemName:     req.ItemName,
		ItemSubTitle: req.ItemSubTitle,
		ItemImage:    req.ItemImage,
		Sort:         req.Sort,
		Status:       req.Status,
		ExtJSON:      req.ExtJSON,
		DeptID:       baizeContext.GetDeptId(c),
		State:        commonStatus.NORMAL,
	}
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ShopHomeModuleItemDaoImpl) update(c *gin.Context, req *models.HomeModuleItemSet) error {
	existing, err := s.GetByID(c, req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("首页模块明细不存在")
	}
	now := time.Now()
	return s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", req.ID).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"module_id":      req.ModuleID,
			"business_type":  req.BusinessType,
			"link_id":        req.LinkID,
			"item_name":      req.ItemName,
			"item_sub_title": req.ItemSubTitle,
			"item_image":     req.ItemImage,
			"sort":           req.Sort,
			"status":         req.Status,
			"ext_json":       req.ExtJSON,
			"update_by":      baizeContext.GetUserId(c),
			"update_time":    &now,
		}).Error
}
