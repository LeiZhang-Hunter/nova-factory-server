package grasp

// OrderSyncRequest 订单同步请求体
type OrderSyncRequest struct {
	Orders []*OrderSyncOrder `json:"orders"`
}

// OrderSyncOrder 单笔订单数据
type OrderSyncOrder struct {
	Tid              string              `json:"tid"`
	Weight           *float64            `json:"weight"`
	Size             *float64            `json:"size"`
	BuyerNick        *string             `json:"buyernick"`
	BuyerMessage     *string             `json:"buyermessage"`
	SellerMemo       *string             `json:"sellermemo"`
	Total            *float64            `json:"total"`
	Privilege        float64             `json:"privilege"`
	PostFee          float64             `json:"postfee"`
	ReceiverName     string              `json:"receivername"`
	ReceiverState    string              `json:"receiverstate"`
	ReceiverCity     string              `json:"receivercity"`
	ReceiverDistrict string              `json:"receiverdistrict"`
	ReceiverAddress  string              `json:"receiveraddress"`
	ReceiverPhone    *string             `json:"receiverphone"`
	ReceiverMobile   string              `json:"receivermobile"`
	Created          string              `json:"created"`
	Status           string              `json:"status"`
	Type             string              `json:"type"`
	InvoiceName      *string             `json:"invoicename"`
	SellerFlag       *string             `json:"sellerflag"`
	PayTime          *string             `json:"paytime"`
	LogistBTypeCode  *string             `json:"logistbtypecode"`
	LogistBillCode   *string             `json:"logistbillcode"`
	BTypeCode        *string             `json:"btypecode"`
	Details          []*OrderSyncDetail  `json:"details"`
	Accounts         []*OrderSyncAccount `json:"accounts,omitempty"`
}

// OrderSyncDetail 订单明细行
type OrderSyncDetail struct {
	OID            string   `json:"oid"`
	Barcode        *string  `json:"barcode"`
	EShopGoodsID   *string  `json:"eshopgoodsid"`
	OuterIID       *string  `json:"outeriid"`
	EShopGoodsName string   `json:"eshopgoodsname"`
	EShopSKUId     *string  `json:"eshopskuid"`
	EShopSKUName   *string  `json:"eshopskuname"`
	NumIID         *int64   `json:"numiid"`
	SKUId          *int64   `json:"skuid"`
	Num            float64  `json:"num"`
	Payment        float64  `json:"payment"`
	PicPath        *string  `json:"picpath"`
	Weight         *float64 `json:"weight"`
	Size           *float64 `json:"size"`
	UnitID         *int64   `json:"unitid"`
	UnitQty        *float64 `json:"unitqty"`
}

// OrderSyncAccount 订单收款账户信息
type OrderSyncAccount struct {
	FinanceCode string  `json:"financeCode"`
	Total       float64 `json:"total"`
}

// OrderSyncResponse 订单同步响应体
type OrderSyncResponse struct {
	Code    int64              `json:"code"`
	Message string             `json:"message"`
	Orders  []*OrderSyncResult `json:"orders"`
}

// OrderSyncResult 单笔订单同步结果
type OrderSyncResult struct {
	Tid      string `json:"tid"`
	BillCode string `json:"billcode"`
	Message  string `json:"message"`
}
