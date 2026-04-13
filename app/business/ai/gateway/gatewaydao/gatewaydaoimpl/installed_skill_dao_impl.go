package gatewaydaoimpl

import (
	"errors"
	"time"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InstalledSkillDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewInstalledSkillDao(db *gorm.DB) gatewaydao.IInstalledSkillDao {
	return &InstalledSkillDaoImpl{
		db:    db,
		table: "installed_skills",
	}
}

func (i *InstalledSkillDaoImpl) Create(c *gin.Context, req *gatewaymodels.InstalledSkillUpsert) (*gatewaymodels.InstalledSkill, error) {
	item := &gatewaymodels.InstalledSkill{
		ID:          snowflake.GenID(),
		Name:        req.Name,
		Slug:        req.Slug,
		Version:     req.Version,
		Source:      req.Source,
		Description: req.Description,
		Enabled:     req.Enabled,
		DeptID:      baizeContext.GetDeptId(c),
		State:       commonStatus.NORMAL,
	}
	item.SetCreateBy(baizeContext.GetUserId(c))
	if err := i.db.WithContext(c).Table(i.table).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (i *InstalledSkillDaoImpl) Update(c *gin.Context, req *gatewaymodels.InstalledSkillUpsert) (*gatewaymodels.InstalledSkill, error) {
	item := &gatewaymodels.InstalledSkill{
		ID:          req.ID,
		Name:        req.Name,
		Slug:        req.Slug,
		Version:     req.Version,
		Source:      req.Source,
		Description: req.Description,
		Enabled:     req.Enabled,
	}
	item.SetUpdateBy(baizeContext.GetUserId(c))
	if err := i.db.WithContext(c).Table(i.table).Where("id = ?", item.ID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Select("name", "slug", "version", "source", "description", "enabled", "update_by", "update_time").
		Updates(item).Error; err != nil {
		return nil, err
	}
	return i.GetByID(c, item.ID)
}

func (i *InstalledSkillDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	now := time.Now()
	return i.db.WithContext(c).Table(i.table).Where("id IN ?", ids).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": now,
		}).Error
}

func (i *InstalledSkillDaoImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.InstalledSkill, error) {
	var item gatewaymodels.InstalledSkill
	if err := i.db.WithContext(c).Table(i.table).Where("id = ?", id).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (i *InstalledSkillDaoImpl) GetBySlug(c *gin.Context, slug string) (*gatewaymodels.InstalledSkill, error) {
	var item gatewaymodels.InstalledSkill
	if err := i.db.WithContext(c).Table(i.table).Where("slug = ?", slug).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (i *InstalledSkillDaoImpl) List(c *gin.Context, req *gatewaymodels.InstalledSkillQuery) (*gatewaymodels.InstalledSkillListData, error) {
	db := i.db.WithContext(c).Table(i.table).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL)
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Slug != "" {
		db = db.Where("slug LIKE ?", "%"+req.Slug+"%")
	}
	if req.Enabled != nil {
		db = db.Where("enabled = ?", req.Enabled)
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*gatewaymodels.InstalledSkill, 0)
	if err := db.Order("id DESC").Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	return &gatewaymodels.InstalledSkillListData{
		Rows:  rows,
		Total: total,
	}, nil
}
