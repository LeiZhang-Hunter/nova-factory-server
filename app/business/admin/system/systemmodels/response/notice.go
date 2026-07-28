package response

import (
	"nova-factory-server/app/baize"
	"time"
)

type SysNoticeVo struct {
	Id         int64       `json:"id,string" db:"id"`                   //通知ID
	Title      string      `json:"title" db:"title" binding:"required"` //通知标题
	Txt        string      `json:"txt" db:"txt" binding:"required"`     //通知文本
	Type       string      `json:"type" db:"type" binding:"required"`   //通知文本
	DeptId     int64       `json:"-" db:"dept_id"`                      //部门ID
	DeptIds    *baize.List `json:"deptIds" db:"dept_ids"`               //接收部门列表
	CreateName string      `json:"createName" db:"create_name"`         //创建人
	baize.BaseEntity
}

type ConsumptionNoticeVo struct {
	Id         int64      `json:"id,string" db:"id"`  //通知ID
	Title      string     `json:"title" db:"title" `  //通知标题
	Txt        string     `json:"txt" db:"txt" `      //通知文本
	Status     string     `json:"status" db:"status"` //通知状态 1未读 2 已读
	Type       string     `json:"type" db:"type" binding:"required"`
	CreateName string     `json:"createName" db:"create_name"`
	CreateTime *time.Time `json:"createTime" db:"create_time"` //创建时间
}
