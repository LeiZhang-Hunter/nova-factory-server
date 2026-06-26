package models

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderSyncReq 同步订单
type OrderSyncReq struct {
	Orders   []OrderItem    `json:"orders" form:"orders"`
	cache    cache.Cache    `json:"-"`
	callback event.Callback `json:"-"`
}

func (o *OrderSyncReq) GetOrders() []event.OrderData {
	orders := make([]event.OrderData, 0, len(o.Orders))
	for i := range o.Orders {
		orders = append(orders, &o.Orders[i])
	}
	return orders
}

func (o *OrderSyncReq) Metadata() map[string]any {
	return map[string]any{}
}

func (o *OrderSyncReq) Ptr() any {
	return o
}

func (o *OrderSyncReq) Config() config.Config {
	return nil
}

func (o *OrderSyncReq) Action() event.EventType {
	return event.EventOrderCreated
}

func (o *OrderSyncReq) Cache() cache.Cache {
	return o.cache
}

func (o *OrderSyncReq) GetCallback() event.Callback {
	return o.callback
}

// OrderItem 单条订单数据，实现 event.OrderData 接口
type OrderItem struct {
	OrderNo          string  `json:"orderno" form:"orderno"`
	Weight           float64 `json:"weight" form:"weight"`
	Size             float64 `json:"size" form:"size"`
	BuyerNick        string  `json:"buyernick" form:"buyernick"`
	BuyerMessage     string  `json:"buyermessage" form:"buyermessage"`
	SellerMemo       string  `json:"sellermemo" form:"sellermemo"`
	TotalAmount      float64 `json:"totalamount" form:"totalamount"`
	Privilege        float64 `json:"privilege" form:"privilege"`
	PostFee          float64 `json:"postfee" form:"postfee"`
	ReceiverName     string  `json:"receivername" form:"receivername"`
	ReceiverState    string  `json:"receiverstat" form:"receiverstat"`
	ReceiverCity     string  `json:"receivercity" form:"receivercity"`
	ReceiverDistrict string  `json:"receiverdistrict" form:"receiverdistrict"`
	ReceiverAddress  string  `json:"receiveraddress" form:"receiveraddress"`
	ReceiverPhone    string  `json:"receiverphone" form:"receiverphone"`
	ReceiverMobile   string  `json:"receivermobile" form:"receivermobile"`
	Created          string  `json:"created" form:"created"`
	Type             string  `json:"type" form:"type"`
	Status           string  `json:"status" form:"status"`
	InvoiceName      string  `json:"invoicename" form:"invoicename"`
	SellerFlag       string  `json:"sellerflag" form:"sellerflag"`
	PayTime          string  `json:"paytime" form:"paytime"`
	LogistBTypeCode  string  `json:"logistbtypecode" form:"logistbtypecode"`
	LogistBillCode   string  `json:"logistbillcode" form:"logistbillcode"`
	BTypeCode        string  `json:"btypecode" form:"btypecode"`
	db               *gorm.DB
}

func (o *OrderItem) GetTransactionId() string {
	//TODO implement me
	panic("implement me")
}

func (o *OrderItem) GetNotifyRaw() string {
	return ""
}

func (o *OrderItem) GetMchId() string {
	return ""
}

func (o *OrderItem) GetAppid() string {
	return ""
}

func (o *OrderItem) GetPayerOpenid() string {
	return ""
}

func (o *OrderItem) GetPayChannel() int {
	return 0
}

func (o *OrderItem) GetUserId() uint64 {
	return 0
}

func (o *OrderItem) GetOrderNo() string          { return o.OrderNo }
func (o *OrderItem) GetWeight() float64          { return o.Weight }
func (o *OrderItem) GetSize() float64            { return o.Size }
func (o *OrderItem) GetBuyerNick() string        { return o.BuyerNick }
func (o *OrderItem) GetBuyerMessage() string     { return o.BuyerMessage }
func (o *OrderItem) GetSellerMemo() string       { return o.SellerMemo }
func (o *OrderItem) GetTotalAmount() float64     { return o.TotalAmount }
func (o *OrderItem) GetPrivilege() float64       { return o.Privilege }
func (o *OrderItem) GetPostFee() float64         { return o.PostFee }
func (o *OrderItem) GetReceiverName() string     { return o.ReceiverName }
func (o *OrderItem) GetReceiverState() string    { return o.ReceiverState }
func (o *OrderItem) GetReceiverCity() string     { return o.ReceiverCity }
func (o *OrderItem) GetReceiverDistrict() string { return o.ReceiverDistrict }
func (o *OrderItem) GetReceiverAddress() string  { return o.ReceiverAddress }
func (o *OrderItem) GetReceiverPhone() string    { return o.ReceiverPhone }
func (o *OrderItem) GetReceiverMobile() string   { return o.ReceiverMobile }
func (o *OrderItem) GetCreated() string          { return o.Created }
func (o *OrderItem) GetType() string             { return o.Type }
func (o *OrderItem) GetStatus() string           { return o.Status }
func (o *OrderItem) GetInvoiceName() string      { return o.InvoiceName }
func (o *OrderItem) GetSellerFlag() string       { return o.SellerFlag }
func (o *OrderItem) GetPayTime() string          { return o.PayTime }
func (o *OrderItem) GetLogIstBTypeCode() string  { return o.LogistBTypeCode }
func (o *OrderItem) GetLogIstBillCode() string   { return o.LogistBillCode }
func (o *OrderItem) GetBTypeCode() string        { return o.BTypeCode }

func (o *OrderItem) GetDetails() []event.GoodsDetail {
	return nil
}

func (o *OrderItem) GetAccounts() []event.Account {
	return nil
}

func (o *OrderItem) Metadata() map[string]any {
	return map[string]any{}
}

func (o *OrderItem) Ptr() any {
	return o
}

// OrderSendReq 订单发货
type OrderSendReq struct {
	Tid         string            `json:"tid" form:"tid"`
	Companycode string            `json:"companycode" form:"companycode"`
	Issplit     int               `json:"issplit" form:"issplit"`
	Outsid      string            `json:"outsid" form:"outsid"`
	Subtid      string            `json:"subtid" form:"subtid"`
	Details     []OrderSendDetail `json:"details" form:"details"`
	DB          *gorm.DB          `json:"-" form:"-"`
	ctx         *gin.Context
}

func (o *OrderSendReq) GetDetails() []event.OrderSendDetail {
	if len(o.Details) == 0 {
		return make([]event.OrderSendDetail, 0)
	}
	details := make([]event.OrderSendDetail, len(o.Details))
	for _, detail := range o.Details {
		details = append(details, &detail)
	}
	return details
}

func (o *OrderSendReq) GetTid() string {
	return o.Tid
}
func (o *OrderSendReq) GetCompanyCode() string {
	return o.Companycode
}
func (o *OrderSendReq) GetIsSplit() int {
	return o.Issplit
}
func (o *OrderSendReq) GetOutSid() string {
	return o.Outsid
}
func (o *OrderSendReq) GetSubTid() string {
	return o.Subtid
}

func (o *OrderSendReq) WithDB(tx *gorm.DB) {
	o.DB = tx
}
func (o *OrderSendReq) ToEvent() event.OrderSendEvent {
	return o
}

func (o *OrderSendReq) GetCtx() *gin.Context {
	return o.ctx
}

// Config 返回本次事件关联的集成配置，可能为 nil（表示未配置集成）
func (o *OrderSendReq) Config() config.Config {
	return nil
}

// Metadata 返回业务数据关联的扩展元数据键值对
func (o *OrderSendReq) Metadata() map[string]any {
	return make(map[string]any)
}

// Ptr 返回当前数据对象的指针，供序列化或深拷贝使用
func (o *OrderSendReq) Ptr() any {
	return o
}

// Action 返回事件类型（创建、更新、删除等），观察者据此决定同步策略
func (o *OrderSendReq) Action() event.EventType {
	return event.EventOrderSendChanged
}

// Cache 返回缓存实例，用于在同步过程中读写临时数据（如 OAuth Token）
func (o *OrderSendReq) GetCache() cache.Cache {
	return nil
}

// GetCallback 返回处理完成后的回调接口，用于通知上游同步结果
func (o *OrderSendReq) GetCallback() event.Callback {
	return nil
}

// GetDB 读取DB
func (o *OrderSendReq) GetDB() *gorm.DB {
	return nil
}

// GetTransaction 是否打开事物
func (o *OrderSendReq) GetTransaction() bool {
	return false
}

type OrderSendDetail struct {
	SubTid string `json:"subtid" form:"subtid"`
	Qty    int    `json:"qty" form:"qty"`
}

func (o *OrderSendDetail) GetSubTid() string {
	return o.SubTid
}

// GetQty 读取库存
func (o *OrderSendDetail) GetQty() int {
	return o.Qty
}

// OrderSendResponse 订单发货返回接口
type OrderSendResponse struct {
	Iserror  bool   `json:"iserror"`
	Errormsg string `json:"errormsg"`
	Tid      string `json:"tid"`
	Issplit  int    `json:"issplit"`
	Subtids  string `json:"subtids"`
}
