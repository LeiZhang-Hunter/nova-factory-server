package event

// Account 订单账户信息
type Account interface {
	// FinanceCode 账户编码 会检测编码在全渠道中是否存在
	GetFinanceCode() string
	// Total 收款金额 有账户编码的时候金额为必填
	GetTotal() float64
	GetRawData() string
}

// GoodsDetail 商品详情
type GoodsDetail interface {
	// GetOid 网店订单明细编号
	GetOid() string
	// GetBarcode 商品条码
	GetBarcode() string
	// GetEshopGoodsId 网店商品ID（对应商品列表接口中的productid），方案A必填
	GetEshopGoodsId() string
	// GetOuterIid 网店商家编码
	GetOuterIid() string
	// GetEshopGoodsName 网店商品名称
	GetEshopGoodsName() string
	// GetEshopSkuId 网店商品SKUID（对应商品列表接口中的skuid），方案A且商品有SKU则必填
	GetEshopSkuId() string
	// GetEshopSkuName 网店商品SKU名称，方案A且商品有SKU则必填
	GetEshopSkuName() string
	// GetNumIid 商品ID，方案B必填，方案A不要传值
	GetNumIid() int64
	// GetSkuId 	规格ID，方案B必填，方案A不要传值
	GetSkuId() int64
	// GetNum 基本单位数量（不能为0）
	GetNum() float64
	// GetPayment 商品总额
	GetPayment() float64
	// GetPicPath 商品图片路径
	GetPicPath() string
	// GetWeight 重量
	GetWeight() float64
	// GetSize 尺寸，体积
	GetSize() float64
	// GetUniTid 销售单位ID
	GetUniTid() int64
	// GetUnitQty 销售单位数量（无多单位可不填）
	GetUnitQty() float64
	GetRawData() string
}

type OrderData interface {
	// GetOrderNo 通用订单字段（所有业务模块都能提供）
	GetOrderNo() string // 订单号
	GetWeight() float64
	GetSize() float64
	GetBuyerNick() string    // 买家账号
	GetBuyerMessage() string // 卖家留言
	GetSellerMemo() string   // 卖家备注

	GetTotalAmount() float64 // 订单总金额

	GetPrivilege() float64 // 订单享受优惠的金额（订单总金额-实际支付金额）

	GetPostFee() float64 // 运费

	GetReceiverName() string // 收货人名称
	GetReceiverState() string
	GetReceiverCity() string
	GetReceiverDistrict() string
	GetReceiverAddress() string
	GetReceiverPhone() string
	GetReceiverMobile() string
	GetCreated() string
	GetType() string // 订单类型（Cod=货到付款, NoCod=非货到付款）

	// GetStatus  订单状态 NoPay = 未付款 Payed = 已付款（货到付款传已付款） Sended = 已发货 TradeSuccess = 交易成功 TradeClosed = 交易关闭 PartSend = 部分发货
	GetStatus() string // 状态

	GetInvoiceName() string // 发票抬头
	GetSellerFlag() string  // 	卖家旗帜（数值型）

	// GetPayTime 付款时间（时间格式：yyyy-MM-dd HH:mm:ss）
	GetPayTime() string
	// GetLogIstBTypeCode 物流公司编码（填入信息需和全渠道中的物流公司编码相同才可以匹配到，如 “ZTO”）
	GetLogIstBTypeCode() string

	// GetLogIstBillCode 物流单号
	GetLogIstBillCode() string

	// GetBTypeCode 往来单位编码
	GetBTypeCode() string

	// Details 订单商品信息
	GetDetails() []GoodsDetail

	// Accounts 订单账户信息
	GetAccounts() []Account

	Base
}

type OrderEvent interface {
	Event
	Base
	GetOrders() []OrderData
}
