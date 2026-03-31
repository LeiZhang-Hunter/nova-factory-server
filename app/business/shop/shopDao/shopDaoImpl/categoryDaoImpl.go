package shopDaoImpl

import (
	"fmt"
	"nova-factory-server/app/business/shop/shopDao"
	"nova-factory-server/app/business/shop/shopModels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopCategoryDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewShopCategoryDao(ms *gorm.DB) shopDao.IShopCategoryDao {
	return &ShopCategoryDaoImpl{
		db:        ms,
		tableName: "shop_category",
	}
}

func (s *ShopCategoryDaoImpl) Create(c *gin.Context, req *shopModels.CategoryUpsert) (*shopModels.Category, error) {
	model := &shopModels.Category{
		ParentID:     req.ParentID,
		CategoryName: req.CategoryName,
		CategoryCode: req.CategoryCode,
		Sort:         req.Sort,
		Status:       req.Status,
		Depth:        1,
		AncestorPath: "/",
	}
	if req.ParentID > 0 {
		var parent shopModels.Category
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

func (s *ShopCategoryDaoImpl) Update(c *gin.Context, req *shopModels.CategoryUpsert) (*shopModels.Category, error) {
	updates := map[string]interface{}{
		"category_name": req.CategoryName,
		"category_code": req.CategoryCode,
		"sort":          req.Sort,
		"status":        req.Status,
	}
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, int64(req.ID))
}

func (s *ShopCategoryDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.db.WithContext(c).Table(s.tableName).Where("id IN ?", ids).Update("is_deleted", 1).Error
}

func (s *ShopCategoryDaoImpl) GetByID(c *gin.Context, id int64) (*shopModels.Category, error) {
	var item shopModels.Category
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ShopCategoryDaoImpl) List(c *gin.Context, req *shopModels.CategoryQuery) (*shopModels.CategoryListData, error) {
	db := s.db.WithContext(c).Table(s.tableName).Where("is_deleted = 0")
	if req.CategoryName != "" {
		db = db.Where("category_name LIKE ?", "%"+req.CategoryName+"%")
	}
	if req.CategoryCode != "" {
		db = db.Where("category_code = ?", req.CategoryCode)
	}
	if req.Status == 0 || req.Status == 1 {
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
	rows := make([]*shopModels.Category, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return &shopModels.CategoryListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func parentIDPath(id uint64) string {
	return fmt.Sprintf("%d/", id)
}
