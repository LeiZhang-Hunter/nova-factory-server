package models

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/utils/observer/integration/event"
	"time"
)

// OrderSend 订单发货主表
type OrderSend struct {
	ID          uint64             `json:"id,string" gorm:"column:id"`
	OrderID     uint64             `json:"order_id,string" gorm:"column:order_id"`
	Tid         string             `json:"tid" gorm:"column:tid"`
	IsSplit     int32              `json:"is_split" gorm:"column:is_split"`
	Subtid      string             `json:"subtid" gorm:"column:subtid"`
	CompanyCode string             `json:"company_code" gorm:"column:company_code"`
	Outsid      string             `json:"outsid" gorm:"column:outsid"`
	SyncStatus  int32              `json:"sync_status" gorm:"column:sync_status"`
	SyncMessage string             `json:"sync_message" gorm:"column:sync_message"`
	SyncTime    *time.Time         `json:"sync_time" gorm:"column:sync_time"`
	ExtJSON     string             `json:"ext_json" gorm:"column:ext_json"`
	DeptID      int64              `json:"dept_id" gorm:"column:dept_id"`
	Details     []*OrderSendDetail `json:"details" gorm:"-"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// TableName 表名
func (OrderSend) TableName() string {
	return "shop_order_send"
}

// OrderSendDetail 订单发货明细表
type OrderSendDetail struct {
	ID      uint64  `json:"id,string" gorm:"column:id"`
	SendID  uint64  `json:"send_id,string" gorm:"column:send_id"`
	OrderID uint64  `json:"order_id,string" gorm:"column:order_id"`
	Tid     string  `json:"tid" gorm:"column:tid"`
	Oid     string  `json:"oid" gorm:"column:oid"`
	Qty     float64 `json:"qty" gorm:"column:qty"`
	ExtJSON string  `json:"ext_json" gorm:"column:ext_json"`
	DeptID  int64   `json:"dept_id" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// TableName 表名
func (OrderSendDetail) TableName() string {
	return "shop_order_send_detail"
}

// ToOrderSendByEvent 将发货事件转换为 OrderSend 模型。
func ToOrderSendByEvent(e event.OrderSendEvent) *OrderSend {
	if e == nil {
		return nil
	}
	return &OrderSend{
		Tid:         e.GetTid(),
		CompanyCode: e.GetCompanyCode(),
		IsSplit:     int32(e.GetIsSplit()),
		Outsid:      e.GetOutSid(),
		Subtid:      e.GetSubTid(),
		Details:     toOrderSendDetails(e.GetTid(), e.GetDetails()),
	}
}

// toOrderSendDetails 将事件发货明细列表转换为模型明细列表。
func toOrderSendDetails(tid string, eventDetails []event.OrderSendDetail) []*OrderSendDetail {
	if len(eventDetails) == 0 {
		return nil
	}
	details := make([]*OrderSendDetail, 0, len(eventDetails))
	for _, d := range eventDetails {
		if d == nil {
			continue
		}
		details = append(details, &OrderSendDetail{
			Tid: tid,
			Oid: d.GetSubTid(),
			Qty: float64(d.GetQty()),
		})
	}
	return details
}
