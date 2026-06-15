package models

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"

	"gorm.io/gorm"
)

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
