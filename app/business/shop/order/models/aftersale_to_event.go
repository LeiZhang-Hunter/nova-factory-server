package models

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AftersaleSyncEvent 售后同步事件，实现 event.ZAfterSaleOrderSyncReqEvent 接口。
type AftersaleSyncEvent struct {
	db          *gorm.DB
	cache       cache.Cache
	callback    event.Callback
	metadata    map[string]any
	cfg         config.Config
	action      event.EventType
	transaction bool
	orders      []AftersaleSyncReqData
	ctx         *gin.Context
	userId      int64
}

// AftersaleSyncReqData 单条售后同步数据，实现 event.ZAfterSaleOrderSyncReqData。
type AftersaleSyncReqData struct {
	Rtid           string                         `json:"rtid"`
	Tid            string                         `json:"tid"`
	Total          float64                        `json:"total"`
	Privilege      float64                        `json:"privilege"`
	PostFee        float64                        `json:"post_fee"`
	Created        string                         `json:"created"`
	AftSaleType    string                         `json:"aft_sale_type"`
	ReasonCode     string                         `json:"reason_code"`
	LogistBillCode string                         `json:"logist_bill_code"`
	AftSaleRemark  string                         `json:"aft_sale_remark"`
	Details        []AftersaleSyncDetailWrapper   `json:"details"`
	ExDetails      []AftersaleSyncExDetailWrapper `json:"ex_details"`
}

// AftersaleSyncDetailWrapper 售后退款明细。
type AftersaleSyncDetailWrapper struct {
	Oid            string  `json:"oid"`
	EshopGoodsName string  `json:"eshop_goods_name"`
	EshopSkuName   string  `json:"eshop_sku_name"`
	BackQty        float64 `json:"back_qty"`
	BackTotal      float64 `json:"back_total"`
	OuterIid       string  `json:"outer_iid"`
}

// AftersaleSyncExDetailWrapper 售后换货明细。
type AftersaleSyncExDetailWrapper struct {
	Oid            string  `json:"oid"`
	EshopGoodsName string  `json:"eshop_goods_name"`
	EshopSkuName   string  `json:"eshop_sku_name"`
	ExchangeQty    float64 `json:"exchange_qty"`
	BackTotal      float64 `json:"back_total"`
	OuterIid       string  `json:"outer_iid"`
}

// --- ZAfterSaleOrderSyncReqData interface ---

func (d AftersaleSyncReqData) GetRtid() string        { return d.Rtid }
func (d AftersaleSyncReqData) GetTid() string         { return d.Tid }
func (d AftersaleSyncReqData) GetTotal() float64      { return d.Total }
func (d AftersaleSyncReqData) GetPrivilege() float64  { return d.Privilege }
func (d AftersaleSyncReqData) GetPostFee() float64    { return d.PostFee }
func (d AftersaleSyncReqData) GetCreated() string     { return d.Created }
func (d AftersaleSyncReqData) GetAftSaleType() string { return d.AftSaleType }

func (d AftersaleSyncReqData) GetReasonCode() string     { return d.ReasonCode }
func (d AftersaleSyncReqData) GetLogistBillCode() string { return d.LogistBillCode }
func (d AftersaleSyncReqData) GetAftSaleRemark() string  { return d.AftSaleRemark }

func (d AftersaleSyncReqData) GetDetails() *[]event.ZAfterSaleOrderDetail {
	if len(d.Details) == 0 {
		return nil
	}
	result := make([]event.ZAfterSaleOrderDetail, len(d.Details))
	for i := range d.Details {
		result[i] = d.Details[i]
	}
	return &result
}

func (d AftersaleSyncReqData) GetExDetails() *[]event.ZAfterSaleOrderExDetail {
	if len(d.ExDetails) == 0 {
		return nil
	}
	result := make([]event.ZAfterSaleOrderExDetail, len(d.ExDetails))
	for i := range d.ExDetails {
		result[i] = d.ExDetails[i]
	}
	return &result
}

// --- ZAfterSaleOrderDetail interface ---

func (d AftersaleSyncDetailWrapper) GetOid() string            { return d.Oid }
func (d AftersaleSyncDetailWrapper) GetEshopGoodsName() string { return d.EshopGoodsName }
func (d AftersaleSyncDetailWrapper) GetEshopSkuName() string   { return d.EshopSkuName }
func (d AftersaleSyncDetailWrapper) GetBackQty() float64       { return d.BackQty }
func (d AftersaleSyncDetailWrapper) GetBackTotal() float64     { return d.BackTotal }
func (d AftersaleSyncDetailWrapper) GetOuterIid() string       { return d.OuterIid }

// --- ZAfterSaleOrderExDetail interface ---

func (d AftersaleSyncExDetailWrapper) GetOid() string            { return d.Oid }
func (d AftersaleSyncExDetailWrapper) GetEshopGoodsName() string { return d.EshopGoodsName }
func (d AftersaleSyncExDetailWrapper) GetEshopSkuName() string   { return d.EshopSkuName }
func (d AftersaleSyncExDetailWrapper) GetExchangeQty() float64   { return d.ExchangeQty }
func (d AftersaleSyncExDetailWrapper) GetBackTotal() float64     { return d.BackTotal }
func (d AftersaleSyncExDetailWrapper) GetOuterIid() string       { return d.OuterIid }

// --- AftersaleSyncEvent methods ---

func NewAftersaleSyncEvent(aftersale *OrderRefund, order *Order) *AftersaleSyncEvent {
	data := buildAftersaleSyncData(aftersale, order)
	return &AftersaleSyncEvent{
		orders: []AftersaleSyncReqData{data},
	}
}

func buildAftersaleSyncData(aftersale *OrderRefund, order *Order) AftersaleSyncReqData {
	created := ""
	if aftersale.CreateTime != nil {
		created = aftersale.CreateTime.Format("2006-01-02 15:04:05")
	}

	details := make([]AftersaleSyncDetailWrapper, 0, len(order.Details))
	for _, d := range order.Details {
		details = append(details, AftersaleSyncDetailWrapper{
			Oid:            d.OID,
			EshopGoodsName: d.EShopGoodsName,
			EshopSkuName:   d.EShopSkuName,
			BackQty:        d.Num,
			BackTotal:      d.Payment,
			OuterIid:       d.OuterIID,
		})
	}

	return AftersaleSyncReqData{
		Rtid:           aftersale.OutRefundNo,
		Tid:            aftersale.Tid,
		Total:          order.Total,
		Privilege:      order.Privilege,
		PostFee:        order.PostFee,
		Created:        created,
		AftSaleType:    "RefundAndGoods",
		ReasonCode:     "01",
		AftSaleRemark:  aftersale.Reason,
		LogistBillCode: "",
		Details:        details,
	}
}

// -- TransactionEvent --
func (e *AftersaleSyncEvent) GetDB() *gorm.DB                            { return e.db }
func (e *AftersaleSyncEvent) WithDB(tx *gorm.DB)                         { e.db = tx }
func (e *AftersaleSyncEvent) GetUserId() int64                           { return e.userId }
func (e *AftersaleSyncEvent) WithUserId(id int64)                        { e.userId = id }
func (e *AftersaleSyncEvent) ToEvent() event.ZAfterSaleOrderSyncReqEvent { return e }

// -- Event --
func (e *AftersaleSyncEvent) Config() config.Config         { return e.cfg }
func (e *AftersaleSyncEvent) WithConfig(cfg config.Config)  { e.cfg = cfg }
func (e *AftersaleSyncEvent) Action() event.EventType       { return e.action }
func (e *AftersaleSyncEvent) GetCache() cache.Cache         { return e.cache }
func (e *AftersaleSyncEvent) WithCache(c cache.Cache)       { e.cache = c }
func (e *AftersaleSyncEvent) GetCallback() event.Callback   { return e.callback }
func (e *AftersaleSyncEvent) WithCallback(f event.Callback) { e.callback = f }
func (e *AftersaleSyncEvent) GetTransaction() bool          { return e.transaction }
func (e *AftersaleSyncEvent) WithTransaction(b bool)        { e.transaction = b }
func (e *AftersaleSyncEvent) GetCtx() *gin.Context          { return e.ctx }
func (e *AftersaleSyncEvent) WithCtx(ctx *gin.Context)      { e.ctx = ctx }

// -- Base --
func (e *AftersaleSyncEvent) Metadata() map[string]any {
	if e.metadata == nil {
		return make(map[string]any)
	}
	return e.metadata
}
func (e *AftersaleSyncEvent) Ptr() any { return e }

// -- ZAfterSaleOrderSyncReqEvent --
func (e *AftersaleSyncEvent) GetOrders() *[]event.ZAfterSaleOrderSyncReqData {
	if e.orders == nil {
		return nil
	}
	result := make([]event.ZAfterSaleOrderSyncReqData, len(e.orders))
	for i := range e.orders {
		result[i] = e.orders[i]
	}
	return &result
}

// -- Builder --
func (e *AftersaleSyncEvent) WithMetadata(m map[string]any) *AftersaleSyncEvent {
	e.metadata = m
	return e
}

func (e *AftersaleSyncEvent) WithOrders(orders []AftersaleSyncReqData) *AftersaleSyncEvent {
	e.orders = orders
	return e
}
