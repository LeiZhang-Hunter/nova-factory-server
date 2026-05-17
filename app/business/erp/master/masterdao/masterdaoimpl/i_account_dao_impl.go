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

type AccountDaoImpl struct {
	db *gorm.DB
}

func NewAccountDao(db *gorm.DB) masterdao.IAccountDao {
	return &AccountDaoImpl{db: db}
}

func (d *AccountDaoImpl) Create(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := mastermodels.AccountUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table("erp_account").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *AccountDaoImpl) Update(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.No != "" {
		updates["no"] = req.No
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	updates["status"] = req.Status
	updates["sort"] = req.Sort
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table("erp_account").Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, req.ID)
}

func (d *AccountDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table("erp_account").
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *AccountDaoImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Account, error) {
	item := new(mastermodels.Account)
	if err := d.db.WithContext(c).Table("erp_account").
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

func (d *AccountDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Account, error) {
	if column == "" {
		return nil, nil
	}
	item := new(mastermodels.Account)
	if err := d.db.WithContext(c).Table("erp_account").
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

func (d *AccountDaoImpl) ListPage(c *gin.Context, req *mastermodels.AccountQuery) (*mastermodels.AccountListData, error) {
	if req == nil {
		req = new(mastermodels.AccountQuery)
	}
	db := d.db.WithContext(c).Table("erp_account").Where("state = ?", commonStatus.NORMAL)
	db = applyAccountFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]mastermodels.Account, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*mastermodels.Account, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &mastermodels.AccountListData{Rows: result, Total: total}, nil
}

func (d *AccountDaoImpl) List(c *gin.Context, req *mastermodels.AccountQuery) (*mastermodels.AccountListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.AccountListData{Rows: result.Rows, Total: result.Total}, nil
}
