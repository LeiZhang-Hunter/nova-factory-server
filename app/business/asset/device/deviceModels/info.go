package deviceModels

import (
	"encoding/json"
	"github.com/gogf/gf/util/gconv"
	"nova-factory-server/app/baize"
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
	DeviceId         uint64             `json:"deviceId,string" db:"device_id"`
	DeviceGroupId    uint64             `json:"deviceGroupId,string" db:"device_group_id"`
	DeviceClassId    uint64             `json:"deviceClassId,string" db:"device_class_id"`
	DeviceProtocolId uint64             `json:"deviceProtocolId,string" db:"device_protocol_id"`
	DeviceBuildingId uint64             `json:"deviceBuildingId,string" db:"device_building_id"`
	Name             *string            `json:"Name" db:"name"`
	Number           *string            `json:"Number" db:"number"`
	Type             *string            `json:"Type" db:"type"`
	Action           []*string          `json:"Action" db:"action"`
	Extension        map[string]*string `json:"Extension" db:"extension"`
	ControlType      int                `json:"ControlType" db:"control_type"`
}

type DeviceVO struct {
	DeviceId         uint64  `json:"deviceId,string" db:"device_id"`
	DeviceGroupId    uint64  `json:"deviceGroupId,string" db:"device_group_id"`
	DeviceClassId    uint64  `json:"deviceClassId,string" db:"device_class_id"`
	DeviceProtocolId uint64  `json:"deviceProtocolId,string" db:"device_protocol_id"`
	DeviceBuildingId uint64  `json:"deviceBuildingId,string" db:"device_building_id"`
	Name             *string `json:"Name" db:"name"`
	Number           *string `json:"Number" db:"number"`
	Type             *string `json:"Type" db:"type"`
	Action           string  `json:"Action" db:"action"`
	Extension        string  `json:"Extension" db:"extension"`
	ControlType      int     `json:"ControlType" db:"control_type"`
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
		extension, _ := json.Marshal(device.Extension)
		vo.Extension = string(extension)
	}
	return &vo
}

type DeviceValue struct {
	DeviceId         uint64            `json:"deviceId,string"`
	DeviceGroupId    uint64            `json:"deviceGroupId,string"`
	DeviceClassId    uint64            `json:"deviceClassId,string"`
	DeviceProtocolId uint64            `json:"deviceProtocolId,string"`
	DeviceBuildingId uint64            `json:"deviceBuildingId,string"`
	Name             string            `json:"name"`
	DeviceGroupName  string            `json:"deviceGroupName"`
	Number           string            `json:"number"`
	Type             string            `json:"type"`
	Action           []string          `json:"action"`
	Extension        map[string]string `json:"extension"`
	ControlType      int               `json:"controlType"`
	CreateUserName   string            `json:"createUserName"`
	UpdateUserName   string            `json:"updateUserName"`
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
