package settingDaoImpl

import (
	"errors"

	"nova-factory-server/app/business/erp/setting/settingDao"
	"nova-factory-server/app/business/erp/setting/settingModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IntegrationConfigDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIntegrationConfigDao(db *gorm.DB) settingDao.IIntegrationConfigDao {
	return &IntegrationConfigDaoImpl{
		db:    db,
		table: "erp_integration_config",
	}
}

func (i *IntegrationConfigDaoImpl) Set(c *gin.Context, req *settingModels.IntegrationConfigSet) (*settingModels.IntegrationConfig, error) {
	if req.Type == "" {
		return nil, errors.New("type不能为空")
	}
	var exists settingModels.IntegrationConfig
	err := i.db.WithContext(c).Table(i.table).Where("type = ?", req.Type).Where("state = ?", commonStatus.NORMAL).First(&exists).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		data := &settingModels.IntegrationConfig{
			Type:   req.Type,
			Data:   req.Data,
			Status: req.Status,
			DeptID: baizeContext.GetDeptId(c),
			State:  commonStatus.NORMAL,
		}
		if data.Status == nil {
			status := true
			data.Status = &status
		}
		data.SetCreateBy(baizeContext.GetUserId(c))
		if err = i.db.WithContext(c).Table(i.table).Create(data).Error; err != nil {
			return nil, err
		}
		return data, nil
	}
	if err != nil {
		return nil, err
	}
	update := &settingModels.IntegrationConfig{
		ID:     exists.ID,
		Type:   exists.Type,
		Data:   req.Data,
		Status: exists.Status,
	}
	if req.Status != nil {
		update.Status = req.Status
	}
	update.SetUpdateBy(baizeContext.GetUserId(c))
	if err = i.db.WithContext(c).Table(i.table).Where("id = ?", exists.ID).Where("state = ?", commonStatus.NORMAL).
		Select("data", "status", "update_by", "update_time").Updates(update).Error; err != nil {
		return nil, err
	}
	if err = i.db.WithContext(c).Table(i.table).Where("id = ?", exists.ID).Where("state = ?", commonStatus.NORMAL).First(&exists).Error; err != nil {
		return nil, err
	}
	return &exists, nil
}

func (i *IntegrationConfigDaoImpl) List(c *gin.Context, req *settingModels.IntegrationConfigQuery) (*settingModels.IntegrationConfigListData, error) {
	db := i.db.WithContext(c).Table(i.table).Where("state = ?", commonStatus.NORMAL)
	if req != nil && req.Type != "" {
		db = db.Where("type LIKE ?", "%"+req.Type+"%")
	}
	if req != nil && req.Status != nil {
		db = db.Where("status = ?", req.Status)
	}
	page := int64(1)
	size := int64(20)
	if req != nil && req.Page > 0 {
		page = req.Page
	}
	if req != nil && req.Size > 0 {
		size = req.Size
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*settingModels.IntegrationConfig, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	return &settingModels.IntegrationConfigListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func (i *IntegrationConfigDaoImpl) GetEnabled(c *gin.Context) (*settingModels.IntegrationConfig, error) {
	var item settingModels.IntegrationConfig
	err := i.db.WithContext(c).Table(i.table).
		Where("state = ?", commonStatus.NORMAL).
		Where("status = ?", true).
		Order("id DESC").
		First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}
