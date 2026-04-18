package ordermodels

import (
	"nova-factory-server/app/baize"
	"time"
)

// CheckLoginStateReq 检查登录态请求
type CheckLoginStateReq struct {
	CheckURL string `form:"checkUrl"`
}

// CheckLoginStateResp 检查登录态响应
type CheckLoginStateResp struct {
	Online   bool   `json:"online"`
	Message  string `json:"message"`
	Type     string `json:"type"`
	CheckURL string `json:"checkUrl"`
}

// Order ERP订单主表
type Order struct {
	ID               uint64          `json:"id,string" gorm:"column:id"`
	Tid              string          `json:"tid" gorm:"column:tid"`
	Weight           float64         `json:"weight" gorm:"column:weight"`
	Size             float64         `json:"size" gorm:"column:size"`
	BuyerNick        string          `json:"buyernick" gorm:"column:buyernick"`
	BuyerMessage     string          `json:"buyermessage" gorm:"column:buyermessage"`
	SellerMemo       string          `json:"sellermemo" gorm:"column:sellermemo"`
	Total            float64         `json:"total" gorm:"column:total"`
	Privilege        float64         `json:"privilege" gorm:"column:privilege"`
	PostFee          float64         `json:"postfee" gorm:"column:postfee"`
	ReceiverName     string          `json:"receivername" gorm:"column:receivername"`
	ReceiverState    string          `json:"receiverstate" gorm:"column:receiverstate"`
	ReceiverCity     string          `json:"receivercity" gorm:"column:receivercity"`
	ReceiverDistrict string          `json:"receiverdistrict" gorm:"column:receiverdistrict"`
	ReceiverAddress  string          `json:"receiveraddress" gorm:"column:receiveraddress"`
	ReceiverPhone    string          `json:"receiverphone" gorm:"column:receiverphone"`
	ReceiverMobile   string          `json:"receivermobile" gorm:"column:receivermobile"`
	ReceiverZip      string          `json:"receiverzip" gorm:"column:receiverzip"`
	Status           string          `json:"status" gorm:"column:status"`
	Type             string          `json:"type" gorm:"column:type"`
	InvoiceName      string          `json:"invoicename" gorm:"column:invoicename"`
	SellerFlag       string          `json:"sellerflag" gorm:"column:sellerflag"`
	PayTime          *time.Time      `json:"paytime" gorm:"column:paytime"`
	LogistBTypeCode  string          `json:"logistbtypecode" gorm:"column:logistbtypecode"`
	LogistBillCode   string          `json:"logistbillcode" gorm:"column:logistbillcode"`
	BTypeCode        string          `json:"btypecode" gorm:"column:btypecode"`
	BillCode         string          `json:"billcode" gorm:"column:billcode"`
	SyncMessage      string          `json:"syncMessage" gorm:"column:sync_message"`
	SyncStatus       int32           `json:"syncStatus" gorm:"column:sync_status"`
	SyncTime         *time.Time      `json:"syncTime" gorm:"column:sync_time"`
	DeptID           int64           `json:"deptId" gorm:"column:dept_id"`
	Details          []*OrderDetail  `json:"details" gorm:"-"`
	Accounts         []*OrderAccount `json:"accounts" gorm:"-"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// OrderDetail ERP订单明细
type OrderDetail struct {
	ID             uint64  `json:"id,string" gorm:"column:id"`
	OrderID        uint64  `json:"orderId,string" gorm:"column:order_id"`
	Tid            string  `json:"tid" gorm:"column:tid"`
	OID            string  `json:"oid" gorm:"column:oid"`
	Barcode        string  `json:"barcode" gorm:"column:barcode"`
	EShopGoodsID   string  `json:"eshopgoodsid" gorm:"column:eshopgoodsid"`
	OuterIID       string  `json:"outeriid" gorm:"column:outeriid"`
	EShopGoodsName string  `json:"eshopgoodsname" gorm:"column:eshopgoodsname"`
	EShopSkuID     string  `json:"eshopskuid" gorm:"column:eshopskuid"`
	EShopSkuName   string  `json:"eshopskuname" gorm:"column:eshopskuname"`
	NumIID         int64   `json:"numiid" gorm:"column:numiid"`
	SkuID          int64   `json:"skuid" gorm:"column:skuid"`
	Num            float64 `json:"num" gorm:"column:num"`
	Payment        float64 `json:"payment" gorm:"column:payment"`
	PicPath        string  `json:"picpath" gorm:"column:picpath"`
	Weight         float64 `json:"weight" gorm:"column:weight"`
	Size           float64 `json:"size" gorm:"column:size"`
	UnitID         int64   `json:"unitid" gorm:"column:unitid"`
	UnitQty        float64 `json:"unitqty" gorm:"column:unitqty"`
	DeptID         int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// OrderAccount ERP订单账户
type OrderAccount struct {
	ID          uint64  `json:"id,string" gorm:"column:id"`
	OrderID     uint64  `json:"orderId,string" gorm:"column:order_id"`
	Tid         string  `json:"tid" gorm:"column:tid"`
	FinanceCode string  `json:"financeCode" gorm:"column:finance_code"`
	Total       float64 `json:"total" gorm:"column:total"`
	DeptID      int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// OrderSet ERP订单保存参数
type OrderSet struct {
	ID           uint64  `json:"id,string"`
	Tid          string  `json:"tid"`
	Weight       float64 `json:"weight"`
	Size         float64 `json:"size"`
	BuyerNick    string  `json:"buyernick"`
	BuyerMessage string  `json:"buyermessage"`
	SellerMemo   string  `json:"sellermemo"`
	Total        float64 `json:"total"`
	Privilege    float64 `json:"privilege"`
	PostFee      float64 `json:"postfee"`
	ReceiverName string  `json:"receivername"`

	ReceiverState     string `json:"receiverstate"`
	ReceiverStateName string `json:"receiverstate_name"`

	ReceiverStreet     string `json:"receiverstreet"`
	ReceiverStreetName string `json:"receiverstreet_name"`

	ReceiverCity     string `json:"receivercity"`
	ReceiverCityName string `json:"receivercity_name"`

	ReceiverDistrict     string `json:"receiverdistrict"`
	ReceiverDistrictName string `json:"receiverdistrict_name"`

	ReceiverAddress string             `json:"receiveraddress"`
	ReceiverPhone   string             `json:"receiverphone"`
	ReceiverMobile  string             `json:"receivermobile"`
	ReceiverZip     string             `json:"receiverzip"`
	Status          string             `json:"status"`
	Type            string             `json:"type"`
	InvoiceName     string             `json:"invoicename"`
	SellerFlag      string             `json:"sellerflag"`
	PayTime         string             `json:"paytime"`
	LogistBTypeCode string             `json:"logistbtypecode"`
	LogistBillCode  string             `json:"logistbillcode"`
	BTypeCode       string             `json:"btypecode"`
	BillCode        string             `json:"billcode"`
	SyncMessage     string             `json:"syncMessage"`
	SyncStatus      int32              `json:"syncStatus"`
	SyncTime        string             `json:"syncTime"`
	Details         []*OrderDetailSet  `json:"details"`
	Accounts        []*OrderAccountSet `json:"accounts"`
}

// OrderDetailSet ERP订单明细保存参数
type OrderDetailSet struct {
	OID            string  `json:"oid"`
	Barcode        string  `json:"barcode"`
	EShopGoodsID   string  `json:"eshopgoodsid"`
	OuterIID       string  `json:"outeriid"`
	EShopGoodsName string  `json:"eshopgoodsname"`
	EShopSkuID     string  `json:"eshopskuid"`
	EShopSkuName   string  `json:"eshopskuname"`
	NumIID         int64   `json:"numiid"`
	SkuID          int64   `json:"skuid"`
	Num            float64 `json:"num"`
	Payment        float64 `json:"payment"`
	PicPath        string  `json:"picpath"`
	Weight         float64 `json:"weight"`
	Size           float64 `json:"size"`
	UnitID         int64   `json:"unitid"`
	UnitQty        float64 `json:"unitqty"`
}

// OrderAccountSet ERP订单账户保存参数
type OrderAccountSet struct {
	FinanceCode string  `json:"financeCode"`
	Total       float64 `json:"total"`
}

// OrderQuery ERP订单查询参数
type OrderQuery struct {
	Tid          string `form:"tid"`
	Status       string `form:"status"`
	SyncStatus   *int32 `form:"syncStatus"`
	ReceiverName string `form:"receiverName"`
	Page         int64  `form:"page"`
	Size         int64  `form:"size"`
}

// OrderListData ERP订单列表结果
type OrderListData struct {
	Rows  []*Order `json:"rows"`
	Total int64    `json:"total"`
}
