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
	RtspUrl   string `gorm:"column:rtsp_url" json:"rtsp_url"`
	Status    int    `gorm:"column:status" json:"status"`
	DeptId    int64  `gorm:"column:dept_id" json:"dept_id"`
	baize.BaseEntity
	State int `gorm:"column:state" json:"state"`
}

type IotCameraListReq struct {
	Name      string `gorm:"column:name" json:"name"`
	IpAddress string `gorm:"column:ip_address" json:"ip_address"`
	Brand     string `gorm:"column:brand" json:"brand"`
}

type IotCameraList struct {
	Rows  []*IotCamera `json:"rows"`
	Total int64        `json:"total"`
}
