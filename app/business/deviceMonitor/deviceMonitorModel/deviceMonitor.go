package deviceMonitorModel

type DeviceMonitorMetricReq struct {
	DeviceId   string `json:"device_id"`
	TemplateId string `json:"template_id"`
	DataId     string `json:"data_id"`
}
