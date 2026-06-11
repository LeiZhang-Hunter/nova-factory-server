package event

import (
	"time"
)

// Account 订单账户信息
type Account interface {
	// FinanceCode 账户编码 会检测编码在全渠道中是否存在
	FinanceCode() string
	// Total 收款金额 有账户编码的时候金额为必填
	Total() float64
	RawData() string
}

// GoodsDetail 商品详情
type GoodsDetail interface {
	// Oid 网店订单明细编号
	Oid() string
	// Barcode 商品条码
	Barcode() string
	// EshopGoodsId 网店商品ID（对应商品列表接口中的productid），方案A必填
	EshopGoodsId() string
	// OuterIid 网店商家编码
	OuterIid() string
	// EshopGoodsName 网店商品名称
	EshopGoodsName() string
	// EshopSkuId 网店商品SKUID（对应商品列表接口中的skuid），方案A且商品有SKU则必填
	EshopSkuId() string
	// EshopSkuName 网店商品SKU名称，方案A且商品有SKU则必填
	EshopSkuName() string
	// NumIid 商品ID，方案B必填，方案A不要传值
	NumIid() int64
	// SkuId 	规格ID，方案B必填，方案A不要传值
	SkuId() int64
	// Num 基本单位数量（不能为0）
	Num() float64
	// Payment 商品总额
	Payment() float64
	// PicPath 商品图片路径
	PicPath() string
	// Weight 重量
	Weight() float64
	// Size 尺寸，体积
	Size() float64
	// UniTid 销售单位ID
	UniTid() int64
	// UnitQty 销售单位数量（无多单位可不填）
	UnitQty() float64
	RawData() string
}

type OrderData interface {
	// OrderNo 通用订单字段（所有业务模块都能提供）
	OrderNo() string // 订单号
	Weight() float64
	Size() float64
	BuyerNick() string    // 买家账号
	BuyerMessage() string // 卖家留言
	SellerMemo() string   // 卖家备注

	TotalAmount() float64 // 订单总金额

	Privilege() float64 // 订单享受优惠的金额（订单总金额-实际支付金额）

	PostFee() float64 // 运费

	ReceiverName() string // 收货人名称
	ReceiverState() string
	ReceiverCity() string
	ReceiverDistrict() string
	ReceiverAddress() string
	ReceiverPhone() string
	ReceiverMobile() string
	Created() time.Time
	Type() string // 订单类型（Cod=货到付款, NoCod=非货到付款）

	// Status  订单状态 NoPay = 未付款 Payed = 已付款（货到付款传已付款） Sended = 已发货 TradeSuccess = 交易成功 TradeClosed = 交易关闭 PartSend = 部分发货
	Status() string // 状态

	InvoiceName() string // 发票抬头
	SellerFlag() string  // 	卖家旗帜（数值型）

	// PayTime 付款时间（时间格式：yyyy-MM-dd HH:mm:ss）
	PayTime() string
	// LogIstBTypeCode 物流公司编码（填入信息需和全渠道中的物流公司编码相同才可以匹配到，如 “ZTO”）
	LogIstBTypeCode() string

	// LogIstBillCode 物流单号
	LogIstBillCode() string

	// BTypeCode 往来单位编码
	BTypeCode() string

	// Details 订单商品信息
	Details() []GoodsDetail

	// Accounts 订单账户信息
	Accounts() []Account

	Base
}

type OrderEvent interface {
	Event
	Base
	Orders() []OrderData
}
