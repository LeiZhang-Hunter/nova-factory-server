package deviceModels

import (
	"encoding/json"
	"github.com/gogf/gf/util/gconv"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/constant/device"
)

type DeviceListReq struct {
	DeviceId         int64   `json:"deviceId,string" form:"deviceId" jsonschema:"description=设备id，数据库主键"`
	DeviceGroupId    int64   `json:"deviceGroupId,string" form:"deviceGroupId" jsonschema:"description=设备分组id"`
	DeviceClassId    int64   `json:"deviceClassId,string" form:"deviceClassId" jsonschema:"description=设备分类id"`
	DeviceProtocolId int64   `json:"deviceProtocolId,string" form:"deviceProtocolId" jsonschema:"description=设备模板id,设备的协议解析规范，目前支持modbus协议解析和mqtt"`
	DeviceBuildingId int64   `json:"deviceBuildingId,string" form:"deviceBuildingId" jsonschema:"description=设备所在的建筑物id"`
	Start            uint64  `json:"start"`
	End              uint64  `json:"end"`
	Name             *string `json:"name" form:"name" jsonschema:"description=设备名字"`
	Number           *string `json:"number" form:"number" jsonschema:"description=设备标签，用来给设备打上标识"`
	Type             *string `json:"type" form:"type" jsonschema:"description=设备类型"`
	ControlType      *int    `json:"controlType" form:"controlType" jsonschema:"description=设备控制类型"`
	baize.BaseEntityDQL
}

type DeviceInfo struct {
	DeviceId          uint64    `json:"deviceId,string" db:"device_id"`
	DeviceGroupId     uint64    `json:"deviceGroupId,string" db:"device_group_id"`
	DeviceClassId     uint64    `json:"deviceClassId,string" db:"device_class_id"`
	CommunicationType string    `gorm:"column:communication_type;comment:通信方式" json:"communication_type"`               // 通信方式
	ProtocolType      string    `gorm:"column:protocol_type;comment:协议类型" json:"protocol_type"`                         // 协议类型
	DeviceGatewayID   int64     `gorm:"column:device_gateway_id;not null;comment:网关id" json:"device_gateway_id,string"` // 网关id
	DeviceProtocolId  uint64    `json:"deviceProtocolId,string" db:"device_protocol_id"`
	DeviceBuildingId  uint64    `json:"deviceBuildingId,string" db:"device_building_id"`
	Name              *string   `json:"Name" db:"name"`
	Number            *string   `json:"Number" db:"number"`
	Type              *string   `json:"Type" db:"type"`
	Action            []*string `json:"Action" db:"action"`
	Extension         *string   `json:"Extension" db:"extension"`
	Enable            *bool     `json:"enable"`
	ControlType       int       `json:"ControlType" db:"control_type"`
}

type LocalInfo struct {
	Slave    int    `json:"slave,omitempty"`
	Address  string `json:"address"`
	Quantity uint16 `json:"quantity"`
}

type NetworkInfo struct {
	Slave int `json:"Name,omitempty"`
}

type MqttInfo struct {
	Address  string `json:"address"`
	Topic    string `json:"topic"`
	ClientId string `json:"client_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// "{\"localInfo\":[{\"salve\":1,\"address\":\"81.71.98.26:11808\"}]}"
type ExtensionInfo struct {
	LocalInfo     []LocalInfo `json:"localInfo,omitempty"`
	LocalMqttInfo []MqttInfo  `json:"mqttInfo,omitempty"`
	//NetworkInfo NetworkInfo `json:"networkInfo,omitempty"`
}

type DeviceVO struct {
	DeviceId          uint64                                               `json:"deviceId,string" db:"device_id"`
	DeviceGroupId     uint64                                               `json:"deviceGroupId,string" db:"device_group_id"`
	DeviceClassId     uint64                                               `json:"deviceClassId,string" db:"device_class_id"`
	DeviceProtocolId  uint64                                               `json:"deviceProtocolId,string" db:"device_protocol_id"`
	DeviceBuildingId  uint64                                               `json:"deviceBuildingId,string" db:"device_building_id"`
	CommunicationType string                                               `gorm:"column:communication_type;comment:通信方式" json:"communication_type"`               // 通信方式
	ProtocolType      string                                               `gorm:"column:protocol_type;comment:协议类型" json:"protocol_type"`                         // 协议
	DeviceGatewayID   int64                                                `gorm:"column:device_gateway_id;not null;comment:网关id" json:"device_gateway_id,string"` // 网关id
	Name              *string                                              `json:"Name" db:"name"`
	Number            *string                                              `json:"Number" db:"number"`
	Type              *string                                              `json:"Type" db:"type"`
	Action            string                                               `json:"Action" db:"action"`
	Enable            *bool                                                `json:"enable" gorm:"column:enable;comment:是否启用 0 禁用 1 启用"`
	Status            device.RUN_STATUS                                    `json:"status" gorm:"column:status;comment:是否启用 0 禁用 1 启用" `
	Extension         string                                               `json:"Extension" db:"extension"`
	ControlType       int                                                  `json:"ControlType" db:"control_type"`
	ExtensionInfo     *ExtensionInfo                                       `json:"extension_info,omitempty" gorm:"-" db:"control_type"`
	TemplateList      map[uint64]map[uint64]*metricModels.DeviceMetricData `json:"template_list" gorm:"-" db:"template_list"`
	Active            bool                                                 `json:"active" gorm:"-" db:"active"`
	Exception         bool                                                 `json:"exception" gorm:"-"`
	baize.BaseEntity
}

func NewDeviceVO(device *DeviceInfo) *DeviceVO {
	var vo DeviceVO
	gconv.Scan(device, &vo)
	if device.Action != nil {
		action, _ := json.Marshal(device.Action)
		vo.Action = string(action)
	}
	if device.Extension != nil {
		vo.Extension = *device.Extension
	}
	return &vo
}

type DeviceValue struct {
	DeviceId          uint64            `json:"deviceId,string"`
	DeviceGroupId     uint64            `json:"deviceGroupId,string"`
	DeviceClassId     uint64            `json:"deviceClassId,string"`
	DeviceProtocolId  uint64            `json:"deviceProtocolId,string"`
	DeviceBuildingId  uint64            `json:"deviceBuildingId,string"`
	Name              string            `json:"name"`
	CommunicationType string            `gorm:"column:communication_type;comment:通信方式" json:"communication_type"`               // 通信方式
	ProtocolType      string            `gorm:"column:protocol_type;comment:协议类型" json:"protocol_type"`                         // 协议
	DeviceGatewayID   int64             `gorm:"column:device_gateway_id;not null;comment:网关id" json:"device_gateway_id,string"` // 网关id
	DeviceGroupName   string            `json:"deviceGroupName"`
	Number            string            `json:"number"`
	Type              string            `json:"type"`
	Action            []string          `json:"action"`
	Extension         string            `json:"extension"`
	ControlType       int               `json:"controlType"`
	Enable            *bool             `json:"enable" gorm:"column:enable;comment:是否启用 0 禁用 1 启用"`
	Status            device.RUN_STATUS `json:"status"`
	CreateUserName    string            `json:"createUserName"`
	UpdateUserName    string            `json:"updateUserName"`
	baize.BaseEntity
}

type DeviceInfoListData struct {
	Rows  []*DeviceVO `json:"rows"`
	Total int64       `json:"total"`
}

type DeviceInfoListValue struct {
	Rows  []*DeviceValue `json:"rows"`
	Total int64          `json:"total"`
}

// DeviceTagListReq 根据标签读取设备
type DeviceTagListReq struct {
	Tags string `form:"tags"`
	baize.BaseEntityDQL
}

// DeviceMetricInfoListValue 根据tag 返回数据结构
type DeviceMetricInfoListValue struct {
	Rows  []*DeviceVO `json:"rows"`
	Total int64       `json:"total"`
}

// DeviceStatData 根据tag 返回数据结构
type DeviceStatData struct {
	Total       int64 `json:"total"`
	Online      int64 `json:"online"`
	OffLine     int64 `json:"offline"`
	Exception   int64 `json:"exception"`
	Maintenance int64 `json:"maintenance"`
}
