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

// StockMoveItemDaoImpl 提供数据访问能力。
type StockMoveItemDaoImpl struct {
	db *gorm.DB
}

// NewStockMoveItemDao 创建 DAO。
func NewStockMoveItemDao(db *gorm.DB) stockdao.IStockMoveItemDao {
	return &StockMoveItemDaoImpl{db: db}
}

func (d *StockMoveItemDaoImpl) Create(c *gin.Context, req *stockmodels.StockMoveItemUpsert) (*stockmodels.StockMoveItem, error) {
	model := new(stockmodels.StockMoveItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_stock_move_item").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *StockMoveItemDaoImpl) Update(c *gin.Context, req *stockmodels.StockMoveItemUpsert) (*stockmodels.StockMoveItem, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(stockmodels.StockMoveItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_stock_move_item").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *StockMoveItemDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_stock_move_item").Where("id IN ?", ids)
	if erpbiz.HasField(new(stockmodels.StockMoveItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *StockMoveItemDaoImpl) GetByID(c *gin.Context, id int64) (*stockmodels.StockMoveItem, error) {
	item := new(stockmodels.StockMoveItem)
	db := d.db.WithContext(c).Table("erp_stock_move_item").Where("id = ?", id)
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

func (d *StockMoveItemDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockMoveItem, error) {
	if column == "" {
		return nil, nil
	}
	item := new(stockmodels.StockMoveItem)
	db := d.db.WithContext(c).Table("erp_stock_move_item").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *StockMoveItemDaoImpl) ListPage(c *gin.Context, req *stockmodels.StockMoveItemQuery) (*erpbiz.PageResult[stockmodels.StockMoveItem], error) {
	if req == nil {
		req = new(stockmodels.StockMoveItemQuery)
	}
	db := d.db.WithContext(c).Table("erp_stock_move_item")
	if erpbiz.HasField(new(stockmodels.StockMoveItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]stockmodels.StockMoveItem, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*stockmodels.StockMoveItem, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[stockmodels.StockMoveItem]{Rows: result, Total: total}, nil
}

func (d *StockMoveItemDaoImpl) List(c *gin.Context, req *stockmodels.StockMoveItemQuery) (*stockmodels.StockMoveItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockMoveItemListData{Rows: result.Rows, Total: result.Total}, nil
}
