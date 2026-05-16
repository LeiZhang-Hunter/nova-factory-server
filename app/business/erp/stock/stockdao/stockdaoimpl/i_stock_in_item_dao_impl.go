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

type StockInItemDaoImpl struct {
	db *gorm.DB
}

func NewStockInItemDao(db *gorm.DB) stockdao.IStockInItemDao {
	return &StockInItemDaoImpl{db: db}
}

func (d *StockInItemDaoImpl) Create(c *gin.Context, req *stockmodels.StockInItemUpsert) (*stockmodels.StockInItem, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := stockmodels.StockInItemUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_stock_in_item").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *StockInItemDaoImpl) Update(c *gin.Context, req *stockmodels.StockInItemUpsert) (*stockmodels.StockInItem, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.InID > 0 {
		updates["in_id"] = req.InID
	}
	if req.WarehouseID > 0 {
		updates["warehouse_id"] = req.WarehouseID
	}
	if req.ProductID > 0 {
		updates["product_id"] = req.ProductID
	}
	if req.ProductUnitID > 0 {
		updates["product_unit_id"] = req.ProductUnitID
	}
	if req.ProductPrice != 0 {
		updates["product_price"] = req.ProductPrice
	}
	if req.Count != 0 {
		updates["count"] = req.Count
	}
	if req.TotalPrice != 0 {
		updates["total_price"] = req.TotalPrice
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_stock_in_item").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *StockInItemDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_stock_in_item").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *StockInItemDaoImpl) GetByID(c *gin.Context, id int64) (*stockmodels.StockInItem, error) {
	item := new(stockmodels.StockInItem)
	if err := d.db.WithContext(c).Table("erp_stock_in_item").
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

func (d *StockInItemDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockInItem, error) {
	if column == "" {
		return nil, nil
	}
	item := new(stockmodels.StockInItem)
	if err := d.db.WithContext(c).Table("erp_stock_in_item").
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

func (d *StockInItemDaoImpl) ListPage(c *gin.Context, req *stockmodels.StockInItemQuery) (*stockmodels.StockInItemListData, error) {
	if req == nil {
		req = new(stockmodels.StockInItemQuery)
	}
	db := d.db.WithContext(c).Table("erp_stock_in_item").Where("state = ?", commonStatus.NORMAL)
	db = applyStockInItemFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]stockmodels.StockInItem, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*stockmodels.StockInItem, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &stockmodels.StockInItemListData{Rows: result, Total: total}, nil
}

func (d *StockInItemDaoImpl) List(c *gin.Context, req *stockmodels.StockInItemQuery) (*stockmodels.StockInItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockInItemListData{Rows: result.Rows, Total: result.Total}, nil
}
