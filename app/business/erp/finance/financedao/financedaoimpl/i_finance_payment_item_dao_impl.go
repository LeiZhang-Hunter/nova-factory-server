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

type FinancePaymentItemDaoImpl struct {
	db *gorm.DB
}

func NewFinancePaymentItemDao(db *gorm.DB) financedao.IFinancePaymentItemDao {
	return &FinancePaymentItemDaoImpl{db: db}
}

func (d *FinancePaymentItemDaoImpl) Create(c *gin.Context, req *financemodels.FinancePaymentItemUpsert) (*financemodels.FinancePaymentItem, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := financemodels.FinancePaymentItemUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_finance_payment_item").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *FinancePaymentItemDaoImpl) Update(c *gin.Context, req *financemodels.FinancePaymentItemUpsert) (*financemodels.FinancePaymentItem, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.PaymentID > 0 {
		updates["payment_id"] = req.PaymentID
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
	if req.PaidPrice != 0 {
		updates["paid_price"] = req.PaidPrice
	}
	if req.PaymentPrice != 0 {
		updates["payment_price"] = req.PaymentPrice
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_finance_payment_item").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *FinancePaymentItemDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_finance_payment_item").Where("id IN ?", ids).Where("state = ?", commonStatus.NORMAL).Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *FinancePaymentItemDaoImpl) GetByID(c *gin.Context, id int64) (*financemodels.FinancePaymentItem, error) {
	item := new(financemodels.FinancePaymentItem)
	if err := d.db.WithContext(c).Table("erp_finance_payment_item").Where("id = ?", id).Where("state = ?", commonStatus.NORMAL).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *FinancePaymentItemDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinancePaymentItem, error) {
	if column == "" {
		return nil, nil
	}
	item := new(financemodels.FinancePaymentItem)
	if err := d.db.WithContext(c).Table("erp_finance_payment_item").Where(fmt.Sprintf("%s = ?", column), value).Where("state = ?", commonStatus.NORMAL).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *FinancePaymentItemDaoImpl) ListPage(c *gin.Context, req *financemodels.FinancePaymentItemQuery) (*financemodels.FinancePaymentItemListData, error) {
	if req == nil {
		req = new(financemodels.FinancePaymentItemQuery)
	}
	db := d.db.WithContext(c).Table("erp_finance_payment_item").Where("state = ?", commonStatus.NORMAL)
	db = applyFinancePaymentItemFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]financemodels.FinancePaymentItem, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*financemodels.FinancePaymentItem, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &financemodels.FinancePaymentItemListData{Rows: result, Total: total}, nil
}

func (d *FinancePaymentItemDaoImpl) List(c *gin.Context, req *financemodels.FinancePaymentItemQuery) (*financemodels.FinancePaymentItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinancePaymentItemListData{Rows: result.Rows, Total: result.Total}, nil
}
