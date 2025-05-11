package deviceModels

import (
	"github.com/gogf/gf/util/gconv"
	"nova-factory-server/app/baize"
)

type DeviceGroupDQL struct {
	Name string `form:"name"`
	baize.BaseEntityDQL
}

type DeviceGroup struct {
	GroupId uint64  `json:"groupId,string" db:"device_id"`
	Name    *string `json:"name" db:"name"`
	Remark  *string `json:"remark" db:"remark"`
}

type DeviceGroupVO struct {
	GroupId        uint64 `json:"groupId,string" db:"device_id"`
	DeptId         uint64 `json:"deptId,string" db:"dept_id"`
	Name           string `json:"name" db:"name"`
	Remark         string `json:"remark" db:"remark"`
	CreateUserName string `json:"create_user_name" db:"create_user_name" gorm:"-"`
	UpdateUserName string `json:"update_user_name" db:"update_user_name" gorm:"-"`
	baize.BaseEntity
}

type DeviceGroupListData struct {
	Rows  []*DeviceGroupVO `json:"rows"`
	Total int64            `json:"total"`
}

func NewDeviceGroupVO(group *DeviceGroup) *DeviceGroupVO {
	var vo DeviceGroupVO
	gconv.Scan(group, &vo)
	return &vo
}
