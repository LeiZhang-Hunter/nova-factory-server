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

type FinancePaymentDaoImpl struct {
	db *gorm.DB
}

func NewFinancePaymentDao(db *gorm.DB) financedao.IFinancePaymentDao {
	return &FinancePaymentDaoImpl{db: db}
}

func (d *FinancePaymentDaoImpl) Create(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := financemodels.FinancePaymentUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_finance_payment").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *FinancePaymentDaoImpl) Update(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.PaymentTime != "" {
		if parsed, err := time.ParseInLocation("2006-01-02 15:04:05", req.PaymentTime, time.Local); err == nil {
			updates["payment_time"] = &parsed
		}
	}
	if req.No != "" {
		updates["no"] = req.No
	}
	updates["status"] = req.Status
	if req.FinanceUserID > 0 {
		updates["finance_user_id"] = req.FinanceUserID
	}
	if req.SupplierID > 0 {
		updates["supplier_id"] = req.SupplierID
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
	if req.PaymentPrice != 0 {
		updates["payment_price"] = req.PaymentPrice
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_finance_payment").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *FinancePaymentDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_finance_payment").Where("id IN ?", ids).Where("state = ?", commonStatus.NORMAL).Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *FinancePaymentDaoImpl) GetByID(c *gin.Context, id int64) (*financemodels.FinancePayment, error) {
	item := new(financemodels.FinancePayment)
	if err := d.db.WithContext(c).Table("erp_finance_payment").Where("id = ?", id).Where("state = ?", commonStatus.NORMAL).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *FinancePaymentDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinancePayment, error) {
	if column == "" {
		return nil, nil
	}
	item := new(financemodels.FinancePayment)
	if err := d.db.WithContext(c).Table("erp_finance_payment").Where(fmt.Sprintf("%s = ?", column), value).Where("state = ?", commonStatus.NORMAL).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *FinancePaymentDaoImpl) ListPage(c *gin.Context, req *financemodels.FinancePaymentQuery) (*financemodels.FinancePaymentListData, error) {
	if req == nil {
		req = new(financemodels.FinancePaymentQuery)
	}
	db := d.db.WithContext(c).Table("erp_finance_payment").Where("state = ?", commonStatus.NORMAL)
	db = applyFinancePaymentFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]financemodels.FinancePayment, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*financemodels.FinancePayment, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &financemodels.FinancePaymentListData{Rows: result, Total: total}, nil
}

func (d *FinancePaymentDaoImpl) List(c *gin.Context, req *financemodels.FinancePaymentQuery) (*financemodels.FinancePaymentListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinancePaymentListData{Rows: result.Rows, Total: result.Total}, nil
}
