package metricdao

import "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"

// IIotStorageTemplateDao 设备模板
type IIotStorageTemplateDao interface {
	Delete(template entity.DeviceTemplate) error

	// Update 更新模板
	Update(template []entity.DeviceTemplate) error
}
