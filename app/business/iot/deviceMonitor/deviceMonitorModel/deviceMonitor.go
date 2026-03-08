package deviceMonitorModel

type DeviceMonitorMetricReq struct {
	DeviceId   string `json:"device_id"`
	TemplateId string `json:"template_id"`
	DataId     string `json:"data_id"`
}

type TypeDeviceCounterRankValue struct {
	Time     int64  `json:"time"`
	Dev      string `json:"dev"`
	DevName  string `json:"dev_name"`
	DataName string `json:"data_name"`
	Value    int64  `json:"value"`
}

// TypeDeviceCounterRank 设备上报次数排名
type TypeDeviceCounterRank struct {
	Rows []*TypeDeviceCounterRankValue `json:"rows"`
}
