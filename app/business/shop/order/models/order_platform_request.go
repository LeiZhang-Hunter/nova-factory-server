package models

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/stringUtils"
	"nova-factory-server/app/utils/time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderSyncRequest 同步请求
type OrderSyncRequest struct {
	Orders      []*OrderSyncOrder `json:"orders"`
	cfg         config.Config     `json:"-"`
	action      event.EventType   `json:"-"`
	c           cache.Cache       `json:"-"`
	db          *gorm.DB          `json:"-"`
	transaction bool              `json:"-"`
	callback    event.Callback
	ctx         *gin.Context
}

func (o *OrderSyncRequest) WithDB(db *gorm.DB) {
	o.db = db
}

func (o *OrderSyncRequest) GetDB() *gorm.DB {
	return o.db
}
func (o *OrderSyncRequest) ToEvent() event.OrderEvent {
	return o
}
func (o *OrderSyncRequest) GetCtx() *gin.Context {
	return o.ctx
}
func (o *OrderSyncRequest) GetTransaction() bool {
	return o.transaction
}

func (o *OrderSyncRequest) WithTransaction(transaction bool) {
	o.transaction = transaction
}

func (o *OrderSyncRequest) GetCache() cache.Cache {
	return o.c
}

func (o *OrderSyncRequest) GetCallback() event.Callback {
	return o.callback
}
func (o *OrderSyncRequest) WithCallback(callback event.Callback) {
	o.callback = callback
}
func (o *OrderSyncRequest) GetOrders() []event.OrderData {
	if o.Orders == nil {
		return make([]event.OrderData, 0)
	}
	result := make([]event.OrderData, 0, len(o.Orders))
	for _, d := range o.Orders {
		result = append(result, d)
	}
	return result
}
func (o *OrderSyncRequest) Metadata() map[string]any {
	return make(map[string]any)
}
func (o *OrderSyncRequest) Ptr() any {
	return o
}

func (o *OrderSyncRequest) WithConfig(c config.Config) {
	o.cfg = c
	return
}

func (o *OrderSyncRequest) Config() config.Config {
	return o.cfg
}

func (o *OrderSyncRequest) WithAction(action event.EventType) {
	o.action = action
}

func (o *OrderSyncRequest) Action() event.EventType {
	return o.action
}

func (o *OrderSyncRequest) WithCache(c cache.Cache) {
	o.c = c
	return
}

func (o *OrderSyncRequest) Cache() cache.Cache {
	return o.c
}

type OrderSyncOrder struct {
	Tid              string              `json:"tid"`
	UserId           int64               `json:"user_id"`
	Weight           float64             `json:"weight"`
	Size             float64             `json:"size"`
	BuyerNick        string              `json:"buyernick"`
	BuyerMessage     string              `json:"buyermessage"`
	SellerMemo       string              `json:"sellermemo"`
	Total            float64             `json:"total"`
	Privilege        float64             `json:"privilege"`
	PostFee          float64             `json:"postfee"`
	ReceiverName     string              `json:"receivername"`
	ReceiverState    string              `json:"receiverstate"`
	ReceiverCity     string              `json:"receivercity"`
	ReceiverDistrict string              `json:"receiverdistrict"`
	ReceiverAddress  string              `json:"receiveraddress"`
	ReceiverPhone    string              `json:"receiverphone"`
	ReceiverMobile   string              `json:"receivermobile"`
	Created          string              `json:"created"`
	Status           string              `json:"status"`
	Type             string              `json:"type"`
	InvoiceName      string              `json:"invoicename"`
	SellerFlag       string              `json:"sellerflag"`
	PayTime          string              `json:"paytime"`
	LogistBTypeCode  string              `json:"logistbtypecode"`
	LogistBillCode   string              `json:"logistbillcode"`
	BTypeCode        string              `json:"btypecode"`
	Details          []*OrderSyncDetail  `json:"details"`
	Accounts         []*OrderSyncAccount `json:"accounts,omitempty"`
	TransactionId    string              `json:"transaction_id"`
	NotifyRaw        string              `json:"notify_raw"`
	MchId            string              `json:"mch_id"`
	Appid            string              `json:"appid"`
	PayerOpenid      string              `json:"payer_openid"`
	PayChannel       int                 `json:"pay_channel"`
}

func (o *OrderSyncOrder) Metadata() map[string]any {
	return map[string]any{}
}
func (o *OrderSyncOrder) Ptr() any {
	return o
}

// GetOrderNo 通用订单字段（所有业务模块都能提供）
func (o *OrderSyncOrder) GetOrderNo() string {
	return o.Tid
}
func (o *OrderSyncOrder) GetWeight() float64 {
	return o.Weight
}
func (o *OrderSyncOrder) GetSize() float64 {
	return o.Size
}

// 买家账号
func (o *OrderSyncOrder) GetBuyerNick() string {
	return o.BuyerNick
}

// 卖家留言
func (o *OrderSyncOrder) GetBuyerMessage() string {
	return o.BuyerMessage
}

// 卖家备注
func (o *OrderSyncOrder) GetSellerMemo() string {
	return o.SellerMemo
}

// 订单总金额
func (o *OrderSyncOrder) GetTotalAmount() float64 {
	return o.Total
}

// 订单享受优惠的金额（订单总金额-实际支付金额）
func (o *OrderSyncOrder) GetPrivilege() float64 {
	return o.Privilege
}

// 运费
func (o *OrderSyncOrder) GetPostFee() float64 {
	return o.PostFee
}

// 收货人名称
func (o *OrderSyncOrder) GetReceiverName() string {
	return o.ReceiverName
}

func (o *OrderSyncOrder) GetReceiverState() string {
	return o.ReceiverState
}

func (o *OrderSyncOrder) GetReceiverCity() string {
	return o.ReceiverCity
}
func (o *OrderSyncOrder) GetReceiverDistrict() string {
	return o.ReceiverDistrict
}
func (o *OrderSyncOrder) GetReceiverAddress() string {
	return o.ReceiverAddress
}
func (o *OrderSyncOrder) GetReceiverPhone() string {
	return o.ReceiverPhone
}
func (o *OrderSyncOrder) GetReceiverMobile() string {
	return o.ReceiverMobile
}
func (o *OrderSyncOrder) GetCreated() string {
	return o.Created
}

// 订单类型（Cod=货到付款, NoCod=非货到付款）
func (o *OrderSyncOrder) GetType() string {
	return o.Type
}

// GetStatus  订单状态 NoPay = 未付款 Payed = 已付款（货到付款传已付款） Sended = 已发货 TradeSuccess = 交易成功 TradeClosed = 交易关闭 PartSend = 部分发货
func (o *OrderSyncOrder) GetStatus() string {
	return o.Status
}

func (o *OrderSyncOrder) GetInvoiceName() string {
	return o.InvoiceName
}

// 卖家旗帜（数值型）
func (o *OrderSyncOrder) GetSellerFlag() string {
	return o.SellerFlag
}

// GetPayTime 付款时间（时间格式：yyyy-MM-dd HH:mm:ss）
func (o *OrderSyncOrder) GetPayTime() string {
	return o.PayTime
}

// GetLogIstBTypeCode 物流公司编码（填入信息需和全渠道中的物流公司编码相同才可以匹配到，如 “ZTO”）
func (o *OrderSyncOrder) GetLogIstBTypeCode() string {
	return o.LogistBTypeCode
}

// GetLogIstBillCode 物流单号
func (o *OrderSyncOrder) GetLogIstBillCode() string {
	return o.LogistBillCode
}

// GetBTypeCode 往来单位编码
func (o *OrderSyncOrder) GetBTypeCode() string {
	return o.BTypeCode
}
func (o *OrderSyncOrder) GetUserId() uint64 {
	return uint64(o.UserId)
}

// Details 订单商品信息
func (o *OrderSyncOrder) GetDetails() []event.GoodsDetail {
	if o.Details == nil {
		return nil
	}
	result := make([]event.GoodsDetail, 0, len(o.Details))
	for _, d := range o.Details {
		result = append(result, d)
	}
	return result
}

// Accounts 订单账户信息
func (o *OrderSyncOrder) GetAccounts() []event.Account {
	if o.Accounts == nil {
		return nil
	}
	result := make([]event.Account, 0, len(o.Accounts))
	for _, a := range o.Accounts {
		result = append(result, a)
	}
	return result
}

func (o *OrderSyncOrder) GetTransactionId() string {
	return o.TransactionId
}
func (o *OrderSyncOrder) GetNotifyRaw() string {
	return o.NotifyRaw
}
func (o *OrderSyncOrder) GetMchId() string {
	return o.MchId
}
func (o *OrderSyncOrder) GetAppid() string {
	return o.Appid
}
func (o *OrderSyncOrder) GetPayerOpenid() string {
	return o.PayerOpenid
}
func (o *OrderSyncOrder) GetPayChannel() int {
	return o.PayChannel
}

type OrderSyncDetail struct {
	OID            string  `json:"oid"`
	Barcode        string  `json:"barcode"`
	EShopGoodsID   string  `json:"eshopgoodsid"`
	OuterIID       string  `json:"outeriid"`
	EShopGoodsName string  `json:"eshopgoodsname"`
	EShopSKUId     string  `json:"eshopskuid"`
	EShopSKUName   string  `json:"eshopskuname"`
	NumIID         int64   `json:"numiid"`
	SKUId          int64   `json:"skuid"`
	Num            float64 `json:"num"`
	Payment        float64 `json:"payment"`
	PicPath        string  `json:"picpath"`
	Weight         float64 `json:"weight"`
	Size           float64 `json:"size"`
	UnitID         int64   `json:"unitid"`
	UnitQty        float64 `json:"unitqty"`
}

// GetOid 网店订单明细编号
func (o *OrderSyncDetail) GetOid() string {
	return o.OID
}

// GetBarcode 商品条码
func (o *OrderSyncDetail) GetBarcode() string {
	return o.Barcode
}

// GetEshopGoodsId 网店商品ID（对应商品列表接口中的productid），方案A必填
func (o *OrderSyncDetail) GetEshopGoodsId() string {
	return o.EShopGoodsID
}

// GetOuterIid 网店商家编码
func (o *OrderSyncDetail) GetOuterIid() string {
	return o.OuterIID
}

// GetEshopGoodsName 网店商品名称
func (o *OrderSyncDetail) GetEshopGoodsName() string {
	return o.EShopGoodsName
}

// GetEshopSkuId 网店商品SKUID（对应商品列表接口中的skuid），方案A且商品有SKU则必填
func (o *OrderSyncDetail) GetEshopSkuId() string {
	return o.EShopSKUId
}

// GetEshopSkuName 网店商品SKU名称，方案A且商品有SKU则必填
func (o *OrderSyncDetail) GetEshopSkuName() string {
	return o.EShopSKUName
}

// GetNumIid 商品ID，方案B必填，方案A不要传值
func (o *OrderSyncDetail) GetNumIid() int64 {
	return o.NumIID
}

// GetSkuId 	规格ID，方案B必填，方案A不要传值
func (o *OrderSyncDetail) GetSkuId() int64 {
	return o.SKUId
}

// GetNum 基本单位数量（不能为0）
func (o *OrderSyncDetail) GetNum() float64 {
	return o.Num
}

// GetPayment 商品总额
func (o *OrderSyncDetail) GetPayment() float64 {
	return o.Payment
}

// GetPicPath 商品图片路径
func (o *OrderSyncDetail) GetPicPath() string {
	return o.PicPath
}

// GetWeight 重量
func (o *OrderSyncDetail) GetWeight() float64 {
	return o.Weight
}

// GetSize 尺寸，体积
func (o *OrderSyncDetail) GetSize() float64 {
	return o.Size
}

// GetUniTid 销售单位ID
func (o *OrderSyncDetail) GetUniTid() int64 {
	return o.UnitID
}

// GetUnitQty 销售单位数量（无多单位可不填）
func (o *OrderSyncDetail) GetUnitQty() float64 {
	return o.UnitQty
}
func (o *OrderSyncDetail) GetRawData() string {
	return ""
}

type OrderSyncAccount struct {
	FinanceCode string  `json:"financeCode"`
	Total       float64 `json:"total"`
}

// FinanceCode 账户编码 会检测编码在全渠道中是否存在
func (o *OrderSyncAccount) GetFinanceCode() string {
	return o.FinanceCode
}

// Total 收款金额 有账户编码的时候金额为必填
func (o *OrderSyncAccount) GetTotal() float64 {
	return o.Total
}
func (o *OrderSyncAccount) GetRawData() string {
	return ""
}

type OrderSyncResponse struct {
	Code    int64              `json:"code"`
	Message string             `json:"message"`
	Orders  []*OrderSyncResult `json:"orders"`
}

type OrderSyncResult struct {
	Tid      string `json:"tid"`
	BillCode string `json:"billcode"`
	Message  string `json:"message"`
}

func ToOrderSyncOrder(orderInfo *Order, nd PayNotifyDataInterface) *OrderSyncOrder {
	if orderInfo == nil {
		return nil
	}
	orderData := &OrderSyncOrder{
		Tid:              orderInfo.Tid,
		Weight:           orderInfo.Weight,
		Size:             orderInfo.Size,
		BuyerNick:        orderInfo.BuyerNick,
		BuyerMessage:     orderInfo.BuyerMessage,
		SellerMemo:       orderInfo.SellerMemo,
		Total:            orderInfo.Total,
		Privilege:        orderInfo.Privilege,
		PostFee:          orderInfo.PostFee,
		ReceiverName:     orderInfo.ReceiverName,
		ReceiverState:    stringUtils.FirstNonEmpty(orderInfo.ReceiverProvinceName, orderInfo.ReceiverProvince),
		ReceiverCity:     stringUtils.FirstNonEmpty(orderInfo.ReceiverCityName, orderInfo.ReceiverCity),
		ReceiverDistrict: stringUtils.FirstNonEmpty(orderInfo.ReceiverDistrictName, orderInfo.ReceiverDistrict),
		ReceiverAddress:  orderInfo.ReceiverAddress,
		ReceiverPhone:    orderInfo.ReceiverPhone,
		ReceiverMobile:   orderInfo.ReceiverMobile,
		Created:          time.FormatTime(orderInfo.CreateTime),
		Status:           orderInfo.Status,
		Type:             orderInfo.Type,
		InvoiceName:      orderInfo.InvoiceName,
		SellerFlag:       orderInfo.SellerFlag,
		PayTime:          time.FormatTime(orderInfo.PayTime),
		LogistBTypeCode:  orderInfo.LogistBTypeCode,
		LogistBillCode:   orderInfo.LogistBillCode,
		BTypeCode:        orderInfo.BTypeCode,
		Details:          toOrderSyncDetails(orderInfo.Details),
		Accounts:         toOrderSyncAccounts(orderInfo.Accounts),
	}
	if nd != nil {
		orderData.TransactionId = nd.GetTransactionID()
		orderData.NotifyRaw = nd.GetNotifyRaw()
		orderData.MchId = nd.GetMchID()
		orderData.Appid = nd.GetAppID()
		orderData.PayerOpenid = nd.GetPayerOpenID()
		orderData.PayChannel = nd.GetPayChannel()
	}
	return orderData
}

func toOrderSyncDetails(details []*OrderDetail) []*OrderSyncDetail {
	if len(details) == 0 {
		return []*OrderSyncDetail{}
	}
	result := make([]*OrderSyncDetail, 0, len(details))
	for _, detail := range details {
		if detail == nil {
			continue
		}
		result = append(result, &OrderSyncDetail{
			OID:            detail.OID,
			Barcode:        detail.Barcode,
			EShopGoodsID:   detail.EShopGoodsID,
			OuterIID:       detail.OuterIID,
			EShopGoodsName: detail.EShopGoodsName,
			EShopSKUId:     detail.EShopSkuID,
			EShopSKUName:   detail.EShopSkuName,
			NumIID:         detail.NumIID,
			SKUId:          detail.SkuID,
			Num:            detail.Num,
			Payment:        detail.Payment,
			PicPath:        detail.PicPath,
			Weight:         detail.Weight,
			Size:           detail.Size,
			UnitID:         detail.UnitID,
			UnitQty:        detail.UnitQty,
		})
	}
	return result
}

func toOrderSyncAccounts(accounts []*OrderAccount) []*OrderSyncAccount {
	if len(accounts) == 0 {
		return []*OrderSyncAccount{}
	}
	result := make([]*OrderSyncAccount, 0, len(accounts))
	for _, account := range accounts {
		if account == nil {
			continue
		}
		result = append(result, &OrderSyncAccount{
			FinanceCode: account.FinanceCode,
			Total:       account.Total,
		})
	}
	return result
}
