package stockdaoimpl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StockRecordDaoImpl 提供数据访问能力。
type StockRecordDaoImpl struct {
	db *gorm.DB
}

// NewStockRecordDao 创建 DAO。
func NewStockRecordDao(db *gorm.DB) stockdao.IStockRecordDao {
	return &StockRecordDaoImpl{db: db}
}

func (d *StockRecordDaoImpl) Create(c *gin.Context, req *stockmodels.StockRecordUpsert) (*stockmodels.StockRecord, error) {
	model := new(stockmodels.StockRecord)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_stock_record").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *StockRecordDaoImpl) Update(c *gin.Context, req *stockmodels.StockRecordUpsert) (*stockmodels.StockRecord, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(stockmodels.StockRecord)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_stock_record").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *StockRecordDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_stock_record").Where("id IN ?", ids)
	if erpbiz.HasField(new(stockmodels.StockRecord), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *StockRecordDaoImpl) GetByID(c *gin.Context, id int64) (*stockmodels.StockRecord, error) {
	item := new(stockmodels.StockRecord)
	db := d.db.WithContext(c).Table("erp_stock_record").Where("id = ?", id)
	if erpbiz.HasField(item, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.First(item).Error; err != nil {
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
	db := d.db.WithContext(c).Table("erp_stock_record").Where(fmt.Sprintf("%s = ?", column), value)
	if erpbiz.HasField(item, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *StockRecordDaoImpl) ListPage(c *gin.Context, req *stockmodels.StockRecordQuery) (*erpbiz.PageResult[stockmodels.StockRecord], error) {
	if req == nil {
		req = new(stockmodels.StockRecordQuery)
	}
	db := d.db.WithContext(c).Table("erp_stock_record")
	if erpbiz.HasField(new(stockmodels.StockRecord), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]stockmodels.StockRecord, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*stockmodels.StockRecord, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[stockmodels.StockRecord]{Rows: result, Total: total}, nil
}

func (d *StockRecordDaoImpl) List(c *gin.Context, req *stockmodels.StockRecordQuery) (*stockmodels.StockRecordListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockRecordListData{Rows: result.Rows, Total: result.Total}, nil
}
