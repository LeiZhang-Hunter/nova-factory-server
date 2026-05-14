package order

const (
	ERPStatusNoPay        string = "NoPay"        // 未付款
	ERPStatusPayed        string = "Payed"        // 已付款（货到付款）
	ERPStatusSended       string = "Sended"       // 已发货
	ERPStatusPartSend     string = "PartSend"     // 部分发货
	ERPStatusTradeSuccess string = "TradeSuccess" // 交易成功
	ERPStatusTradeClosed  string = "TradeClosed"  // 交易关闭
	ERPStatusAftersale    string = "Aftersale"    // 售后/退款
)
