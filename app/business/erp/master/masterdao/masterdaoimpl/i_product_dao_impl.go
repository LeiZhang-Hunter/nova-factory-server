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

type ProductDaoImpl struct {
	db *gorm.DB
}

func NewProductDao(db *gorm.DB) masterdao.IProductDao {
	return &ProductDaoImpl{db: db}
}

func (d *ProductDaoImpl) Create(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := mastermodels.ProductUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_product").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *ProductDaoImpl) Update(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.BarCode != "" {
		updates["bar_code"] = req.BarCode
	}
	if req.CategoryId > 0 {
		updates["category_id"] = req.CategoryId
	}
	if req.UnitId > 0 {
		updates["unit_id"] = req.UnitId
	}
	updates["status"] = req.Status
	if req.Standard != "" {
		updates["standard"] = req.Standard
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	updates["expiry_day"] = req.ExpiryDay
	if req.Weight != 0 {
		updates["weight"] = req.Weight
	}
	if req.PurchasePrice != 0 {
		updates["purchase_price"] = req.PurchasePrice
	}
	if req.SalePrice != 0 {
		updates["sale_price"] = req.SalePrice
	}
	if req.MinPrice != 0 {
		updates["min_price"] = req.MinPrice
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_product").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *ProductDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_product").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *ProductDaoImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Product, error) {
	item := new(mastermodels.Product)
	if err := d.db.WithContext(c).Table("erp_product").
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

func (d *ProductDaoImpl) GetByIDs(c *gin.Context, ids []int64) ([]*mastermodels.Product, error) {
	if len(ids) == 0 {
		return make([]*mastermodels.Product, 0), nil
	}
	rows := make([]*mastermodels.Product, 0, len(ids))
	if err := d.db.WithContext(c).Table("erp_product").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *ProductDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Product, error) {
	if column == "" {
		return nil, nil
	}
	item := new(mastermodels.Product)
	if err := d.db.WithContext(c).Table("erp_product").
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

func (d *ProductDaoImpl) ListPage(c *gin.Context, req *mastermodels.ProductQuery) (*mastermodels.ProductListData, error) {
	if req == nil {
		req = new(mastermodels.ProductQuery)
	}
	db := d.db.WithContext(c).Table("erp_product").Where("state = ?", commonStatus.NORMAL)
	db = applyProductFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]mastermodels.Product, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*mastermodels.Product, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &mastermodels.ProductListData{Rows: result, Total: total}, nil
}

func (d *ProductDaoImpl) List(c *gin.Context, req *mastermodels.ProductQuery) (*mastermodels.ProductListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductListData{Rows: result.Rows, Total: result.Total}, nil
}

func (d *ProductDaoImpl) UpsertByID(c *gin.Context, id int64, updates map[string]any) error {
	existing := new(mastermodels.Product)
	err := d.db.WithContext(c).Table("erp_product").Where("id = ?", id).First(existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		updates["id"] = id
		updates["state"] = commonStatus.NORMAL
		return d.db.WithContext(c).Table("erp_product").Create(updates).Error
	}
	if err != nil {
		return err
	}
	return d.db.WithContext(c).Table("erp_product").Where("id = ?", id).Updates(updates).Error
}
