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

// StockOutItemDaoImpl 提供数据访问能力。
type StockOutItemDaoImpl struct {
	db *gorm.DB
}

// NewStockOutItemDao 创建 DAO。
func NewStockOutItemDao(db *gorm.DB) stockdao.IStockOutItemDao {
	return &StockOutItemDaoImpl{db: db}
}

func (d *StockOutItemDaoImpl) Create(c *gin.Context, req *stockmodels.StockOutItemUpsert) (*stockmodels.StockOutItem, error) {
	model := new(stockmodels.StockOutItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_stock_out_item").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *StockOutItemDaoImpl) Update(c *gin.Context, req *stockmodels.StockOutItemUpsert) (*stockmodels.StockOutItem, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(stockmodels.StockOutItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_stock_out_item").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *StockOutItemDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_stock_out_item").Where("id IN ?", ids)
	if erpbiz.HasField(new(stockmodels.StockOutItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *StockOutItemDaoImpl) GetByID(c *gin.Context, id int64) (*stockmodels.StockOutItem, error) {
	item := new(stockmodels.StockOutItem)
	db := d.db.WithContext(c).Table("erp_stock_out_item").Where("id = ?", id)
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

func (d *StockOutItemDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockOutItem, error) {
	if column == "" {
		return nil, nil
	}
	item := new(stockmodels.StockOutItem)
	db := d.db.WithContext(c).Table("erp_stock_out_item").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *StockOutItemDaoImpl) ListPage(c *gin.Context, req *stockmodels.StockOutItemQuery) (*erpbiz.PageResult[stockmodels.StockOutItem], error) {
	if req == nil {
		req = new(stockmodels.StockOutItemQuery)
	}
	db := d.db.WithContext(c).Table("erp_stock_out_item")
	if erpbiz.HasField(new(stockmodels.StockOutItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]stockmodels.StockOutItem, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*stockmodels.StockOutItem, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[stockmodels.StockOutItem]{Rows: result, Total: total}, nil
}

func (d *StockOutItemDaoImpl) List(c *gin.Context, req *stockmodels.StockOutItemQuery) (*stockmodels.StockOutItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockOutItemListData{Rows: result.Rows, Total: result.Total}, nil
}
