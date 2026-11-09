package devicedaoImpl

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	catalogresponse "nova-factory-server/app/business/iot/daemonize/daemonizemodels/response"
)

// GetGatewayCatalog reads assigned local assets, including cloud-disabled ones.
func (s *sysDeviceDataDao) GetGatewayCatalog(ctx context.Context, id int64) (*catalogresponse.CatalogSnapshot, error) {
	result := &catalogresponse.CatalogSnapshot{}
	err := s.ms.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(s.tableName).
			Where("device_gateway_id = ? AND state = ? AND communication_type = ?", id, 0, "local").
			Order("device_id").
			Find(&result.Devices).Error; err != nil {
			return err
		}
		ids := make([]uint64, 0, len(result.Devices))
		for _, d := range result.Devices {
			ids = append(ids, d.DeviceProtocolId)
		}
		if len(ids) == 0 {
			return nil
		}
		if err := tx.Table("sys_device_template").
			Where("template_id IN ? AND state = ?", ids, 0).
			Order("template_id").
			Find(&result.Templates).Error; err != nil {
			return err
		}
		return tx.Table("sys_modbus_device_config_data").
			Where("template_id IN ? AND state = ?", ids, 0).
			Order("device_config_id").
			Find(&result.Points).Error
	}, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	return result, err
}
