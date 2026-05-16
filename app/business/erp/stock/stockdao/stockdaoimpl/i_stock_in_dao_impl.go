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

// StockInDaoImpl 提供数据访问能力。
type StockInDaoImpl struct {
	db *gorm.DB
}

// NewStockInDao 创建 DAO。
func NewStockInDao(db *gorm.DB) stockdao.IStockInDao {
	return &StockInDaoImpl{db: db}
}

func (d *StockInDaoImpl) Create(c *gin.Context, req *stockmodels.StockInUpsert) (*stockmodels.StockIn, error) {
	model := new(stockmodels.StockIn)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_stock_in").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *StockInDaoImpl) Update(c *gin.Context, req *stockmodels.StockInUpsert) (*stockmodels.StockIn, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(stockmodels.StockIn)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_stock_in").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *StockInDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_stock_in").Where("id IN ?", ids)
	if erpbiz.HasField(new(stockmodels.StockIn), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *StockInDaoImpl) GetByID(c *gin.Context, id int64) (*stockmodels.StockIn, error) {
	item := new(stockmodels.StockIn)
	db := d.db.WithContext(c).Table("erp_stock_in").Where("id = ?", id)
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

func (d *StockInDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockIn, error) {
	if column == "" {
		return nil, nil
	}
	item := new(stockmodels.StockIn)
	db := d.db.WithContext(c).Table("erp_stock_in").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *StockInDaoImpl) ListPage(c *gin.Context, req *stockmodels.StockInQuery) (*erpbiz.PageResult[stockmodels.StockIn], error) {
	if req == nil {
		req = new(stockmodels.StockInQuery)
	}
	db := d.db.WithContext(c).Table("erp_stock_in")
	if erpbiz.HasField(new(stockmodels.StockIn), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]stockmodels.StockIn, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*stockmodels.StockIn, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[stockmodels.StockIn]{Rows: result, Total: total}, nil
}

func (d *StockInDaoImpl) List(c *gin.Context, req *stockmodels.StockInQuery) (*stockmodels.StockInListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockInListData{Rows: result.Rows, Total: result.Total}, nil
}
