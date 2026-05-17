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

type FinanceReceiptDaoImpl struct {
	db *gorm.DB
}

func NewFinanceReceiptDao(db *gorm.DB) financedao.IFinanceReceiptDao {
	return &FinanceReceiptDaoImpl{db: db}
}

func (d *FinanceReceiptDaoImpl) Create(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := financemodels.FinanceReceiptUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_finance_receipt").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *FinanceReceiptDaoImpl) Update(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.ReceiptTime != "" {
		if parsed, err := time.ParseInLocation("2006-01-02 15:04:05", req.ReceiptTime, time.Local); err == nil {
			updates["receipt_time"] = &parsed
		}
	}
	if req.No != "" {
		updates["no"] = req.No
	}
	updates["status"] = req.Status
	if req.FinanceUserID > 0 {
		updates["finance_user_id"] = req.FinanceUserID
	}
	if req.CustomerID > 0 {
		updates["customer_id"] = req.CustomerID
	}
	if req.AccountID > 0 {
		updates["account_id"] = req.AccountID
	}
	if req.TotalPrice != 0 {
		updates["total_price"] = req.TotalPrice
	}
	if req.DiscountPrice != 0 {
		updates["discount_price"] = req.DiscountPrice
	}
	if req.ReceiptPrice != 0 {
		updates["receipt_price"] = req.ReceiptPrice
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_finance_receipt").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *FinanceReceiptDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_finance_receipt").Where("id IN ?", ids).Where("state = ?", commonStatus.NORMAL).Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *FinanceReceiptDaoImpl) GetByID(c *gin.Context, id int64) (*financemodels.FinanceReceipt, error) {
	item := new(financemodels.FinanceReceipt)
	if err := d.db.WithContext(c).Table("erp_finance_receipt").Where("id = ?", id).Where("state = ?", commonStatus.NORMAL).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *FinanceReceiptDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinanceReceipt, error) {
	if column == "" {
		return nil, nil
	}
	item := new(financemodels.FinanceReceipt)
	if err := d.db.WithContext(c).Table("erp_finance_receipt").Where(fmt.Sprintf("%s = ?", column), value).Where("state = ?", commonStatus.NORMAL).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *FinanceReceiptDaoImpl) ListPage(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*financemodels.FinanceReceiptListData, error) {
	if req == nil {
		req = new(financemodels.FinanceReceiptQuery)
	}
	db := d.db.WithContext(c).Table("erp_finance_receipt").Where("state = ?", commonStatus.NORMAL)
	db = applyFinanceReceiptFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]financemodels.FinanceReceipt, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*financemodels.FinanceReceipt, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &financemodels.FinanceReceiptListData{Rows: result, Total: total}, nil
}

func (d *FinanceReceiptDaoImpl) List(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*financemodels.FinanceReceiptListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinanceReceiptListData{Rows: result.Rows, Total: result.Total}, nil
}
