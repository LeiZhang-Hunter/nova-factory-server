package models

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderStatusData 订单状态变更数据，实现 event.OrderStratusEventData 接口。
type OrderStatusData struct {
	Tid          string `json:"tid"`
	Status       string `json:"status"`
	RefundStatus string `json:"refund_status"`
}

//GetTid() string
//GetStatus() string
//GetRefundstatus() string

func NewOrderStatusData(tid, status, refundstatus string) []event.ZOrderStatusSyncReqData {
	return []event.ZOrderStatusSyncReqData{
		&OrderStatusData{
			Tid:          tid,
			Status:       status,
			RefundStatus: refundstatus,
		},
	}
}
func (d *OrderStatusData) GetTid() string          { return d.Tid }
func (d *OrderStatusData) GetStatus() string       { return d.Status }
func (d *OrderStatusData) GetRefundstatus() string { return d.RefundStatus }

// OrderStatusEvent 订单状态变更事件，同时实现 TransactionEvent[OrderStratusEvent] 和 OrderStratusEvent。
type OrderStatusEvent struct {
	db          *gorm.DB
	cache       cache.Cache
	callback    event.Callback
	metadata    map[string]any
	cfg         config.Config
	action      event.EventType
	transaction bool
	orders      []event.ZOrderStatusSyncReqData
	ctx         *gin.Context
}

// -- TransactionEvent[OrderStratusEvent] --

func (e *OrderStatusEvent) GetDB() *gorm.DB                         { return e.db }
func (e *OrderStatusEvent) WithDB(tx *gorm.DB)                      { e.db = tx }
func (e *OrderStatusEvent) ToEvent() event.ZOrderStatusSyncReqEvent { return e }

// -- Event --

func (e *OrderStatusEvent) Config() config.Config       { return e.cfg }
func (e *OrderStatusEvent) Action() event.EventType     { return e.action }
func (e *OrderStatusEvent) GetCache() cache.Cache       { return e.cache }
func (e *OrderStatusEvent) GetCallback() event.Callback { return e.callback }
func (e *OrderStatusEvent) GetTransaction() bool        { return e.transaction }
func (e *OrderStatusEvent) GetCtx() *gin.Context        { return e.ctx }
func (e *OrderStatusEvent) WithCtx(ctx *gin.Context) {
	e.ctx = ctx
}

// -- Base --

func (e *OrderStatusEvent) Metadata() map[string]any {
	if e.metadata == nil {
		return make(map[string]any)
	}
	return e.metadata
}
func (e *OrderStatusEvent) Ptr() any { return e }

// -- OrderStratusEvent --

func (e *OrderStatusEvent) GetOrders() []event.ZOrderStatusSyncReqData {
	if e.orders == nil {
		return make([]event.ZOrderStatusSyncReqData, 0)
	}
	return e.orders
}

// -- Builder --

func (e *OrderStatusEvent) WithOrders(orders []event.ZOrderStatusSyncReqData) *OrderStatusEvent {
	e.orders = orders
	return e
}

func (e *OrderStatusEvent) WithMetadata(m map[string]any) *OrderStatusEvent {
	e.metadata = m
	return e
}

func (e *OrderStatusEvent) WithCallback(f event.Callback) {
	e.callback = f
}
