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

type StockRecordDaoImpl struct {
	db *gorm.DB
}

func NewStockRecordDao(db *gorm.DB) stockdao.IStockRecordDao {
	return &StockRecordDaoImpl{db: db}
}

func (d *StockRecordDaoImpl) Create(c *gin.Context, req *stockmodels.StockRecordUpsert) (*stockmodels.StockRecord, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := stockmodels.StockRecordUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_stock_record").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *StockRecordDaoImpl) Update(c *gin.Context, req *stockmodels.StockRecordUpsert) (*stockmodels.StockRecord, error) {
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
	if req.TotalCount != 0 {
		updates["total_count"] = req.TotalCount
	}
	updates["biz_type"] = req.BizType
	if req.BizID > 0 {
		updates["biz_id"] = req.BizID
	}
	if req.BizItemId > 0 {
		updates["biz_item_id"] = req.BizItemId
	}
	if req.BizNo != "" {
		updates["biz_no"] = req.BizNo
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_stock_record").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *StockRecordDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_stock_record").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *StockRecordDaoImpl) GetByID(c *gin.Context, id int64) (*stockmodels.StockRecord, error) {
	item := new(stockmodels.StockRecord)
	if err := d.db.WithContext(c).Table("erp_stock_record").
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

func (d *StockRecordDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockRecord, error) {
	if column == "" {
		return nil, nil
	}
	item := new(stockmodels.StockRecord)
	if err := d.db.WithContext(c).Table("erp_stock_record").
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

func (d *StockRecordDaoImpl) ListPage(c *gin.Context, req *stockmodels.StockRecordQuery) (*stockmodels.StockRecordListData, error) {
	if req == nil {
		req = new(stockmodels.StockRecordQuery)
	}
	db := d.db.WithContext(c).Table("erp_stock_record").Where("state = ?", commonStatus.NORMAL)
	db = applyStockRecordFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]stockmodels.StockRecord, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*stockmodels.StockRecord, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &stockmodels.StockRecordListData{Rows: result, Total: total}, nil
}

func (d *StockRecordDaoImpl) List(c *gin.Context, req *stockmodels.StockRecordQuery) (*stockmodels.StockRecordListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockRecordListData{Rows: result.Rows, Total: result.Total}, nil
}
