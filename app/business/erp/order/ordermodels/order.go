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
	CheckURL string `json:"check_url"`
}

// Order ERP订单主表
type Order struct {
	ID                   uint64          `json:"id,string" gorm:"column:id"`
	Tid                  string          `json:"tid" gorm:"column:tid"`
	Weight               float64         `json:"weight" gorm:"column:weight"`
	Size                 float64         `json:"size" gorm:"column:size"`
	BuyerNick            string          `json:"buyer_nick" gorm:"column:buyer_nick"`
	BuyerMessage         string          `json:"buyer_message" gorm:"column:buyer_message"`
	SellerMemo           string          `json:"seller_memo" gorm:"column:seller_memo"`
	Total                float64         `json:"total" gorm:"column:total"`
	Privilege            float64         `json:"privilege" gorm:"column:privilege"`
	PostFee              float64         `json:"post_fee" gorm:"column:post_fee"`
	ReceiverName         string          `json:"receiver_name" gorm:"column:receiver_name"`
	ReceiverProvince     string          `json:"receiver_province" gorm:"column:receiver_province"`
	ReceiverProvinceName string          `json:"receiver_province_name" gorm:"column:receiver_province_name"`
	ReceiverCity         string          `json:"receiver_city" gorm:"column:receiver_city"`
	ReceiverCityName     string          `json:"receiver_city_name" gorm:"column:receiver_city_name"`
	ReceiverDistrict     string          `json:"receiver_district" gorm:"column:receiver_district"`
	ReceiverDistrictName string          `json:"receiver_district_name" gorm:"column:receiver_district_name"`
	ReceiverStreet       string          `json:"receiver_street" gorm:"column:receiver_street"`
	ReceiverStreetName   string          `json:"receiver_street_name" gorm:"column:receiver_street_name"`
	ReceiverAddress      string          `json:"receiver_address" gorm:"column:receiver_address"`
	ReceiverPhone        string          `json:"receiver_phone" gorm:"column:receiver_phone"`
	ReceiverMobile       string          `json:"receiver_mobile" gorm:"column:receiver_mobile"`
	ReceiverZip          string          `json:"receiver_zip" gorm:"column:receiver_zip"`
	Status               string          `json:"status" gorm:"column:status"`
	Type                 string          `json:"order_type" gorm:"column:order_type"`
	InvoiceName          string          `json:"invoice_name" gorm:"column:invoice_name"`
	SellerFlag           string          `json:"seller_flag" gorm:"column:seller_flag"`
	PayTime              *time.Time      `json:"pay_time" gorm:"column:pay_time"`
	LogistBTypeCode      string          `json:"logist_b_type_code" gorm:"column:logist_b_type_code"`
	LogistBillCode       string          `json:"logist_bill_code" gorm:"column:logist_bill_code"`
	BTypeCode            string          `json:"b_type_code" gorm:"column:b_type_code"`
	BillCode             string          `json:"bill_code" gorm:"column:bill_code"`
	SyncMessage          string          `json:"sync_message" gorm:"column:sync_message"`
	SyncStatus           int32           `json:"sync_status" gorm:"column:sync_status"`
	SyncTime             *time.Time      `json:"sync_time" gorm:"column:sync_time"`
	DeptID               int64           `json:"dept_id" gorm:"column:dept_id"`
	Details              []*OrderDetail  `json:"details" gorm:"-"`
	Accounts             []*OrderAccount `json:"accounts" gorm:"-"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// OrderDetail ERP订单明细
type OrderDetail struct {
	ID             uint64  `json:"id,string" gorm:"column:id"`
	OrderID        uint64  `json:"order_id,string" gorm:"column:order_id"`
	Tid            string  `json:"tid" gorm:"column:tid"`
	OID            string  `json:"oid" gorm:"column:oid"`
	Barcode        string  `json:"barcode" gorm:"column:barcode"`
	EShopGoodsID   string  `json:"eshop_goods_id" gorm:"column:eshop_goods_id"`
	OuterIID       string  `json:"outer_iid" gorm:"column:outer_iid"`
	EShopGoodsName string  `json:"eshop_goods_name" gorm:"column:eshop_goods_name"`
	EShopSkuID     string  `json:"eshop_sku_id" gorm:"column:eshop_sku_id"`
	EShopSkuName   string  `json:"eshop_sku_name" gorm:"column:eshop_sku_name"`
	NumIID         int64   `json:"num_iid" gorm:"column:num_iid"`
	SkuID          int64   `json:"sku_id" gorm:"column:sku_id"`
	Num            float64 `json:"num" gorm:"column:num"`
	Payment        float64 `json:"payment" gorm:"column:payment"`
	PicPath        string  `json:"pic_path" gorm:"column:pic_path"`
	Weight         float64 `json:"weight" gorm:"column:weight"`
	Size           float64 `json:"size" gorm:"column:size"`
	UnitID         int64   `json:"unit_id" gorm:"column:unit_id"`
	UnitQty        float64 `json:"unit_qty" gorm:"column:unit_qty"`
	DeptID         int64   `json:"dept_id" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// OrderAccount ERP订单账户
type OrderAccount struct {
	ID          uint64  `json:"id,string" gorm:"column:id"`
	OrderID     uint64  `json:"order_id,string" gorm:"column:order_id"`
	Tid         string  `json:"tid" gorm:"column:tid"`
	FinanceCode string  `json:"finance_code" gorm:"column:finance_code"`
	Total       float64 `json:"total" gorm:"column:total"`
	DeptID      int64   `json:"dept_id" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// OrderSet ERP订单保存参数
type OrderSet struct {
	ID           uint64  `json:"id,string"`
	Tid          string  `json:"tid"`
	Weight       float64 `json:"weight"`
	Size         float64 `json:"size"`
	BuyerNick    string  `json:"buyer_nick"`
	BuyerMessage string  `json:"buyer_message"`
	SellerMemo   string  `json:"seller_memo"`
	Total        float64 `json:"total"`
	Privilege    float64 `json:"privilege"`
	PostFee      float64 `json:"post_fee"`
	ReceiverName string  `json:"receiver_name"`

	ReceiverProvince     string `json:"receiver_province"`
	ReceiverProvinceName string `json:"receiver_province_name"`

	ReceiverCity     string `json:"receiver_city"`
	ReceiverCityName string `json:"receiver_city_name"`

	ReceiverDistrict     string `json:"receiver_district"`
	ReceiverDistrictName string `json:"receiver_district_name"`

	ReceiverStreet     string `json:"receiver_street"`
	ReceiverStreetName string `json:"receiver_street_name"`

	ReceiverAddress string             `json:"receiver_address"`
	ReceiverPhone   string             `json:"receiver_phone"`
	ReceiverMobile  string             `json:"receiver_mobile"`
	ReceiverZip     string             `json:"receiver_zip"`
	Status          string             `json:"status"`
	OrderType       string             `json:"order_type"`
	InvoiceName     string             `json:"invoice_name"`
	SellerFlag      string             `json:"seller_flag"`
	PayTime         string             `json:"pay_time"`
	LogistBTypeCode string             `json:"logist_b_type_code"`
	LogistBillCode  string             `json:"logist_bill_code"`
	BTypeCode       string             `json:"b_type_code"`
	BillCode        string             `json:"bill_code"`
	SyncMessage     string             `json:"sync_message"`
	SyncStatus      int32              `json:"sync_status"`
	SyncTime        string             `json:"sync_time"`
	Details         []*OrderDetailSet  `json:"details"`
	Accounts        []*OrderAccountSet `json:"accounts"`
}

// OrderDetailSet ERP订单明细保存参数
type OrderDetailSet struct {
	OID            string  `json:"oid"`
	Barcode        string  `json:"barcode"`
	EShopGoodsID   string  `json:"eshop_goods_id"`
	OuterIID       string  `json:"outeriid"`
	EShopGoodsName string  `json:"eshop_goods_name"`
	EShopSkuID     string  `json:"eshop_sku_id"`
	EShopSkuName   string  `json:"eshop_sku_name"`
	NumIID         int64   `json:"numiid"`
	SkuID          int64   `json:"skuid"`
	Num            float64 `json:"num"`
	Payment        float64 `json:"payment"`
	PicPath        string  `json:"pic_path"`
	Weight         float64 `json:"weight"`
	Size           float64 `json:"size"`
	UnitID         int64   `json:"unit_id"`
	UnitQty        float64 `json:"unit_qty"`
}

// OrderAccountSet ERP订单账户保存参数
type OrderAccountSet struct {
	FinanceCode string  `json:"finance_code"`
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
