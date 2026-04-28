package models

import "nova-factory-server/app/baize"

type Pink struct {
	ID               int64   `json:"id,string" db:"id"`
	UID              int64   `json:"uid,string" db:"uid"`
	Nickname         string  `json:"nickname" db:"nickname"`
	Avatar           string  `json:"avatar" db:"avatar"`
	OrderID          string  `json:"orderId" db:"order_id"`
	OrderIDKey       int64   `json:"orderIdKey,string" db:"order_id_key"`
	TotalNum         int64   `json:"totalNum" db:"total_num"`
	TotalPrice       float64 `json:"totalPrice" db:"total_price"`
	CID              int64   `json:"cid,string" db:"cid"`
	PID              int64   `json:"pid,string" db:"pid"`
	People           int64   `json:"people" db:"people"`
	Price            float64 `json:"price" db:"price"`
	AddTime          string  `json:"addTime" db:"add_time"`
	StopTime         string  `json:"stopTime" db:"stop_time"`
	KID              int64   `json:"kId,string" db:"k_id"`
	IsTpl            int32   `json:"isTpl" db:"is_tpl"`
	IsRefund         int32   `json:"isRefund" db:"is_refund"`
	Status           int32   `json:"status" db:"status"`
	IsVirtual        int32   `json:"isVirtual" db:"is_virtual"`
	CombinationTitle string  `json:"combinationTitle" db:"combination_title"`
	CombinationImage string  `json:"combinationImage" db:"combination_image"`
	DeptID           int64   `json:"deptId" db:"dept_id"`
	baize.BaseEntity
	State int32 `json:"state" db:"state"`
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
