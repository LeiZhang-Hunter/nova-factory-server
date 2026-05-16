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

type StockCheckDaoImpl struct {
	db *gorm.DB
}

func NewStockCheckDao(db *gorm.DB) stockdao.IStockCheckDao {
	return &StockCheckDaoImpl{db: db}
}

func (d *StockCheckDaoImpl) Create(c *gin.Context, req *stockmodels.StockCheckUpsert) (*stockmodels.StockCheck, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := stockmodels.StockCheckUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_stock_check").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *StockCheckDaoImpl) Update(c *gin.Context, req *stockmodels.StockCheckUpsert) (*stockmodels.StockCheck, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.No != "" {
		updates["no"] = req.No
	}
	if req.CheckTime != "" {
		if t, err := parseTime(req.CheckTime); err == nil {
			updates["check_time"] = t
		}
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
	db := d.db.WithContext(c).Table("erp_stock_check").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *StockCheckDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_stock_check").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *StockCheckDaoImpl) GetByID(c *gin.Context, id int64) (*stockmodels.StockCheck, error) {
	item := new(stockmodels.StockCheck)
	if err := d.db.WithContext(c).Table("erp_stock_check").
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

func (d *StockCheckDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockCheck, error) {
	if column == "" {
		return nil, nil
	}
	item := new(stockmodels.StockCheck)
	if err := d.db.WithContext(c).Table("erp_stock_check").
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

func (d *StockCheckDaoImpl) ListPage(c *gin.Context, req *stockmodels.StockCheckQuery) (*stockmodels.StockCheckListData, error) {
	if req == nil {
		req = new(stockmodels.StockCheckQuery)
	}
	db := d.db.WithContext(c).Table("erp_stock_check").Where("state = ?", commonStatus.NORMAL)
	db = applyStockCheckFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]stockmodels.StockCheck, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*stockmodels.StockCheck, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &stockmodels.StockCheckListData{Rows: result, Total: total}, nil
}

func (d *StockCheckDaoImpl) List(c *gin.Context, req *stockmodels.StockCheckQuery) (*stockmodels.StockCheckListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockCheckListData{Rows: result.Rows, Total: result.Total}, nil
}
