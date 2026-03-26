package orderDaoImpl

import (
	"errors"

	"nova-factory-server/app/business/erp/order/orderDao"
	"nova-factory-server/app/business/erp/setting/settingModels"
	"nova-factory-server/app/constant/commonStatus"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewOrderDao(db *gorm.DB) orderDao.IOrderDao {
	return &OrderDaoImpl{
		db:    db,
		table: "erp_integration_config",
	}
}

func (o *OrderDaoImpl) GetEnabledGJPCfg(c *gin.Context) (*settingModels.IntegrationConfig, error) {
	var item settingModels.IntegrationConfig
	err := o.db.WithContext(c).Table(o.table).
		Where("state = ?", commonStatus.NORMAL).
		Where("status = ?", true).
		Where("type LIKE ?", "%管家婆%").
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
