package masterdaoimpl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AccountDaoImpl 提供数据访问能力。
type AccountDaoImpl struct {
	db *gorm.DB
}

// NewAccountDao 创建 DAO。
func NewAccountDao(db *gorm.DB) masterdao.IAccountDao {
	return &AccountDaoImpl{db: db}
}

func (d *AccountDaoImpl) Create(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error) {
	model := new(mastermodels.Account)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_account").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *AccountDaoImpl) Update(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(mastermodels.Account)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_account").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *AccountDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_account").Where("id IN ?", ids)
	if erpbiz.HasField(new(mastermodels.Account), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *AccountDaoImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Account, error) {
	item := new(mastermodels.Account)
	db := d.db.WithContext(c).Table("erp_account").Where("id = ?", id)
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

func (d *AccountDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Account, error) {
	if column == "" {
		return nil, nil
	}
	item := new(mastermodels.Account)
	db := d.db.WithContext(c).Table("erp_account").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *AccountDaoImpl) ListPage(c *gin.Context, req *mastermodels.AccountQuery) (*erpbiz.PageResult[mastermodels.Account], error) {
	if req == nil {
		req = new(mastermodels.AccountQuery)
	}
	db := d.db.WithContext(c).Table("erp_account")
	if erpbiz.HasField(new(mastermodels.Account), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]mastermodels.Account, 0)
	orderBy := strings.TrimSpace("sort ASC, id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*mastermodels.Account, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[mastermodels.Account]{Rows: result, Total: total}, nil
}

func (d *AccountDaoImpl) List(c *gin.Context, req *mastermodels.AccountQuery) (*mastermodels.AccountListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.AccountListData{Rows: result.Rows, Total: result.Total}, nil
}
