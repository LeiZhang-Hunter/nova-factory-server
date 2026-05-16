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

// StockMoveDaoImpl 提供数据访问能力。
type StockMoveDaoImpl struct {
	db *gorm.DB
}

// NewStockMoveDao 创建 DAO。
func NewStockMoveDao(db *gorm.DB) stockdao.IStockMoveDao {
	return &StockMoveDaoImpl{db: db}
}

func (d *StockMoveDaoImpl) Create(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error) {
	model := new(stockmodels.StockMove)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_stock_move").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *StockMoveDaoImpl) Update(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(stockmodels.StockMove)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_stock_move").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *StockMoveDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_stock_move").Where("id IN ?", ids)
	if erpbiz.HasField(new(stockmodels.StockMove), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *StockMoveDaoImpl) GetByID(c *gin.Context, id int64) (*stockmodels.StockMove, error) {
	item := new(stockmodels.StockMove)
	db := d.db.WithContext(c).Table("erp_stock_move").Where("id = ?", id)
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

func (d *StockMoveDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockMove, error) {
	if column == "" {
		return nil, nil
	}
	item := new(stockmodels.StockMove)
	db := d.db.WithContext(c).Table("erp_stock_move").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *StockMoveDaoImpl) ListPage(c *gin.Context, req *stockmodels.StockMoveQuery) (*erpbiz.PageResult[stockmodels.StockMove], error) {
	if req == nil {
		req = new(stockmodels.StockMoveQuery)
	}
	db := d.db.WithContext(c).Table("erp_stock_move")
	if erpbiz.HasField(new(stockmodels.StockMove), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]stockmodels.StockMove, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*stockmodels.StockMove, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[stockmodels.StockMove]{Rows: result, Total: total}, nil
}

func (d *StockMoveDaoImpl) List(c *gin.Context, req *stockmodels.StockMoveQuery) (*stockmodels.StockMoveListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockMoveListData{Rows: result.Rows, Total: result.Total}, nil
}
