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

// FinanceReceiptDaoImpl 提供数据访问能力。
type FinanceReceiptDaoImpl struct {
	db *gorm.DB
}

// NewFinanceReceiptDao 创建 DAO。
func NewFinanceReceiptDao(db *gorm.DB) financedao.IFinanceReceiptDao {
	return &FinanceReceiptDaoImpl{db: db}
}

func (d *FinanceReceiptDaoImpl) Create(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error) {
	model := new(financemodels.FinanceReceipt)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_finance_receipt").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *FinanceReceiptDaoImpl) Update(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(financemodels.FinanceReceipt)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_finance_receipt").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *FinanceReceiptDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_finance_receipt").Where("id IN ?", ids)
	if erpbiz.HasField(new(financemodels.FinanceReceipt), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *FinanceReceiptDaoImpl) GetByID(c *gin.Context, id int64) (*financemodels.FinanceReceipt, error) {
	item := new(financemodels.FinanceReceipt)
	db := d.db.WithContext(c).Table("erp_finance_receipt").Where("id = ?", id)
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

func (d *FinanceReceiptDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinanceReceipt, error) {
	if column == "" {
		return nil, nil
	}
	item := new(financemodels.FinanceReceipt)
	db := d.db.WithContext(c).Table("erp_finance_receipt").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *FinanceReceiptDaoImpl) ListPage(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*erpbiz.PageResult[financemodels.FinanceReceipt], error) {
	if req == nil {
		req = new(financemodels.FinanceReceiptQuery)
	}
	db := d.db.WithContext(c).Table("erp_finance_receipt")
	if erpbiz.HasField(new(financemodels.FinanceReceipt), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]financemodels.FinanceReceipt, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*financemodels.FinanceReceipt, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[financemodels.FinanceReceipt]{Rows: result, Total: total}, nil
}

func (d *FinanceReceiptDaoImpl) List(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*financemodels.FinanceReceiptListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinanceReceiptListData{Rows: result.Rows, Total: result.Total}, nil
}
