package masterdaoimpl

import (
	"errors"
	"fmt"
	"time"

	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SupplierDaoImpl struct {
	db *gorm.DB
}

func NewSupplierDao(db *gorm.DB) masterdao.ISupplierDao {
	return &SupplierDaoImpl{db: db}
}

func (d *SupplierDaoImpl) Create(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := mastermodels.SupplierUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_supplier").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *SupplierDaoImpl) Update(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Contact != "" {
		updates["contact"] = req.Contact
	}
	if req.Mobile != "" {
		updates["mobile"] = req.Mobile
	}
	if req.Telephone != "" {
		updates["telephone"] = req.Telephone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Fax != "" {
		updates["fax"] = req.Fax
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	updates["status"] = req.Status
	updates["sort"] = req.Sort
	if req.TaxNo != "" {
		updates["tax_no"] = req.TaxNo
	}
	if req.TaxPercent != 0 {
		updates["tax_percent"] = req.TaxPercent
	}
	if req.BankName != "" {
		updates["bank_name"] = req.BankName
	}
	if req.BankAccount != "" {
		updates["bank_account"] = req.BankAccount
	}
	if req.BankAddress != "" {
		updates["bank_address"] = req.BankAddress
	}
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_supplier").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *SupplierDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_supplier").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *SupplierDaoImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Supplier, error) {
	item := new(mastermodels.Supplier)
	if err := d.db.WithContext(c).Table("erp_supplier").
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

func (d *SupplierDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Supplier, error) {
	if column == "" {
		return nil, nil
	}
	item := new(mastermodels.Supplier)
	if err := d.db.WithContext(c).Table("erp_supplier").
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

func (d *SupplierDaoImpl) ListPage(c *gin.Context, req *mastermodels.SupplierQuery) (*mastermodels.SupplierListData, error) {
	if req == nil {
		req = new(mastermodels.SupplierQuery)
	}
	db := d.db.WithContext(c).Table("erp_supplier").Where("state = ?", commonStatus.NORMAL)
	db = applySupplierFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]mastermodels.Supplier, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*mastermodels.Supplier, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &mastermodels.SupplierListData{Rows: result, Total: total}, nil
}

func (d *SupplierDaoImpl) List(c *gin.Context, req *mastermodels.SupplierQuery) (*mastermodels.SupplierListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.SupplierListData{Rows: result.Rows, Total: result.Total}, nil
}
