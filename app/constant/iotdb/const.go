package iotdb

import (
	"fmt"
	"github.com/cespare/xxhash/v2"
	"go.uber.org/zap"
)

const (
	ROOT_DEVICE_TEMPLATE_NAME = "root.device.dev%d"
	NOVA_DEVICE_TEMPLATE      = "nova_device_template"
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
