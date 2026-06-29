package order

import (
	"nova-factory-server/app/utils/observer/integration/event"
)

type orderDetail struct {
	Oid            string  `json:"oid"`
	EShopGoodsName string  `json:"eshopgoodsname"`
	EShopSKUName   string  `json:"eshopskuname"`
	Numiid         int64   `json:"numiid"`
	Skuid          int64   `json:"skuid"`
	Num            float64 `json:"num"`
	Outeriid       string  `json:"outeriid"`
	Payment        float64 `json:"payment"`
	Picpath        string  `json:"picpath"`
	Weight         float64 `json:"weight"`
	Size           float64 `json:"size"`
	Unitid         int64   `json:"unitid"`
	Unitqty        float64 `json:"unitqty"`
}

type orderAccount struct {
	FinanceCode string  `json:"financeCode"`
	Total       float64 `json:"total"`
}

type syncOrder struct {
	Tid              string          `json:"tid"`
	Weight           float64         `json:"weight"`
	Size             float64         `json:"size"`
	Buyernick        string          `json:"buyernick"`
	Buyermessage     string          `json:"buyermessage"`
	Sellermemo       string          `json:"sellermemo"`
	Total            float64         `json:"total"`
	Postfee          float64         `json:"postfee"`
	Receivername     string          `json:"receivername"`
	Receiverstate    string          `json:"receiverstate"`
	Receivercity     string          `json:"receivercity"`
	Receiverdistrict string          `json:"receiverdistrict"`
	Receiveraddress  string          `json:"receiveraddress"`
	Receiverphone    string          `json:"receiverphone"`
	Receivermobile   string          `json:"receivermobile"`
	Receiverzip      string          `json:"receiverzip"`
	Created          string          `json:"created"`
	Status           string          `json:"status"`
	Type             string          `json:"type"`
	Invoicename      string          `json:"invoicename"`
	Sellerflag       string          `json:"sellerflag"`
	Paytime          string          `json:"paytime"`
	Btypecode        string          `json:"btypecode"`
	Getbtypeaddress  int             `json:"getbtypeaddress"`
	Details          []*orderDetail  `json:"details"`
	Accounts         []*orderAccount `json:"accounts"`
}
type syncAfterSaleOrders struct {
	Rtid           string                       `json:"rtid"`
	Tid            string                       `json:"tid"`
	Logistbillcode string                       `json:"logistbillcode"`
	Postfee        float64                      `json:"postfee"`
	Total          float64                      `json:"total"`
	Privilege      float64                      `json:"privilege"`
	Created        string                       `json:"created"`
	Aftsaletype    string                       `json:"aftsaletype"`
	Reasoncode     string                       `json:"reasoncode"`
	Aftsaleremark  string                       `json:"aftsaleremark"`
	Details        []syncAfterSaleOrdersDetail  `json:"details"`
	Exdetails      []syncAfterSaleOrderExdetail `json:"exdetails"`
}
type syncAfterSaleOrdersDetail struct {
	Oid            string  `json:"oid"`
	Backqty        float64 `json:"backqty"`
	Backtotal      float64 `json:"backtotal"`
	Outeriid       string  `json:"outeriid"`
	Eshopgoodsname string  `json:"eshopgoodsname"`
	Eshopskuname   string  `json:"eshopskuname"`
}
type syncAfterSaleOrderExdetail struct {
	Oid            string  `json:"oid"`
	Exchangeqty    float64 `json:"exchangeqty"`
	Backtotal      float64 `json:"backtotal"`
	Outeriid       string  `json:"outeriid"`
	Eshopgoodsname string  `json:"eshopgoodsname"`
	Eshopskuname   string  `json:"eshopskuname"`
}

// syncOrder 同步订单
type syncOrderList struct {
	Orders []syncOrder `json:"orders"`
}

func toOrderSyncOrder(orderInfo event.OrderEvent) *syncOrderList {
	if orderInfo == nil {
		return nil
	}

	orders := make([]syncOrder, 0)
	for _, v := range orderInfo.GetOrders() {
		orderData := &syncOrder{
			Tid:          v.GetOrderNo(),
			Weight:       v.GetWeight(),
			Size:         v.GetSize(),
			Buyernick:    v.GetBuyerNick(),
			Buyermessage: v.GetBuyerMessage(),
			Sellermemo:   v.GetSellerMemo(),
			Total:        v.GetTotalAmount(),
			Postfee:      v.GetPostFee(),
			Receivername: v.GetReceiverName(),

			Receiverstate:    v.GetReceiverStateName(),
			Receivercity:     v.GetReceiverCityName(),
			Receiverdistrict: v.GetReceiverDistrictName(),

			Receiveraddress: v.GetReceiverAddress(),
			Receiverphone:   v.GetReceiverPhone(),
			Receivermobile:  v.GetReceiverMobile(),
			Created:         v.GetCreated(),
			Status:          v.GetStatus(),
			Type:            v.GetType(),
			Invoicename:     v.GetInvoiceName(),
			Sellerflag:      v.GetSellerFlag(),
			Paytime:         v.GetPayTime(),
			Btypecode:       v.GetBTypeCode(),
			Details:         toOrderSyncDetails(v.GetDetails()),
			Accounts:        toOrderSyncAccounts(v.GetAccounts()),
		}
		orders = append(orders, *orderData)
	}

	return &syncOrderList{
		Orders: orders,
	}
}

func toOrderSyncDetails(details []event.GoodsDetail) []*orderDetail {
	if len(details) == 0 {
		return []*orderDetail{}
	}
	result := make([]*orderDetail, 0, len(details))
	for _, detail := range details {
		if detail == nil {
			continue
		}
		result = append(result, &orderDetail{
			Oid:            detail.GetOid(),
			Numiid:         detail.GetNumIid(),
			Skuid:          detail.GetSkuId(),
			Num:            detail.GetNum(),
			EShopGoodsName: detail.GetEshopGoodsName(),
			EShopSKUName:   detail.GetEshopSkuName(),
			Outeriid:       detail.GetOuterIid(),
			Payment:        detail.GetPayment(),
			Picpath:        detail.GetPicPath(),
			Weight:         detail.GetWeight(),
			Size:           detail.GetSize(),
			Unitid:         detail.GetUniTid(),
			Unitqty:        detail.GetUnitQty(),
		})
	}
	return result
}

func toOrderSyncAccounts(accounts []event.Account) []*orderAccount {
	result := make([]*orderAccount, 0, len(accounts))
	if len(accounts) == 0 {
		return result
	}
	for _, account := range accounts {
		if account == nil {
			continue
		}
		result = append(result, &orderAccount{
			FinanceCode: account.GetFinanceCode(),
			Total:       account.GetTotal(),
		})
	}
	return result
}

func toSyncAfterSaleOrders(reqData event.ZAfterSaleOrderSyncReqDataEvent) []syncAfterSaleOrders {
	result := make([]syncAfterSaleOrders, 0)
	for _, v := range *reqData.GetOrders() {
		syncAfterSaleOrderExdetails := make([]syncAfterSaleOrderExdetail, 0)
		if v.GetExDetails() != nil {
			for _, dv := range *v.GetExDetails() {
				syncAfterSaleOrderExdetails = append(syncAfterSaleOrderExdetails, syncAfterSaleOrderExdetail{
					Oid:            dv.GetOid(),
					Exchangeqty:    dv.GetExchangeQty(),
					Backtotal:      dv.GetBackTotal(),
					Outeriid:       dv.GetOuterIid(),
					Eshopgoodsname: dv.GetEshopGoodsName(),
					Eshopskuname:   dv.GetEshopSkuName(),
				})
			}
		}
		syncAfterSaleOrdersDetails := make([]syncAfterSaleOrdersDetail, 0)
		if v.GetDetails() != nil {
			for _, dv := range *v.GetDetails() {
				syncAfterSaleOrdersDetails = append(syncAfterSaleOrdersDetails, syncAfterSaleOrdersDetail{
					Oid:            dv.GetOid(),
					Backqty:        dv.GetBackQty(),
					Backtotal:      dv.GetBackTotal(),
					Outeriid:       dv.GetOuterIid(),
					Eshopgoodsname: dv.GetEshopGoodsName(),
					Eshopskuname:   dv.GetEshopSkuName(),
				})
			}
		}

		result = append(result, syncAfterSaleOrders{
			Rtid:           v.GetRtid(),
			Tid:            v.GetTid(),
			Logistbillcode: v.GetLogistBillCode(),
			Postfee:        v.GetPostFee(),
			Total:          v.GetTotal(),
			Privilege:      v.GetPrivilege(),
			Created:        v.GetCreated(),
			Aftsaletype:    v.GetAftSaleType(),
			Reasoncode:     v.GetReasonCode(),
			Aftsaleremark:  v.GetAftSaleRemark(),
			Details:        syncAfterSaleOrdersDetails,
			Exdetails:      syncAfterSaleOrderExdetails,
		})
	}
	return result
}
