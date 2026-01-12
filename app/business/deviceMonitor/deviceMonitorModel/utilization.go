package deviceMonitorModel

import (
	"nova-factory-server/app/business/metric/device/metricModels"
)

// DeviceUtilizationReq 稼动率请求
type DeviceUtilizationReq struct {
	Start uint64 `json:"start" form:"start"`
	End   uint64 `json:"end" form:"end"`
}

// DeviceStatus 计算设备在某种状态的运行时间
type DeviceStatus struct {
	Time     int64   `json:"time"`
	DeviceId int64   `json:"deviceId"`
	Value    float64 `json:"value"`
	Status   int     `json:"status"`
}

type DeviceStatusList struct {
	List []DeviceStatus
}

func NewDeviceStatusList() *DeviceStatusList {
	return &DeviceStatusList{
		List: make([]DeviceStatus, 0),
	}
}

// DeviceUtilizationData 稼动率数据
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

type DeviceRunProcess struct {
	List            []DeviceProcessStatus    `json:"list"`
	DeviceName      string                   `json:"deviceName"`
	UtilizationRate float64                  `json:"utilization_rate"` //稼动率
	WaitRate        float64                  `json:"wait_rate"`        //待机率
	BuildingName    string                   `json:"buildingName"`
	StatusMap       map[int]DeviceStatusData `json:"statusMap"`
}

type DeviceProcessStatus struct {
	Time     int64           `json:"time"`
	DeviceId int64           `json:"deviceId"`
	Value    map[int]float64 `json:"value"`
}

// DeviceUtilizationPublicDataList 稼动率报表
type DeviceUtilizationPublicDataList struct {
	List        []*DeviceUtilizationData      `json:"list"`
	Total       uint64                        `json:"total"`
	Running     uint64                        `json:"running"`
	WaitTing    uint64                        `json:"waitTing"`
	Stopped     uint64                        `json:"stopped"`
	ProcessList map[string][]DeviceRunProcess `json:"processList"` // 进度列表 key 是建筑物 => 进度
	RunTop5     []*DeviceUtilizationData      `json:"runTop5"`     // 运行TOP5
	WaitTop5    []*DeviceUtilizationData      `json:"waitTop5"`    // 待机top5
	Radio       map[string][]DeviceRadio      `json:"radio"`       // 建筑物 => 比率
	Data        *metricModels.MetricQueryData `json:"data"`
}

// DeviceRunStat 设备运行状态
type DeviceRunStat struct {
	Time   int64  `json:"time"`
	Dev    string `json:"deviceId"`
	Status int    `json:"value"`
}

// DeviceProcessList 设备运行进度
type DeviceProcessList struct {
	List map[string][]DeviceStatus
}

// DeviceRadio 设备比率
type DeviceRadio struct {
	Time    int64   `json:"time"`
	Dev     string  `json:"deviceId"`
	Status  int     `json:"status"`
	Percent float64 `json:"percent"`
}

type DeviceStatusData struct {
	Time    uint64  `json:"time"`     //运行时间
	TimeStr string  `json:"time_str"` //运行时间
	Rate    float64 `json:"rate"`     //稼动率
	RateStr string  `json:"rate_str"` //稼动率
}

// DeviceUtilizationDataV2 稼动率数据,第二版本通用版本
type DeviceUtilizationDataV2 struct {
	DeviceId   int64  `json:"device_id"`   //设备id
	Building   string `json:"building"`    //建筑物渲染
	DeviceName string `json:"device_name"` //设备名字

	StatusMap map[int]DeviceStatusData `json:"status_map"`
}

// DeviceUtilizationPublicDataListV2 稼动率报表数据
type DeviceUtilizationPublicDataListV2 struct {
	List  []*DeviceUtilizationDataV2 `json:"list"`
	Total uint64                     `json:"total"`

	ProcessList map[string][]DeviceRunProcess      `json:"processList"` // 进度列表 key 是建筑物 => 进度
	Top5        map[int][]*DeviceUtilizationDataV2 `json:"runTop5"`     // 运行TOP5
	Radio       map[string][]DeviceRadio           `json:"radio"`       // 建筑物 => 比率
	Data        *metricModels.MetricQueryData      `json:"data"`

	StatusCount map[int]uint64 `json:"statusCount"`
}
