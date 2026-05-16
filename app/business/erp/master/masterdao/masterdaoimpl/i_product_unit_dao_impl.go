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

type ProductUnitDaoImpl struct {
	db *gorm.DB
}

func NewProductUnitDao(db *gorm.DB) masterdao.IProductUnitDao {
	return &ProductUnitDaoImpl{db: db}
}

func (d *ProductUnitDaoImpl) Create(c *gin.Context, req *mastermodels.ProductUnitUpsert) (*mastermodels.ProductUnit, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := mastermodels.ProductUnitUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_product_unit").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *ProductUnitDaoImpl) Update(c *gin.Context, req *mastermodels.ProductUnitUpsert) (*mastermodels.ProductUnit, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.Name != "" {
		updates["name"] = req.Name
	}
	updates["status"] = req.Status
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_product_unit").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *ProductUnitDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_product_unit").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *ProductUnitDaoImpl) GetByID(c *gin.Context, id int64) (*mastermodels.ProductUnit, error) {
	item := new(mastermodels.ProductUnit)
	if err := d.db.WithContext(c).Table("erp_product_unit").
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

func (d *ProductUnitDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*mastermodels.ProductUnit, error) {
	if column == "" {
		return nil, nil
	}
	item := new(mastermodels.ProductUnit)
	if err := d.db.WithContext(c).Table("erp_product_unit").
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

func (d *ProductUnitDaoImpl) ListPage(c *gin.Context, req *mastermodels.ProductUnitQuery) (*mastermodels.ProductUnitListData, error) {
	if req == nil {
		req = new(mastermodels.ProductUnitQuery)
	}
	db := d.db.WithContext(c).Table("erp_product_unit").Where("state = ?", commonStatus.NORMAL)
	db = applyProductUnitFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]mastermodels.ProductUnit, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*mastermodels.ProductUnit, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &mastermodels.ProductUnitListData{Rows: result, Total: total}, nil
}

func (d *ProductUnitDaoImpl) List(c *gin.Context, req *mastermodels.ProductUnitQuery) (*mastermodels.ProductUnitListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductUnitListData{Rows: result.Rows, Total: result.Total}, nil
}
