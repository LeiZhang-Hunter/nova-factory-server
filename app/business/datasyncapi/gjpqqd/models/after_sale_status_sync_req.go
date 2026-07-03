package models

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AfterSaleStatusSyncReq ERP售后状态回写请求，实现 event.ZAfterSaleStatusSyncReqEvent。
type AfterSaleStatusSyncReq struct {
	Tid    string `form:"tid"`
	Rtid   string `form:"rtid"`
	Status string `form:"status"`

	db    *gorm.DB
	cache cache.Cache
	ctx   *gin.Context
}

func (r *AfterSaleStatusSyncReq) GetRtid() string   { return r.Rtid }
func (r *AfterSaleStatusSyncReq) GetTid() string    { return r.Tid }
func (r *AfterSaleStatusSyncReq) GetStatus() string { return r.Status }

// Event
func (r *AfterSaleStatusSyncReq) Config() config.Config { return nil }
func (r *AfterSaleStatusSyncReq) Action() event.EventType {
	return event.EventType("after_sale_status_sync")
}
func (r *AfterSaleStatusSyncReq) GetCache() cache.Cache       { return r.cache }
func (r *AfterSaleStatusSyncReq) GetCallback() event.Callback { return nil }
func (r *AfterSaleStatusSyncReq) GetDB() *gorm.DB             { return r.db }
func (r *AfterSaleStatusSyncReq) GetTransaction() bool        { return false }
func (r *AfterSaleStatusSyncReq) GetCtx() *gin.Context        { return r.ctx }

// Base
func (r *AfterSaleStatusSyncReq) Metadata() map[string]any { return nil }
func (r *AfterSaleStatusSyncReq) Ptr() any                 { return r }

// TransactionEvent
func (r *AfterSaleStatusSyncReq) WithDB(tx *gorm.DB)                          { r.db = tx }
func (r *AfterSaleStatusSyncReq) ToEvent() event.ZAfterSaleStatusSyncReqEvent { return r }
