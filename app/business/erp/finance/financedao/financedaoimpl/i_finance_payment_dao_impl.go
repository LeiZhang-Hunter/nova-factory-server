package financedaoimpl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/finance/financedao"
	"nova-factory-server/app/business/erp/finance/financemodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FinancePaymentDaoImpl 提供数据访问能力。
type FinancePaymentDaoImpl struct {
	db *gorm.DB
}

// NewFinancePaymentDao 创建 DAO。
func NewFinancePaymentDao(db *gorm.DB) financedao.IFinancePaymentDao {
	return &FinancePaymentDaoImpl{db: db}
}

func (d *FinancePaymentDaoImpl) Create(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error) {
	model := new(financemodels.FinancePayment)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_finance_payment").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *FinancePaymentDaoImpl) Update(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(financemodels.FinancePayment)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_finance_payment").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *FinancePaymentDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_finance_payment").Where("id IN ?", ids)
	if erpbiz.HasField(new(financemodels.FinancePayment), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *FinancePaymentDaoImpl) GetByID(c *gin.Context, id int64) (*financemodels.FinancePayment, error) {
	item := new(financemodels.FinancePayment)
	db := d.db.WithContext(c).Table("erp_finance_payment").Where("id = ?", id)
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

func (d *FinancePaymentDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinancePayment, error) {
	if column == "" {
		return nil, nil
	}
	item := new(financemodels.FinancePayment)
	db := d.db.WithContext(c).Table("erp_finance_payment").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *FinancePaymentDaoImpl) ListPage(c *gin.Context, req *financemodels.FinancePaymentQuery) (*erpbiz.PageResult[financemodels.FinancePayment], error) {
	if req == nil {
		req = new(financemodels.FinancePaymentQuery)
	}
	db := d.db.WithContext(c).Table("erp_finance_payment")
	if erpbiz.HasField(new(financemodels.FinancePayment), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]financemodels.FinancePayment, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*financemodels.FinancePayment, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[financemodels.FinancePayment]{Rows: result, Total: total}, nil
}

func (d *FinancePaymentDaoImpl) List(c *gin.Context, req *financemodels.FinancePaymentQuery) (*financemodels.FinancePaymentListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinancePaymentListData{Rows: result.Rows, Total: result.Total}, nil
}
