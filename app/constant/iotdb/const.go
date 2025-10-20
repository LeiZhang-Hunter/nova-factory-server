package iotdb

import (
	"fmt"
	"github.com/cespare/xxhash/v2"
	"go.uber.org/zap"
)

const (
	ROOT_DEVICE_TEMPLATE_NAME            = "root.device.dev%d"
	NOVA_DEVICE_TEMPLATE                 = "nova_device_template"
	NOVA_DEVICE_RUN_TEMPLATE             = "nova_device_running_template"
	ROOT_RUN_STATUS_DEVICE_TEMPLATE_NAME = "root.run_status_device.dev%d"
)

func MakeDeviceTemplateName(deviceId int64, templateId int64, dataId int64) string {
	hash := xxhash.New()
	_, err := hash.WriteString(fmt.Sprintf("%d-%d-%d", deviceId, templateId, dataId))
	if err != nil {
		zap.L().Error("iotdb.MakeDeviceTemplateName", zap.Error(err))
		return ""
	}
	return fmt.Sprintf(ROOT_DEVICE_TEMPLATE_NAME, hash.Sum64())
}

func MakeRunDeviceTemplateName(deviceId int64) string {
	return fmt.Sprintf(ROOT_RUN_STATUS_DEVICE_TEMPLATE_NAME, deviceId)
}
