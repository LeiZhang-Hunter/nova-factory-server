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

// ShopHomeModuleDaoImpl 提供首页模块的数据访问能力。
type ShopHomeModuleDaoImpl struct {
	db          *gorm.DB
	itemDao     dao.IShopHomeModuleItemDao
	moduleTable string
}

// NewShopHomeModuleDao 创建首页模块 DAO。
func NewShopHomeModuleDao(ms *gorm.DB, itemDao dao.IShopHomeModuleItemDao) dao.IShopHomeModuleDao {
	return &ShopHomeModuleDaoImpl{
		db:          ms,
		itemDao:     itemDao,
		moduleTable: "shop_home_module",
	}
}

// Set 新增或修改首页模块。
func (s *ShopHomeModuleDaoImpl) Set(c *gin.Context, req *models.HomeModuleSet) (*models.HomeModule, error) {
	if req.ID > 0 {
		if err := s.update(c, req); err != nil {
			return nil, err
		}
		return s.GetByID(c, req.ID)
	}
	return s.create(c, req)
}

// DeleteByIDs 软删除首页模块；存在明细数据时不允许删除。
func (s *ShopHomeModuleDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	hasItems, err := s.itemDao.HasByModuleIDs(c, ids)
	if err != nil {
		return err
	}
	if hasItems {
		return errors.New("存在首页模块明细数据，不能删除")
	}

	now := time.Now()
	userID := baizeContext.GetUserId(c)
	return s.db.WithContext(c).Table(s.moduleTable).
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   userID,
			"update_time": &now,
		}).Error
}

// GetByID 根据主键获取首页模块详情。
func (s *ShopHomeModuleDaoImpl) GetByID(c *gin.Context, id int64) (*models.HomeModule, error) {
	var item models.HomeModule
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

// List 分页查询首页模块列表。
func (s *ShopHomeModuleDaoImpl) List(c *gin.Context, req *models.HomeModuleQuery) (*models.HomeModuleListData, error) {
	db := s.baseQuery(c)
	if req.PageKey != "" {
		db = db.Where("page_key = ?", req.PageKey)
	}
	if req.ModuleType != "" {
		db = db.Where("module_type = ?", req.ModuleType)
	}
	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
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

	rows := make([]*models.HomeModule, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).
		Limit(int(req.Size)).
		Order("sort ASC").
		Order("id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.HomeModuleListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func (s *ShopHomeModuleDaoImpl) baseQuery(c *gin.Context) *gorm.DB {
	return s.db.WithContext(c).Table(s.moduleTable).
		Where("state = ?", commonStatus.NORMAL)
}

func (s *ShopHomeModuleDaoImpl) create(c *gin.Context, req *models.HomeModuleSet) (*models.HomeModule, error) {
	model := &models.HomeModule{
		ID:         snowflake.GenID(),
		PageKey:    req.PageKey,
		ModuleType: req.ModuleType,
		ModuleName: req.ModuleName,
		Title:      req.Title,
		SubTitle:   req.SubTitle,
		SourceType: req.SourceType,
		LimitNum:   req.LimitNum,
		Sort:       req.Sort,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		ShowMore:   req.ShowMore,
		MoreLink:   req.MoreLink,
		StyleJSON:  req.StyleJSON,
		RuleJSON:   req.RuleJSON,
		ExtJSON:    req.ExtJSON,
		Status:     req.Status,
		DeptID:     baizeContext.GetDeptId(c),
		State:      commonStatus.NORMAL,
	}
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := s.db.WithContext(c).Table(s.moduleTable).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ShopHomeModuleDaoImpl) update(c *gin.Context, req *models.HomeModuleSet) error {
	existing, err := s.GetByID(c, req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("首页模块不存在")
	}
	now := time.Now()
	return s.db.WithContext(c).Table(s.moduleTable).
		Where("id = ?", req.ID).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"page_key":    req.PageKey,
			"module_type": req.ModuleType,
			"module_name": req.ModuleName,
			"title":       req.Title,
			"sub_title":   req.SubTitle,
			"source_type": req.SourceType,
			"limit_num":   req.LimitNum,
			"sort":        req.Sort,
			"start_time":  req.StartTime,
			"end_time":    req.EndTime,
			"show_more":   req.ShowMore,
			"more_link":   req.MoreLink,
			"style_json":  req.StyleJSON,
			"rule_json":   req.RuleJSON,
			"ext_json":    req.ExtJSON,
			"status":      req.Status,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": &now,
		}).Error
}
