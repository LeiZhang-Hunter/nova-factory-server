package cameraModels

import (
	"nova-factory-server/app/baize"
)

// IotCamera 监控摄像头模型
type IotCamera struct {
	Id        int64  `gorm:"primaryKey;column:id" json:"id,string"`
	Name      string `gorm:"column:name" json:"name"`
	Number    string `gorm:"column:number" json:"number"`
	IpAddress string `gorm:"column:ip_address" json:"ip_address"`
	Brand     string `gorm:"column:brand" json:"brand"`
	Port      int    `gorm:"column:port" json:"port"`
	Username  string `gorm:"column:username" json:"username"`
	Password  string `gorm:"column:password" json:"password"`
	Enable    *bool  `gorm:"column:enable" json:"enable"`
	Status    *bool  `gorm:"column:status" json:"status"`
	DeptId    int64  `gorm:"column:dept_id" json:"dept_id"`
	baize.BaseEntity
	State int `gorm:"column:state" json:"state"`
}

type IotCameraListReq struct {
	Name      string `gorm:"column:name" form:"name"`
	IpAddress string `gorm:"column:ip_address" form:"ip_address"`
	Brand     string `gorm:"column:brand" form:"brand"`
}

type IotCameraList struct {
	Rows  []*IotCamera `json:"rows"`
	Total int64        `json:"total"`
}

type IotCameraDetail struct {
	IotCamera
	GatewayRealtime map[string]interface{} `json:"gateway_realtime,omitempty"`
}
