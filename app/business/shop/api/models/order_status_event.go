package models

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"

	"gorm.io/gorm"
)

type OrderStatusEvent struct {
	db       *gorm.DB       `json:"-" form:"-"`
	cache    cache.Cache    `json:"-"`
	callback event.Callback `json:"-"`
	metadata map[string]any
}

func (o *OrderStatusEvent) WithDB(tx *gorm.DB) {
	o.db = tx
}

func (o *OrderStatusEvent) Config() config.Config {
	return nil
}

func (o *OrderStatusEvent) Action() event.EventType {
	return ""
}

func (o *OrderStatusEvent) GetCache() cache.Cache {
	return o.cache
}

func (o *OrderStatusEvent) GetCallback() event.Callback {
	return o.callback
}

func (o *OrderStatusEvent) GetTransaction() bool {
	return true
}

func (o *OrderStatusEvent) Metadata() map[string]any {
	return o.metadata
}
func (o *OrderStatusEvent) WithMetadata(metadata map[string]any) {
	o.metadata = metadata
}

func (o *OrderStatusEvent) Ptr() any {
	return o
}

func (o *OrderStatusEvent) Orders() []event.OrderStratusEventData {
	return nil
}

func (o *OrderStatusEvent) GetDB() *gorm.DB {
	return o.db
}

func (o *OrderStatusEvent) ToEvent() event.OrderStratusEvent {
	return o
}

type OrderStatusEventData struct {
	tid          string
	status       string
	refundstatus string
}

func NewOrderStatusEventData(tid, status, refundstatus string) *OrderStatusEventData {
	return &OrderStatusEventData{tid: tid, status: status, refundstatus: refundstatus}
}
func (o *OrderStatusEventData) GetTid() string {
	return o.tid
}
func (o *OrderStatusEventData) GetStatus() string {
	return o.status
}
func (o *OrderStatusEventData) GetRefundstatus() string {
	return o.refundstatus
}
