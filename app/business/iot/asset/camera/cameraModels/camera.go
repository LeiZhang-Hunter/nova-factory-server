package cameraModels

import (
	"time"
)

// IotCamera 监控摄像头模型
type IotCamera struct {
	Id         int64     `gorm:"primaryKey;column:id" json:"id,string"`
	Name       string    `gorm:"column:name" json:"name"`
	Number     string    `gorm:"column:number" json:"number"`
	IpAddress  string    `gorm:"column:ip_address" json:"ip_address"`
	Brand      string    `gorm:"column:brand" json:"brand"`
	Port       int       `gorm:"column:port" json:"port"`
	Username   string    `gorm:"column:username" json:"username"`
	Password   string    `gorm:"column:password" json:"password"`
	RtspUrl    string    `gorm:"column:rtsp_url" json:"rtsp_url"`
	Status     int       `gorm:"column:status" json:"status"`
	DeptId     int64     `gorm:"column:dept_id" json:"dept_id"`
	CreateBy   int64     `gorm:"column:create_by" json:"create_by"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateBy   int64     `gorm:"column:update_by" json:"update_by"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	State      int       `gorm:"column:state" json:"state"`
}

// TableName 设置表名
func (IotCamera) TableName() string {
	return "iot_camera"
}
