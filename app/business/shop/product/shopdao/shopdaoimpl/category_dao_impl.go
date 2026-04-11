package shopdaoimpl

import (
	"errors"
	"fmt"
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopCategoryDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewShopCategoryDao(ms *gorm.DB) shopdao.IShopCategoryDao {
	return &ShopCategoryDaoImpl{
		db:        ms,
		tableName: "shop_category",
	}
}

func (s *ShopCategoryDaoImpl) Create(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error) {
	model := &shopmodels.Category{
		ID:           snowflake.GenID(),
		ParentID:     req.ParentID,
		CategoryName: req.CategoryName,
		CategoryCode: req.CategoryCode,
		Sort:         req.Sort,
		Status:       req.Status,
		Depth:        1,
		AncestorPath: "/",
	}
	model.SetCreateBy(baizeContext.GetUserId(c))
	model.DeptID = baizeContext.GetDeptId(c)
	if req.ParentID > 0 {
		var parent shopmodels.Category
		if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", req.ParentID).First(&parent).Error; err != nil {
			return nil, err
		}
		model.Depth = parent.Depth + 1
		model.AncestorPath = parent.AncestorPath + parentIDPath(parent.ID)
	}
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ShopCategoryDaoImpl) Update(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error) {
	updates := &shopmodels.Category{
		ID:           req.ID,
		CategoryName: req.CategoryName,
		CategoryCode: req.CategoryCode,
		Sort:         req.Sort,
		Status:       req.Status,
	}
	updates.DeptID = baizeContext.GetDeptId(c)
	updates.SetUpdateBy(baizeContext.GetUserId(c))
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", req.ID).
		Select("category_name", "category_code", "sort", "status", "update_by", "update_time").
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, int64(req.ID))
}

func (s *ShopCategoryDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.db.WithContext(c).Table(s.tableName).Where("id IN ?", ids).Update("state", commonStatus.DELETE).Error
}

func (s *ShopCategoryDaoImpl) GetByID(c *gin.Context, id int64) (*shopmodels.Category, error) {
	var item shopmodels.Category
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (s *ShopCategoryDaoImpl) List(c *gin.Context, req *shopmodels.CategoryQuery) (*shopmodels.CategoryListData, error) {
	db := s.db.WithContext(c).Table(s.tableName).Where("state = 0")
	if req.CategoryName != "" {
		db = db.Where("category_name LIKE ?", "%"+req.CategoryName+"%")
	}
	if req.CategoryCode != "" {
		db = db.Where("category_code = ?", req.CategoryCode)
	}
	if req.Status != nil {
		db = db.Where("status = ?", req.Status)
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*shopmodels.Category, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return &shopmodels.CategoryListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func parentIDPath(id int64) string {
	return fmt.Sprintf("%d/", id)
}
