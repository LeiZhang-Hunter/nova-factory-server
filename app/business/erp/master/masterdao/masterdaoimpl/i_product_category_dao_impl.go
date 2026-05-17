package masterdaoimpl

import (
	"errors"
	"fmt"
	"time"

	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductCategoryDaoImpl struct {
	db *gorm.DB
}

func NewProductCategoryDao(db *gorm.DB) masterdao.IProductCategoryDao {
	return &ProductCategoryDaoImpl{db: db}
}

func (d *ProductCategoryDaoImpl) Create(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := mastermodels.ProductCategoryUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_product_category").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *ProductCategoryDaoImpl) Update(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.ParentID > 0 {
		updates["parent_id"] = req.ParentID
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	updates["sort"] = req.Sort
	updates["status"] = req.Status
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_product_category").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *ProductCategoryDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_product_category").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *ProductCategoryDaoImpl) GetByID(c *gin.Context, id int64) (*mastermodels.ProductCategory, error) {
	item := new(mastermodels.ProductCategory)
	if err := d.db.WithContext(c).Table("erp_product_category").
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *ProductCategoryDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*mastermodels.ProductCategory, error) {
	if column == "" {
		return nil, nil
	}
	item := new(mastermodels.ProductCategory)
	if err := d.db.WithContext(c).Table("erp_product_category").
		Where(fmt.Sprintf("%s = ?", column), value).
		Where("state = ?", commonStatus.NORMAL).
		First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *ProductCategoryDaoImpl) GetByIDs(c *gin.Context, ids []int64) ([]*mastermodels.ProductCategory, error) {
	if len(ids) == 0 {
		return []*mastermodels.ProductCategory{}, nil
	}
	rows := make([]mastermodels.ProductCategory, 0)
	if err := d.db.WithContext(c).Table("erp_product_category").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*mastermodels.ProductCategory, 0, len(rows))
	for i := range rows {
		result = append(result, &rows[i])
	}
	return result, nil
}

func (d *ProductCategoryDaoImpl) ListPage(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*mastermodels.ProductCategoryListData, error) {
	if req == nil {
		req = new(mastermodels.ProductCategoryQuery)
	}
	db := d.db.WithContext(c).Table("erp_product_category").Where("state = ?", commonStatus.NORMAL)
	db = applyProductCategoryFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]mastermodels.ProductCategory, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*mastermodels.ProductCategory, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &mastermodels.ProductCategoryListData{Rows: result, Total: total}, nil
}

func (d *ProductCategoryDaoImpl) List(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*mastermodels.ProductCategoryListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductCategoryListData{Rows: result.Rows, Total: result.Total}, nil
}
