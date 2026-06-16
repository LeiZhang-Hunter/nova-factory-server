package settingdaoimpl

import (
	"errors"
	"nova-factory-server/app/utils/observer/integration/api"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/store/integration"

	"nova-factory-server/app/business/erp/setting/settingdao"
	"nova-factory-server/app/business/erp/setting/settingmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IntegrationConfigDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIntegrationConfigDao(db *gorm.DB) settingdao.IIntegrationConfigDao {
	i := &IntegrationConfigDaoImpl{
		db:    db,
		table: "erp_integration_config",
	}
	integration.RegisterStore(i)
	return i
}

func (i *IntegrationConfigDaoImpl) Set(c *gin.Context, req *settingmodels.IntegrationConfigSet) (*settingmodels.IntegrationConfig, error) {
	if req.Type == "" {
		return nil, errors.New("type不能为空")
	}

	var data *settingmodels.IntegrationConfig
	err := i.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var exists settingmodels.IntegrationConfig
		err := tx.Table(i.table).Where("type = ?", req.Type).Where("state = ?", commonStatus.NORMAL).First(&exists).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			data = &settingmodels.IntegrationConfig{
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
			if data.Status != nil && *data.Status {
				if err = tx.Table(i.table).Where("state = ?", commonStatus.NORMAL).Updates(map[string]any{"status": false}).Error; err != nil {
					return err
				}
			}
			data.SetCreateBy(baizeContext.GetUserId(c))
			if err = tx.Table(i.table).Create(data).Error; err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}

		update := &settingmodels.IntegrationConfig{
			ID:     exists.ID,
			Type:   exists.Type,
			Data:   req.Data,
			Status: exists.Status,
		}
		if req.Status != nil {
			update.Status = req.Status
		}
		if update.Status != nil && *update.Status {
			if err = tx.Table(i.table).Where("state = ?", commonStatus.NORMAL).Where("id <> ?", exists.ID).Updates(map[string]any{"status": false}).Error; err != nil {
				return err
			}
		}
		update.SetUpdateBy(baizeContext.GetUserId(c))
		if err = tx.Table(i.table).Where("id = ?", exists.ID).Where("state = ?", commonStatus.NORMAL).
			Select("data", "status", "update_by", "update_time").Updates(update).Error; err != nil {
			return err
		}
		if err = tx.Table(i.table).Where("id = ?", exists.ID).Where("state = ?", commonStatus.NORMAL).First(&exists).Error; err != nil {
			return err
		}
		data = &exists
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (i *IntegrationConfigDaoImpl) List(c *gin.Context, req *settingmodels.IntegrationConfigQuery) (*settingmodels.IntegrationConfigListData, error) {
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
	rows := make([]*settingmodels.IntegrationConfig, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	return &settingmodels.IntegrationConfigListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func (i *IntegrationConfigDaoImpl) GetEnabled(c *gin.Context) (*settingmodels.IntegrationConfig, error) {
	var item settingmodels.IntegrationConfig
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

func (i *IntegrationConfigDaoImpl) GetService(c *gin.Context) (api.Service, config.Config, error) {
	var item settingmodels.IntegrationConfig
	err := i.db.WithContext(c).Table(i.table).
		Where("state = ?", commonStatus.NORMAL).
		Where("status = ?", true).
		Order("id DESC").
		First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	service, err := item.Service()
	return service, &item, nil
}
