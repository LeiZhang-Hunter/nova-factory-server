package models

import (
	"fmt"
	"nova-factory-server/app/baize"
	searchutil "nova-factory-server/app/utils/vectorsearch"
	"strings"
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
	TransactionID        string          `json:"transaction_id" gorm:"column:transaction_id"`
	NotifyRaw            string          `json:"-" gorm:"column:notify_raw"`
	MchID                string          `json:"mch_id" gorm:"column:mch_id"`
	AppID                string          `json:"appid" gorm:"column:appid"`
	PayerOpenid          string          `json:"payer_openid" gorm:"column:payer_openid"`
	LogistBTypeCode      string          `json:"logist_b_type_code" gorm:"column:logist_b_type_code"`
	LogistBillCode       string          `json:"logist_bill_code" gorm:"column:logist_bill_code"`
	BTypeCode            string          `json:"b_type_code" gorm:"column:b_type_code"`
	DetailsJSON          string          `json:"-" gorm:"column:details_json"`
	AccountsJSON         string          `json:"-" gorm:"column:accounts_json"`
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
	ID                uint64  `json:"id,string" gorm:"column:id"`
	OrderID           uint64  `json:"order_id,string" gorm:"column:order_id"`
	Tid               string  `json:"tid" gorm:"column:tid"`
	OID               string  `json:"oid" gorm:"column:oid"`
	Barcode           string  `json:"barcode" gorm:"column:barcode"`
	OldBarcode        string  `json:"old_barcode" gorm:"-"`
	EShopGoodsID      string  `json:"eshop_goods_id" gorm:"column:eshop_goods_id"`
	OuterIID          string  `json:"outer_iid" gorm:"column:outer_iid"`
	EShopGoodsName    string  `json:"eshop_goods_name" gorm:"column:eshop_goods_name"`
	OldEShopGoodsName string  `json:"old_eshop_goods_name" gorm:"-"`
	EShopSkuID        string  `json:"eshop_sku_id" gorm:"column:eshop_sku_id"`
	OldEShopSkuID     string  `json:"old_eshop_sku_id" gorm:"-"`
	EShopSkuName      string  `json:"eshop_sku_name" gorm:"column:eshop_sku_name"`
	OldEShopSkuName   string  `json:"old_eshop_sku_name" gorm:"-"`
	NumIID            int64   `json:"num_iid" gorm:"column:num_iid"`
	SkuID             int64   `json:"sku_id" gorm:"column:sku_id"`
	Num               float64 `json:"num" gorm:"column:num"`
	Payment           float64 `json:"payment" gorm:"column:payment"`
	OldPayment        float64 `json:"old_payment" gorm:"-"`
	PicPath           string  `json:"pic_path" gorm:"column:pic_path"`
	Weight            float64 `json:"weight" gorm:"column:weight"`
	Size              float64 `json:"size" gorm:"column:size"`
	UnitID            int64   `json:"unit_id" gorm:"column:unit_id"`
	UnitQty           float64 `json:"unit_qty" gorm:"column:unit_qty"`
	DeptID            int64   `json:"dept_id" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// VectorSearchLabeledValues 返回订单明细可参与向量检索的结构化文本字段。
func (d *OrderDetail) VectorSearchLabeledValues() []searchutil.LabeledValue {
	if d == nil {
		return nil
	}
	values := make([]searchutil.LabeledValue, 0, 12)

	productCode := strings.TrimSpace(d.OuterIID)
	if productCode == "" {
		productCode = strings.TrimSpace(d.EShopGoodsID)
	}

	unitParts := make([]string, 0, 2)
	if d.UnitID > 0 {
		unitParts = append(unitParts, fmt.Sprintf("单位ID:%d", d.UnitID))
	}
	if d.UnitQty > 0 {
		unitParts = append(unitParts, fmt.Sprintf("单位数量:%.3f", d.UnitQty))
	}

	remarkParts := make([]string, 0, 6)
	if goodsID := strings.TrimSpace(d.EShopGoodsID); goodsID != "" {
		remarkParts = append(remarkParts, "商品ID:"+goodsID)
	}
	if skuID := strings.TrimSpace(d.EShopSkuID); skuID != "" {
		remarkParts = append(remarkParts, "SKU ID:"+skuID)
	}
	if d.NumIID > 0 {
		remarkParts = append(remarkParts, fmt.Sprintf("商品数字ID:%d", d.NumIID))
	}
	if d.SkuID > 0 {
		remarkParts = append(remarkParts, fmt.Sprintf("SKU数字ID:%d", d.SkuID))
	}
	if d.Size > 0 {
		remarkParts = append(remarkParts, fmt.Sprintf("体积:%.3f", d.Size))
	}

	values = append(values,
		searchutil.LabeledValue{Label: "产品名称", Value: d.EShopGoodsName},
		searchutil.LabeledValue{Label: "产品编码", Value: productCode},
		searchutil.LabeledValue{Label: "产品分类", Value: ""},
		searchutil.LabeledValue{Label: "单位", Value: strings.Join(unitParts, " ")},
		searchutil.LabeledValue{Label: "条码", Value: d.Barcode},
		searchutil.LabeledValue{Label: "规格", Value: d.EShopSkuName},
		searchutil.LabeledValue{Label: "备注", Value: strings.Join(remarkParts, " ")},
	)
	if d.Weight > 0 {
		values = append(values, searchutil.LabeledValue{Label: "重量", Value: fmt.Sprintf("%.3fkg", d.Weight)})
	}
	if d.Payment > 0 {
		values = append(values, searchutil.LabeledValue{Label: "销售价", Value: fmt.Sprintf("%.2f", d.Payment)})
	}
	return values
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
	OID            string  `json:"oid" jsonschema:"description=订单明细编号"`
	Barcode        string  `json:"barcode" jsonschema:"description=商品条码"`
	EShopGoodsID   string  `json:"eshop_goods_id" jsonschema:"description=电商平台商品ID"`
	OuterIID       string  `json:"outeriid" jsonschema:"description=外部商品编号"`
	EShopGoodsName string  `json:"eshop_goods_name" jsonschema:"description=商品名称"`
	EShopSkuID     string  `json:"eshop_sku_id" jsonschema:"description=电商平台SKU ID"`
	EShopSkuName   string  `json:"eshop_sku_name" jsonschema:"description=SKU名称"`
	NumIID         int64   `json:"numiid" jsonschema:"description=商品数字ID"`
	SkuID          int64   `json:"skuid" jsonschema:"description=SKU数字ID"`
	Num            float64 `json:"num" jsonschema:"description=购买数量"`
	Payment        float64 `json:"payment" jsonschema:"description=明细实付金额"`
	PicPath        string  `json:"pic_path" jsonschema:"description=商品图片地址"`
	Weight         float64 `json:"weight" jsonschema:"description=明细重量"`
	Size           float64 `json:"size" jsonschema:"description=明细体积"`
	UnitID         int64   `json:"unit_id" jsonschema:"description=单位ID"`
	UnitQty        float64 `json:"unit_qty" jsonschema:"description=单位数量"`
}

// OrderAccountSet ERP订单账户保存参数
type OrderAccountSet struct {
	FinanceCode string  `json:"finance_code" jsonschema:"description=财务科目编码"`
	Total       float64 `json:"total" jsonschema:"description=账户金额"`
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

// OrderAudit ERP订单审核主表
type OrderAudit struct {
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
	DetailsJSON          string          `json:"-" gorm:"column:details_json"`
	AccountsJSON         string          `json:"-" gorm:"column:accounts_json"`
	SourceJSON           string          `json:"source_json" gorm:"column:source_json"`
	AuditStatus          int32           `json:"audit_status" gorm:"column:audit_status"`
	AuditRemark          string          `json:"audit_remark" gorm:"column:audit_remark"`
	AuditBy              int64           `json:"audit_by" gorm:"column:audit_by"`
	AuditTime            *time.Time      `json:"audit_time" gorm:"column:audit_time"`
	TransferStatus       int32           `json:"transfer_status" gorm:"column:transfer_status"`
	TransferMessage      string          `json:"transfer_message" gorm:"column:transfer_message"`
	TransferTime         *time.Time      `json:"transfer_time" gorm:"column:transfer_time"`
	ERPOrderID           uint64          `json:"erp_order_id,string" gorm:"column:erp_order_id"`
	DeptID               int64           `json:"dept_id" gorm:"column:dept_id"`
	Details              []*OrderDetail  `json:"details" gorm:"-"`
	Accounts             []*OrderAccount `json:"accounts" gorm:"-"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// OrderAuditSet ERP订单审核保存参数
type OrderAuditSet struct {
	ID                   uint64             `json:"id,string" jsonschema:"description=审核订单ID，更新时传"`
	Tid                  string             `json:"tid" jsonschema:"description=网店订单编号"`
	Weight               float64            `json:"weight" jsonschema:"description=订单重量"`
	Size                 float64            `json:"size" jsonschema:"description=订单尺寸或体积"`
	BuyerNick            string             `json:"buyer_nick" jsonschema:"description=买家账号"`
	BuyerMessage         string             `json:"buyer_message" jsonschema:"description=买家留言"`
	SellerMemo           string             `json:"seller_memo" jsonschema:"description=卖家备注"`
	Total                float64            `json:"total" jsonschema:"description=订单总金额，含运费"`
	Privilege            float64            `json:"privilege" jsonschema:"description=订单优惠金额"`
	PostFee              float64            `json:"post_fee" jsonschema:"description=运费"`
	ReceiverName         string             `json:"receiver_name" jsonschema:"description=收货人名称"`
	ReceiverProvince     string             `json:"receiver_province" jsonschema:"description=收货省编码"`
	ReceiverProvinceName string             `json:"receiver_province_name" jsonschema:"description=收货省名称"`
	ReceiverCity         string             `json:"receiver_city" jsonschema:"description=收货市编码"`
	ReceiverCityName     string             `json:"receiver_city_name" jsonschema:"description=收货市名称"`
	ReceiverDistrict     string             `json:"receiver_district" jsonschema:"description=收货区编码"`
	ReceiverDistrictName string             `json:"receiver_district_name" jsonschema:"description=收货区名称"`
	ReceiverStreet       string             `json:"receiver_street" jsonschema:"description=收货街道编码"`
	ReceiverStreetName   string             `json:"receiver_street_name" jsonschema:"description=收货街道名称"`
	ReceiverAddress      string             `json:"receiver_address" jsonschema:"description=收货详细地址"`
	ReceiverPhone        string             `json:"receiver_phone" jsonschema:"description=收货电话"`
	ReceiverMobile       string             `json:"receiver_mobile" jsonschema:"description=收货手机号"`
	ReceiverZip          string             `json:"receiver_zip" jsonschema:"description=收货邮编"`
	Status               string             `json:"status" jsonschema:"description=来源订单状态"`
	OrderType            string             `json:"order_type" jsonschema:"description=订单类型，如 Cod 或 NoCod"`
	InvoiceName          string             `json:"invoice_name" jsonschema:"description=发票抬头"`
	SellerFlag           string             `json:"seller_flag" jsonschema:"description=卖家旗帜"`
	PayTime              string             `json:"pay_time" jsonschema:"description=付款时间，格式为 2006-01-02 15:04:05"`
	LogistBTypeCode      string             `json:"logist_b_type_code" jsonschema:"description=物流公司编码"`
	LogistBillCode       string             `json:"logist_bill_code" jsonschema:"description=物流单号"`
	BTypeCode            string             `json:"b_type_code" jsonschema:"description=往来单位编码"`
	SourceJSON           string             `json:"source_json" jsonschema:"description=原始订单报文JSON，留空时自动生成"`
	Details              []*OrderDetailSet  `json:"details" jsonschema:"description=订单明细列表"`
	Accounts             []*OrderAccountSet `json:"accounts" jsonschema:"description=订单账户列表"`
}

// OrderAuditQuery ERP订单审核查询参数
type OrderAuditQuery struct {
	Tid            string `form:"tid"`
	AuditStatus    *int32 `form:"auditStatus"`
	TransferStatus *int32 `form:"transferStatus"`
	ReceiverName   string `form:"receiverName"`
	Page           int64  `form:"page"`
	Size           int64  `form:"size"`
}

// OrderAuditListData ERP订单审核列表结果
type OrderAuditListData struct {
	Rows  []*OrderAudit `json:"rows"`
	Total int64         `json:"total"`
}

// OrderAuditAction ERP订单审核动作参数
type OrderAuditAction struct {
	ID          uint64 `json:"id,string"`
	AuditRemark string `json:"audit_remark"`
}

// OrderAuditImportReq ERP订单审核批量导入参数
type OrderAuditImportReq struct {
	Records []*OrderAuditSet `json:"records" jsonschema:"description=待导入的订单审核记录列表"`
}

// OrderAuditImportItem ERP订单审核批量导入结果项
type OrderAuditImportItem struct {
	Index   int    `json:"index"`
	Tid     string `json:"tid"`
	ID      uint64 `json:"id,string"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// OrderAuditImportResult ERP订单审核批量导入结果
type OrderAuditImportResult struct {
	Total        int                     `json:"total"`
	SuccessCount int                     `json:"successCount"`
	FailureCount int                     `json:"failureCount"`
	Items        []*OrderAuditImportItem `json:"items"`
}
