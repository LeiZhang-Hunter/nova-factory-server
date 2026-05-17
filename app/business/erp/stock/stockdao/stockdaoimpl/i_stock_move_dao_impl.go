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

type StockMoveDaoImpl struct {
	db *gorm.DB
}

func NewStockMoveDao(db *gorm.DB) stockdao.IStockMoveDao {
	return &StockMoveDaoImpl{db: db}
}

func (d *StockMoveDaoImpl) Create(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := stockmodels.StockMoveUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_stock_move").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *StockMoveDaoImpl) Update(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.MoveTime != "" {
		if parsed, err := time.ParseInLocation("2006-01-02 15:04:05", req.MoveTime, time.Local); err == nil {
			updates["move_time"] = &parsed
		}
	}
	if req.No != "" {
		updates["no"] = req.No
	}
	if req.TotalCount != 0 {
		updates["total_count"] = req.TotalCount
	}
	if req.TotalPrice != 0 {
		updates["total_price"] = req.TotalPrice
	}
	updates["status"] = req.Status
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if req.FileURL != "" {
		updates["file_url"] = req.FileURL
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_stock_move").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *StockMoveDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_stock_move").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *StockMoveDaoImpl) GetByID(c *gin.Context, id int64) (*stockmodels.StockMove, error) {
	item := new(stockmodels.StockMove)
	if err := d.db.WithContext(c).Table("erp_stock_move").
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

func (d *StockMoveDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockMove, error) {
	if column == "" {
		return nil, nil
	}
	item := new(stockmodels.StockMove)
	if err := d.db.WithContext(c).Table("erp_stock_move").
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

func (d *StockMoveDaoImpl) ListPage(c *gin.Context, req *stockmodels.StockMoveQuery) (*stockmodels.StockMoveListData, error) {
	if req == nil {
		req = new(stockmodels.StockMoveQuery)
	}
	db := d.db.WithContext(c).Table("erp_stock_move").Where("state = ?", commonStatus.NORMAL)
	db = applyStockMoveFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]stockmodels.StockMove, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*stockmodels.StockMove, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &stockmodels.StockMoveListData{Rows: result, Total: total}, nil
}

func (d *StockMoveDaoImpl) List(c *gin.Context, req *stockmodels.StockMoveQuery) (*stockmodels.StockMoveListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockMoveListData{Rows: result.Rows, Total: result.Total}, nil
}
