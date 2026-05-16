package financedaoimpl

import (
	"errors"
	"fmt"
	"time"

	"nova-factory-server/app/business/erp/finance/financedao"
	"nova-factory-server/app/business/erp/finance/financemodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FinanceReceiptItemDaoImpl struct {
	db *gorm.DB
}

func NewFinanceReceiptItemDao(db *gorm.DB) financedao.IFinanceReceiptItemDao {
	return &FinanceReceiptItemDaoImpl{db: db}
}

func (d *FinanceReceiptItemDaoImpl) Create(c *gin.Context, req *financemodels.FinanceReceiptItemUpsert) (*financemodels.FinanceReceiptItem, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := financemodels.FinanceReceiptItemUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_finance_receipt_item").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *FinanceReceiptItemDaoImpl) Update(c *gin.Context, req *financemodels.FinanceReceiptItemUpsert) (*financemodels.FinanceReceiptItem, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.ReceiptID > 0 {
		updates["receipt_id"] = req.ReceiptID
	}
	updates["biz_type"] = req.BizType
	if req.BizID > 0 {
		updates["biz_id"] = req.BizID
	}
	if req.BizNo != "" {
		updates["biz_no"] = req.BizNo
	}
	if req.TotalPrice != 0 {
		updates["total_price"] = req.TotalPrice
	}
	if req.ReceiptedPrice != 0 {
		updates["receipted_price"] = req.ReceiptedPrice
	}
	if req.ReceiptPrice != 0 {
		updates["receipt_price"] = req.ReceiptPrice
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_finance_receipt_item").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *FinanceReceiptItemDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_finance_receipt_item").Where("id IN ?", ids).Where("state = ?", commonStatus.NORMAL).Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *FinanceReceiptItemDaoImpl) GetByID(c *gin.Context, id int64) (*financemodels.FinanceReceiptItem, error) {
	item := new(financemodels.FinanceReceiptItem)
	if err := d.db.WithContext(c).Table("erp_finance_receipt_item").Where("id = ?", id).Where("state = ?", commonStatus.NORMAL).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *FinanceReceiptItemDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinanceReceiptItem, error) {
	if column == "" {
		return nil, nil
	}
	item := new(financemodels.FinanceReceiptItem)
	if err := d.db.WithContext(c).Table("erp_finance_receipt_item").Where(fmt.Sprintf("%s = ?", column), value).Where("state = ?", commonStatus.NORMAL).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *FinanceReceiptItemDaoImpl) ListPage(c *gin.Context, req *financemodels.FinanceReceiptItemQuery) (*financemodels.FinanceReceiptItemListData, error) {
	if req == nil {
		req = new(financemodels.FinanceReceiptItemQuery)
	}
	db := d.db.WithContext(c).Table("erp_finance_receipt_item").Where("state = ?", commonStatus.NORMAL)
	db = applyFinanceReceiptItemFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]financemodels.FinanceReceiptItem, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*financemodels.FinanceReceiptItem, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &financemodels.FinanceReceiptItemListData{Rows: result, Total: total}, nil
}

func (d *FinanceReceiptItemDaoImpl) List(c *gin.Context, req *financemodels.FinanceReceiptItemQuery) (*financemodels.FinanceReceiptItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinanceReceiptItemListData{Rows: result.Rows, Total: result.Total}, nil
}
