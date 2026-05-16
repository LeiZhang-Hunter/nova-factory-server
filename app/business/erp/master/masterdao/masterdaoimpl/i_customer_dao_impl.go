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

type CustomerDaoImpl struct {
	db *gorm.DB
}

func NewCustomerDao(db *gorm.DB) masterdao.ICustomerDao {
	return &CustomerDaoImpl{db: db}
}

func (d *CustomerDaoImpl) Create(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := mastermodels.CustomerUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_customer").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *CustomerDaoImpl) Update(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error) {
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
	db := d.db.WithContext(c).Table("erp_customer").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *CustomerDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_customer").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *CustomerDaoImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Customer, error) {
	item := new(mastermodels.Customer)
	if err := d.db.WithContext(c).Table("erp_customer").
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

func (d *CustomerDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Customer, error) {
	if column == "" {
		return nil, nil
	}
	item := new(mastermodels.Customer)
	if err := d.db.WithContext(c).Table("erp_customer").
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

func (d *CustomerDaoImpl) ListPage(c *gin.Context, req *mastermodels.CustomerQuery) (*mastermodels.CustomerListData, error) {
	if req == nil {
		req = new(mastermodels.CustomerQuery)
	}
	db := d.db.WithContext(c).Table("erp_customer").Where("state = ?", commonStatus.NORMAL)
	db = applyCustomerFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]mastermodels.Customer, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*mastermodels.Customer, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &mastermodels.CustomerListData{Rows: result, Total: total}, nil
}

func (d *CustomerDaoImpl) List(c *gin.Context, req *mastermodels.CustomerQuery) (*mastermodels.CustomerListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.CustomerListData{Rows: result.Rows, Total: result.Total}, nil
}
