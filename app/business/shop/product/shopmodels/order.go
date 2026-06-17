package shopmodels

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/utils/observer/integration/event"
	"time"
)

// Order ERP订单主表
type Order struct {
	ID                   uint64          `json:"id,string" gorm:"column:id"`
	UserId               uint64          `json:"user_id,string" gorm:"column:user_id"`
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

// ToOrder 转化订单
func ToOrder(event event.OrderEvent) []*Order {
	if event == nil {
		return nil
	}

	eventOrders := event.GetOrders()
	if len(eventOrders) == 0 {
		return []*Order{}
	}

	orders := make([]*Order, 0, len(eventOrders))
	for _, eventOrder := range eventOrders {
		if eventOrder == nil {
			continue
		}

		order := &Order{
			Tid:                  eventOrder.GetOrderNo(),
			UserId:               eventOrder.GetUserId(),
			Weight:               eventOrder.GetWeight(),
			Size:                 eventOrder.GetSize(),
			BuyerNick:            eventOrder.GetBuyerNick(),
			BuyerMessage:         eventOrder.GetBuyerMessage(),
			SellerMemo:           eventOrder.GetSellerMemo(),
			Total:                eventOrder.GetTotalAmount(),
			Privilege:            eventOrder.GetPrivilege(),
			PostFee:              eventOrder.GetPostFee(),
			ReceiverName:         eventOrder.GetReceiverName(),
			ReceiverProvince:     eventOrder.GetReceiverState(),
			ReceiverProvinceName: eventOrder.GetReceiverState(),
			ReceiverCity:         eventOrder.GetReceiverCity(),
			ReceiverCityName:     eventOrder.GetReceiverCity(),
			ReceiverDistrict:     eventOrder.GetReceiverDistrict(),
			ReceiverDistrictName: eventOrder.GetReceiverDistrict(),
			ReceiverAddress:      eventOrder.GetReceiverAddress(),
			ReceiverPhone:        eventOrder.GetReceiverPhone(),
			ReceiverMobile:       eventOrder.GetReceiverMobile(),
			Status:               eventOrder.GetStatus(),
			Type:                 eventOrder.GetType(),
			InvoiceName:          eventOrder.GetInvoiceName(),
			SellerFlag:           eventOrder.GetSellerFlag(),
			LogistBTypeCode:      eventOrder.GetLogIstBTypeCode(),
			LogistBillCode:       eventOrder.GetLogIstBillCode(),
			BTypeCode:            eventOrder.GetBTypeCode(),
			Details:              toOrderDetails(eventOrder.GetOrderNo(), eventOrder.GetDetails()),
			Accounts:             toOrderAccounts(eventOrder.GetOrderNo(), eventOrder.GetAccounts()),
		}

		if payTime := parseOrderTime(eventOrder.GetPayTime()); payTime != nil {
			order.PayTime = payTime
		}

		orders = append(orders, order)
	}

	return orders
}

func toOrderDetails(tid string, eventDetails []event.GoodsDetail) []*OrderDetail {
	if len(eventDetails) == 0 {
		return []*OrderDetail{}
	}

	details := make([]*OrderDetail, 0, len(eventDetails))
	for _, eventDetail := range eventDetails {
		if eventDetail == nil {
			continue
		}

		details = append(details, &OrderDetail{
			Tid:            tid,
			OID:            eventDetail.GetOid(),
			Barcode:        eventDetail.GetBarcode(),
			EShopGoodsID:   eventDetail.GetEshopGoodsId(),
			OuterIID:       eventDetail.GetOuterIid(),
			EShopGoodsName: eventDetail.GetEshopGoodsName(),
			EShopSkuID:     eventDetail.GetEshopSkuId(),
			EShopSkuName:   eventDetail.GetEshopSkuName(),
			NumIID:         eventDetail.GetNumIid(),
			SkuID:          eventDetail.GetSkuId(),
			Num:            eventDetail.GetNum(),
			Payment:        eventDetail.GetPayment(),
			PicPath:        eventDetail.GetPicPath(),
			Weight:         eventDetail.GetWeight(),
			Size:           eventDetail.GetSize(),
			UnitID:         eventDetail.GetUniTid(),
			UnitQty:        eventDetail.GetUnitQty(),
		})
	}

	return details
}

func toOrderAccounts(tid string, eventAccounts []event.Account) []*OrderAccount {
	if len(eventAccounts) == 0 {
		return []*OrderAccount{}
	}

	accounts := make([]*OrderAccount, 0, len(eventAccounts))
	for _, eventAccount := range eventAccounts {
		if eventAccount == nil {
			continue
		}

		accounts = append(accounts, &OrderAccount{
			Tid:         tid,
			FinanceCode: eventAccount.GetFinanceCode(),
			Total:       eventAccount.GetTotal(),
		})
	}

	return accounts
}

func parseOrderTime(value string) *time.Time {
	if value == "" {
		return nil
	}

	layouts := []string{"2006-01-02 15:04:05", time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return &parsed
		}
	}

	return nil
}
