package stockdaoimpl

import (
	"errors"
	"fmt"
	"time"

	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StockDaoImpl struct {
	db *gorm.DB
}

func NewStockDao(db *gorm.DB) stockdao.IStockDao {
	return &StockDaoImpl{db: db}
}

func (d *StockDaoImpl) Create(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := stockmodels.StockUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_stock").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *StockDaoImpl) Update(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.ProductID > 0 {
		updates["product_id"] = req.ProductID
	}
	if req.WarehouseID > 0 {
		updates["warehouse_id"] = req.WarehouseID
	}
	if req.Count != 0 {
		updates["count"] = req.Count
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_stock").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *StockDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_stock").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *StockDaoImpl) GetByID(c *gin.Context, id int64) (*stockmodels.Stock, error) {
	item := new(stockmodels.Stock)
	if err := d.db.WithContext(c).Table("erp_stock").
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

func (d *StockDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*stockmodels.Stock, error) {
	if column == "" {
		return nil, nil
	}
	item := new(stockmodels.Stock)
	if err := d.db.WithContext(c).Table("erp_stock").
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

func (d *StockDaoImpl) ListPage(c *gin.Context, req *stockmodels.StockQuery) (*stockmodels.StockListData, error) {
	if req == nil {
		req = new(stockmodels.StockQuery)
	}
	db := d.db.WithContext(c).Table("erp_stock").Where("state = ?", commonStatus.NORMAL)
	db = applyStockFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]stockmodels.Stock, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*stockmodels.Stock, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &stockmodels.StockListData{Rows: result, Total: total}, nil
}

func (d *StockDaoImpl) List(c *gin.Context, req *stockmodels.StockQuery) (*stockmodels.StockListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockListData{Rows: result.Rows, Total: result.Total}, nil
}

func (d *StockDaoImpl) UpdateStockByProductIDWithDB(db *gorm.DB, productID int64, count float64) error {
	return db.Table("erp_stock").
		Where("product_id = ?", productID).
		Update("count", count).Error
}

func (d *StockDaoImpl) UpdateStockByID(c *gin.Context, id int64, count float64) error {
	return d.db.WithContext(c).Table("erp_stock").
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		Update("count", count).Error
}

func (d *StockDaoImpl) UpsertByID(c *gin.Context, id int64, updates map[string]any) error {
	existing := new(stockmodels.Stock)
	err := d.db.WithContext(c).Table("erp_stock").Where("id = ?", id).First(existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		updates["id"] = id
		updates["state"] = commonStatus.NORMAL
		return d.db.WithContext(c).Table("erp_stock").Create(updates).Error
	}
	if err != nil {
		return err
	}
	return d.db.WithContext(c).Table("erp_stock").Where("id = ?", id).Updates(updates).Error
}
