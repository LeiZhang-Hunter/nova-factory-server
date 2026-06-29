// 定义订单事件相关的数据类型与接口。
// 包含订单账户（Account）、订单商品明细（GoodsDetail）、
// 订单主体数据（OrderData）及订单事件（OrderEvent），
// 用于将网店订单同步至第三方 ERP/全渠道系统（如管家婆）。
package event

// Account 订单账户信息，记录订单收款的账户及其金额。
// 一般用于货到付款或多账户收款场景，需在全渠道系统中预先配置对应编码。
type Account interface {
	// GetFinanceCode 账户编码，系统会检测该编码在全渠道中是否存在，不存在将报错
	GetFinanceCode() string
	// GetTotal 收款金额，有账户编码时金额为必填字段
	GetTotal() float64
	// GetRawData 返回账户原始数据字符串，用于调试和日志
	GetRawData() string
}

// GoodsDetail 订单商品明细接口，描述订单中每件商品的信息。
// 支持方案A（按网店商品ID+SKUID匹配）和方案B（按ERP商品ID+规格ID匹配）两种对接模式。
type GoodsDetail interface {
	// GetOid 网店订单明细编号，用于唯一标识订单中一条商品记录
	GetOid() string
	// GetBarcode 商品条码
	GetBarcode() string
	// GetEshopGoodsId 网店商品ID（对应商品列表接口中的 productid），方案A必填
	GetEshopGoodsId() string
	// GetOuterIid 网店商家编码，用于商家自定义的商品编号
	GetOuterIid() string
	// GetEshopGoodsName 网店商品名称
	GetEshopGoodsName() string
	// GetEshopSkuId 网店商品SKU ID（对应商品列表接口中的 skuid），方案A且商品有SKU则必填
	GetEshopSkuId() string
	// GetEshopSkuName 网店商品SKU名称，方案A且商品有SKU则必填
	GetEshopSkuName() string
	// GetNumIid ERP商品ID，方案B必填，方案A不要传值
	GetNumIid() int64
	// GetSkuId ERP规格ID，方案B必填，方案A不要传值
	GetSkuId() int64
	// GetNum 基本单位数量，不能为0
	GetNum() float64
	// GetPayment 商品总额（单价*数量）
	GetPayment() float64
	// GetPicPath 商品图片路径（URL）
	GetPicPath() string
	// GetWeight 商品重量
	GetWeight() float64
	// GetSize 商品尺寸/体积
	GetSize() float64
	// GetUniTid 销售单位ID，用于多单位场景
	GetUniTid() int64
	// GetUnitQty 销售单位数量，无多单位时可不填
	GetUnitQty() float64
	// GetRawData 返回商品明细原始数据字符串
	GetRawData() string
}

// OrderData 订单数据接口，封装一条完整订单的所有字段。
// 各业务模块（如不同网店平台）只需实现此接口即可接入订单同步流程。
type OrderData interface {
	// GetOrderNo 订单号，所有业务模块都必须提供
	GetOrderNo() string
	// GetUserId 买家ID
	GetUserId() uint64
	// GetWeight 订单总重量
	GetWeight() float64
	// GetSize 订单总尺寸/体积
	GetSize() float64
	// GetBuyerNick 买家账号/昵称
	GetBuyerNick() string
	// GetBuyerMessage 买家留言
	GetBuyerMessage() string
	// GetSellerMemo 卖家备注
	GetSellerMemo() string
	// GetTotalAmount 订单总金额（优惠前）
	GetTotalAmount() float64
	// GetPrivilege 订单享受的优惠金额（订单总金额 - 实际支付金额）
	GetPrivilege() float64
	// GetPostFee 运费
	GetPostFee() float64
	// GetReceiverName 收货人名称
	GetReceiverName() string
	// GetReceiverState 收货省份
	GetReceiverState() string
	GetReceiverStateName() string
	// GetReceiverCity 收货城市
	GetReceiverCity() string
	GetReceiverCityName() string
	// GetReceiverDistrict 收货区/县
	GetReceiverDistrict() string
	GetReceiverDistrictName() string
	// GetReceiverAddress 收货详细地址
	GetReceiverAddress() string
	// GetReceiverPhone 收货人座机号码
	GetReceiverPhone() string
	// GetReceiverMobile 收货人手机号码
	GetReceiverMobile() string
	// GetCreated 订单创建时间
	GetCreated() string
	// GetType 订单类型，Cod=货到付款, NoCod=非货到付款
	GetType() string
	// GetStatus 订单状态，NoPay=未付款, Payed=已付款（货到付款传已付款）,
	// Sended=已发货, TradeSuccess=交易成功, TradeClosed=交易关闭, PartSend=部分发货
	GetStatus() string
	// GetInvoiceName 发票抬头
	GetInvoiceName() string
	// GetSellerFlag 卖家旗帜（数值型，用于标记订单重要程度等）
	GetSellerFlag() string
	// GetPayTime 付款时间，格式 yyyy-MM-dd HH:mm:ss
	GetPayTime() string
	// GetLogIstBTypeCode 物流公司编码，需与全渠道中的物流公司编码一致，如 "ZTO"（中通）
	GetLogIstBTypeCode() string
	// GetLogIstBillCode 物流单号（运单号）
	GetLogIstBillCode() string
	// GetBTypeCode 往来单位编码，用于关联ERP中的客户档案
	GetBTypeCode() string
	// GetDetails 订单商品明细列表
	GetDetails() []GoodsDetail
	// GetAccounts 订单账户信息列表，记录收款方式与金额
	GetAccounts() []Account
	// GetTransactionId 支付订单号
	GetTransactionId() string
	// GetNotifyRaw 支付回调原始数据
	GetNotifyRaw() string
	// GetMchId 支付商户号
	GetMchId() string
	// GetAppid 支付应用ID
	GetAppid() string
	// GetPayerOpenid 支付用户标识
	GetPayerOpenid() string
	// GetPayChannel 支付渠道
	GetPayChannel() int
	Base
}

// OrderEvent 订单事件接口，表示一次订单创建或状态变更事件。
// 聚合了 Event（事件元信息）和 Base（基础能力），并提供订单数据列表。
type OrderEvent interface {
	Event
	Base
	// GetOrders 返回本次事件涉及的订单数据列表
	GetOrders() []OrderData
}
