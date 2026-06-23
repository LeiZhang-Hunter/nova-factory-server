package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"
)

// OrderStatusSyncEvent 订单状态同步事件，实现 event.ZOrderStatusSyncReqEvent 接口。
type OrderStatusSyncEvent struct {
	db          *gorm.DB
	cache       cache.Cache
	callback    event.Callback
	metadata    map[string]any
	cfg         config.Config
	action      event.EventType
	transaction bool
	orders      []OrderStatusSyncReqData
	ctx         *gin.Context
}

func NewOrderStatusSyncEvent(orders []*Order, RefundStatus string) *OrderStatusSyncEvent {
	data := make([]OrderStatusSyncReqData, 0, len(orders))
	for _, o := range orders {
		data = append(data, OrderStatusSyncReqData{
			DBID:         int64(o.ID),
			Tid:          o.Tid,
			Status:       o.Status,
			RefundStatus: RefundStatus,
		})
	}
	return &OrderStatusSyncEvent{
		orders: data,
	}
}

type OrderStatusSyncReqData struct {
	DBID         int64  `json:"dbid"`
	Tid          string `json:"tid"`
	Status       string `json:"status"`
	RefundStatus string `json:"refundstatus"`
}

func (d *OrderStatusSyncReqData) GetDBID() int64          { return d.DBID }
func (d *OrderStatusSyncReqData) GetTid() string          { return d.Tid }
func (d *OrderStatusSyncReqData) GetStatus() string       { return d.Status }
func (d *OrderStatusSyncReqData) GetRefundstatus() string { return d.RefundStatus }

// -- TransactionEvent --

func (e *OrderStatusSyncEvent) GetDB() *gorm.DB    { return e.db }
func (e *OrderStatusSyncEvent) WithDB(tx *gorm.DB) { e.db = tx }
func (e *OrderStatusSyncEvent) WithCache(c cache.Cache) {
	e.cache = c
}
func (e *OrderStatusSyncEvent) ToEvent() event.ZOrderStatusSyncReqEvent { return e }

// -- Event --

func (e *OrderStatusSyncEvent) Config() config.Config { return e.cfg }

func (e *OrderStatusSyncEvent) WithConfig(cfg config.Config) {
	e.cfg = cfg
}

func (e *OrderStatusSyncEvent) Action() event.EventType     { return e.action }
func (e *OrderStatusSyncEvent) GetCache() cache.Cache       { return e.cache }
func (e *OrderStatusSyncEvent) GetCallback() event.Callback { return e.callback }

func (e *OrderStatusSyncEvent) WithCallback(f event.Callback) {
	e.callback = f
}

func (e *OrderStatusSyncEvent) GetTransaction() bool { return e.transaction }
func (e *OrderStatusSyncEvent) WithTransaction(transaction bool) {
	e.transaction = transaction
}
func (e *OrderStatusSyncEvent) GetCtx() *gin.Context     { return e.ctx }
func (e *OrderStatusSyncEvent) WithCtx(ctx *gin.Context) { e.ctx = ctx }

// -- Base --

func (e *OrderStatusSyncEvent) Metadata() map[string]any {
	if e.metadata == nil {
		return make(map[string]any)
	}
	return e.metadata
}
func (e *OrderStatusSyncEvent) Ptr() any { return e }

// -- ZOrderStatusSyncReqEvent --

func (e *OrderStatusSyncEvent) GetOrders() []event.ZOrderStatusSyncReqData {
	if e.orders == nil {
		return make([]event.ZOrderStatusSyncReqData, 0)
	}
	result := make([]event.ZOrderStatusSyncReqData, len(e.orders))
	for i := range e.orders {
		result[i] = &e.orders[i]
	}
	return result
}

// -- Builder --

func (e *OrderStatusSyncEvent) WithOrders(orders []OrderStatusSyncReqData) *OrderStatusSyncEvent {
	e.orders = orders
	return e
}

func (e *OrderStatusSyncEvent) WithMetadata(m map[string]any) *OrderStatusSyncEvent {
	e.metadata = m
	return e
}
