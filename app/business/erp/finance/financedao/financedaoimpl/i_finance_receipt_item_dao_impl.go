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

// FinanceReceiptItemDaoImpl 提供数据访问能力。
type FinanceReceiptItemDaoImpl struct {
	db *gorm.DB
}

// NewFinanceReceiptItemDao 创建 DAO。
func NewFinanceReceiptItemDao(db *gorm.DB) financedao.IFinanceReceiptItemDao {
	return &FinanceReceiptItemDaoImpl{db: db}
}

func (d *FinanceReceiptItemDaoImpl) Create(c *gin.Context, req *financemodels.FinanceReceiptItemUpsert) (*financemodels.FinanceReceiptItem, error) {
	model := new(financemodels.FinanceReceiptItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_finance_receipt_item").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *FinanceReceiptItemDaoImpl) Update(c *gin.Context, req *financemodels.FinanceReceiptItemUpsert) (*financemodels.FinanceReceiptItem, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(financemodels.FinanceReceiptItem)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_finance_receipt_item").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *FinanceReceiptItemDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_finance_receipt_item").Where("id IN ?", ids)
	if erpbiz.HasField(new(financemodels.FinanceReceiptItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *FinanceReceiptItemDaoImpl) GetByID(c *gin.Context, id int64) (*financemodels.FinanceReceiptItem, error) {
	item := new(financemodels.FinanceReceiptItem)
	db := d.db.WithContext(c).Table("erp_finance_receipt_item").Where("id = ?", id)
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

func (d *FinanceReceiptItemDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinanceReceiptItem, error) {
	if column == "" {
		return nil, nil
	}
	item := new(financemodels.FinanceReceiptItem)
	db := d.db.WithContext(c).Table("erp_finance_receipt_item").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *FinanceReceiptItemDaoImpl) ListPage(c *gin.Context, req *financemodels.FinanceReceiptItemQuery) (*erpbiz.PageResult[financemodels.FinanceReceiptItem], error) {
	if req == nil {
		req = new(financemodels.FinanceReceiptItemQuery)
	}
	db := d.db.WithContext(c).Table("erp_finance_receipt_item")
	if erpbiz.HasField(new(financemodels.FinanceReceiptItem), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]financemodels.FinanceReceiptItem, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*financemodels.FinanceReceiptItem, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[financemodels.FinanceReceiptItem]{Rows: result, Total: total}, nil
}

func (d *FinanceReceiptItemDaoImpl) List(c *gin.Context, req *financemodels.FinanceReceiptItemQuery) (*financemodels.FinanceReceiptItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinanceReceiptItemListData{Rows: result.Rows, Total: result.Total}, nil
}
