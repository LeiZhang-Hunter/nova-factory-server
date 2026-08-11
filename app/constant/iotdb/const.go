package iotdb

import "strconv"

const (
	ROOT_DEVICE_TEMPLATE_NAME            = "root.device.dev"
	NOVA_DEVICE_TEMPLATE                 = "nova_device_template"
	NOVA_DEVICE_RUN_TEMPLATE             = "nova_device_running_template"
	NOVA_DEVICE_METRIC_TEMPLATE_PREFIX   = "nova_device_template_"
	ROOT_RUN_STATUS_DEVICE_TEMPLATE_NAME = "root.run_status_device.dev"
)

func MakeDeviceTemplateName(templateId int64) string {
	return NOVA_DEVICE_METRIC_TEMPLATE_PREFIX + strconv.FormatInt(templateId, 10)
}

func MakeDeviceDataName(deviceId int64) string {
	return ROOT_DEVICE_TEMPLATE_NAME + strconv.FormatInt(deviceId, 10)
}

func MakeDeviceDataPath(deviceId int64, dataId int64) string {
	return MakeDeviceDataName(deviceId) + ".d" + strconv.FormatInt(dataId, 10)
}

func MakeRunDeviceTemplateName(deviceId int64) string {
	return ROOT_RUN_STATUS_DEVICE_TEMPLATE_NAME + strconv.FormatInt(deviceId, 10)
}
