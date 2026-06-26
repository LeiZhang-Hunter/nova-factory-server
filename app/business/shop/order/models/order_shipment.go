package models

import (
	"nova-factory-server/app/baize"
)

// OrderShipment 订单发货物流记录。
// 每次 selfmall.order.send 通知到达时插入一条，支持同一订单多次拆单、多物流公司的场景。
type OrderShipment struct {
	ID          uint64  `json:"id,string" gorm:"column:id"`
	OrderID     uint64  `json:"order_id,string" gorm:"column:order_id"`
	Tid         string  `json:"tid" gorm:"column:tid"`
	Issplit     int     `json:"issplit" gorm:"column:issplit"`
	Outsid      string  `json:"outsid" gorm:"column:outsid"`
	Companycode string  `json:"companycode" gorm:"column:companycode"`
	SubTid      string  `json:"subtid" gorm:"column:subtid"`
	OID         string  `json:"oid" gorm:"column:oid"`
	Qty         float64 `json:"qty" gorm:"column:qty"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// OrderShipmentSet 订单发货物流记录保存参数。
type OrderShipmentSet struct {
	OrderID     uint64  `json:"order_id,string"`
	Tid         string  `json:"tid"`
	Issplit     int     `json:"issplit"`
	Outsid      string  `json:"outsid"`
	Companycode string  `json:"companycode"`
	SubTid      string  `json:"subtid"`
	OID         string  `json:"oid"`
	Qty         float64 `json:"qty"`
}
