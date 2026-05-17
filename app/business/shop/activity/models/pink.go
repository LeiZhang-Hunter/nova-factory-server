package models

import "nova-factory-server/app/baize"

type Pink struct {
	ID               int64   `json:"id,string" gorm:"id"`
	UID              int64   `json:"uid,string" gorm:"uid"`
	Nickname         string  `json:"nickname" gorm:"nickname"`
	Avatar           string  `json:"avatar" gorm:"avatar"`
	OrderID          string  `json:"orderId" gorm:"order_id"`
	OrderIDKey       int64   `json:"orderIdKey,string" gorm:"order_id_key"`
	TotalNum         int64   `json:"totalNum" gorm:"total_num"`
	TotalPrice       float64 `json:"totalPrice" gorm:"total_price"`
	CID              int64   `json:"cid,string" gorm:"cid"`
	PID              int64   `json:"pid,string" gorm:"pid"`
	People           int64   `json:"people" gorm:"people"`
	Price            float64 `json:"price" gorm:"price"`
	AddTime          string  `json:"addTime" gorm:"add_time"`
	StopTime         string  `json:"stopTime" gorm:"stop_time"`
	KID              int64   `json:"kId,string" gorm:"k_id"`
	IsTpl            int32   `json:"isTpl" gorm:"is_tpl"`
	IsRefund         int32   `json:"isRefund" gorm:"is_refund"`
	Status           int32   `json:"status" gorm:"status"`
	IsVirtual        int32   `json:"isVirtual" gorm:"is_virtual"`
	CombinationTitle string  `json:"combinationTitle" gorm:"combination_title"`
	CombinationImage string  `json:"combinationImage" gorm:"combination_image"`
	DeptID           int64   `json:"deptId" gorm:"dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"state"`
}

type PinkQuery struct {
	OrderID  string `form:"orderId"`
	CID      int64  `form:"cid"`
	UID      int64  `form:"uid"`
	Status   *int32 `form:"status"`
	IsRefund *int32 `form:"isRefund"`
	Page     int64  `form:"page"`
	Size     int64  `form:"size"`
}

type PinkListData struct {
	Rows  []*Pink `json:"rows"`
	Total int64   `json:"total"`
}
