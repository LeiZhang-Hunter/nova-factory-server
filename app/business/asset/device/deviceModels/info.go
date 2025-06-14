package deviceModels

import (
	"encoding/json"
	"github.com/gogf/gf/util/gconv"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/metric/device/metricModels"
)

type DeviceListReq struct {
	DeviceId         int64   `json:"deviceId,string" form:"deviceId"`
	DeviceGroupId    int64   `json:"deviceGroupId,string" form:"deviceGroupId"`
	DeviceClassId    int64   `json:"deviceClassId,string" form:"deviceClassId"`
	DeviceProtocolId int64   `json:"deviceProtocolId,string" form:"deviceProtocolId"`
	DeviceBuildingId int64   `json:"deviceBuildingId,string" form:"deviceBuildingId"`
	Name             *string `json:"name" form:"name"`
	Number           *string `json:"number" form:"number"`
	Type             *string `json:"type" form:"type"`
	ControlType      int     `json:"controlType" form:"controlType"`
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
	ControlType       int       `json:"ControlType" db:"control_type"`
}

type LocalInfo struct {
	Slave   int    `json:"slave,omitempty"`
	Address string `json:"address"`
}

type NetworkInfo struct {
	Slave int `json:"Name,omitempty"`
}

// "{\"localInfo\":[{\"salve\":1,\"address\":\"81.71.98.26:11808\"}]}"
type ExtensionInfo struct {
	LocalInfo []LocalInfo `json:"localInfo,omitempty"`
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
	Extension         string                                               `json:"Extension" db:"extension"`
	ControlType       int                                                  `json:"ControlType" db:"control_type"`
	ExtensionInfo     *ExtensionInfo                                       `json:"extension_info,omitempty" gorm:"-" db:"control_type"`
	TemplateList      map[uint64]map[uint64]*metricModels.DeviceMetricData `json:"template_list,omitempty" gorm:"-" db:"template_list"`
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
	DeviceId          uint64   `json:"deviceId,string"`
	DeviceGroupId     uint64   `json:"deviceGroupId,string"`
	DeviceClassId     uint64   `json:"deviceClassId,string"`
	DeviceProtocolId  uint64   `json:"deviceProtocolId,string"`
	DeviceBuildingId  uint64   `json:"deviceBuildingId,string"`
	Name              string   `json:"name"`
	CommunicationType string   `gorm:"column:communication_type;comment:通信方式" json:"communication_type"`               // 通信方式
	ProtocolType      string   `gorm:"column:protocol_type;comment:协议类型" json:"protocol_type"`                         // 协议
	DeviceGatewayID   int64    `gorm:"column:device_gateway_id;not null;comment:网关id" json:"device_gateway_id,string"` // 网关id
	DeviceGroupName   string   `json:"deviceGroupName"`
	Number            string   `json:"number"`
	Type              string   `json:"type"`
	Action            []string `json:"action"`
	Extension         string   `json:"extension"`
	ControlType       int      `json:"controlType"`
	CreateUserName    string   `json:"createUserName"`
	UpdateUserName    string   `json:"updateUserName"`
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
