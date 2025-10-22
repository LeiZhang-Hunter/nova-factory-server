package deviceMonitorModel

import "nova-factory-server/app/constant/device"

// DeviceUtilizationReq 稼动率请求
type DeviceUtilizationReq struct {
	Start uint64 `json:"start"`
	End   uint64 `json:"end"`
}

type DeviceStatus struct {
	DeviceId int64             `json:"deviceId"`
	Value    float64           `json:"value"`
	Status   device.RUN_STATUS `json:"status"`
}

type DeviceStatusList struct {
	List []DeviceStatus
}

func NewDeviceStatusList() *DeviceStatusList {
	return &DeviceStatusList{
		List: make([]DeviceStatus, 0),
	}
}

type DeviceUtilizationData struct {
	DeviceId           int64   `json:"device_id"`            //设备id
	Building           string  `json:"building"`             //建筑物渲染
	DeviceName         string  `json:"device_name"`          //设备名字
	RunTime            uint64  `json:"run_time"`             //运行时间
	RunTimeStr         string  `json:"run_time_str"`         //运行时间
	UtilizationRate    float64 `json:"utilization_rate"`     //稼动率
	UtilizationRateStr string  `json:"utilization_rate_str"` //稼动率
	StopTime           uint64  `json:"stop_time"`            //停机时间
	StopTimeStr        string  `json:"stop_time_str"`        //停机时间
	StopRate           float64 `json:"stop_rate"`            //停机率
	StopRateStr        string  `json:"stop_rate_str"`        //停机率
	WaitTime           uint64  `json:"wait_time"`            //待机时间
	WaitTimeStr        string  `json:"wait_time_str"`        //待机时间
	WaitRate           float64 `json:"wait_rate"`            //待机率
	WaitRateStr        string  `json:"wait_rate_str"`        //待机率
}

// DeviceUtilizationDataList 稼动率报表
type DeviceUtilizationDataList struct {
	List []*DeviceUtilizationData `json:"list"`
}
