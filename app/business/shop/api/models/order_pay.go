package models

import (
	"encoding/json"
	"nova-factory-server/app/constant/order"
	"time"
)

type WechatPayNotifyData struct {
	TransactionId string `json:"transaction_id"`
	Amount        struct {
		PayerTotal    int    `json:"payer_total"`
		Total         int    `json:"total"`
		Currency      string `json:"currency"`
		PayerCurrency string `json:"payer_currency"`
	} `json:"amount"`
	Mchid           string `json:"mchid"`
	TradeState      string `json:"trade_state"`
	BankType        string `json:"bank_type"`
	PromotionDetail []struct {
		Amount              int    `json:"amount"`
		WechatpayContribute int    `json:"wechatpay_contribute"`
		CouponId            string `json:"coupon_id"`
		Scope               string `json:"scope"`
		MerchantContribute  int    `json:"merchant_contribute"`
		Name                string `json:"name"`
		OtherContribute     int    `json:"other_contribute"`
		Currency            string `json:"currency"`
		StockId             string `json:"stock_id"`
		GoodsDetail         []struct {
			GoodsRemark    string `json:"goods_remark"`
			Quantity       int    `json:"quantity"`
			DiscountAmount int    `json:"discount_amount"`
			GoodsId        string `json:"goods_id"`
			UnitPrice      int    `json:"unit_price"`
		} `json:"goods_detail"`
	} `json:"promotion_detail"`
	SuccessTime time.Time `json:"success_time"`
	Payer       struct {
		Openid string `json:"openid"`
	} `json:"payer"`
	OutTradeNo     string `json:"out_trade_no"`
	Appid          string `json:"appid"`
	TradeStateDesc string `json:"trade_state_desc"`
	TradeType      string `json:"trade_type"`
	Attach         string `json:"attach"`
	SceneInfo      struct {
		DeviceId string `json:"device_id"`
	} `json:"scene_info"`
}

func (p *WechatPayNotifyData) GetOutTradeNo() string {
	return p.OutTradeNo
}
func (p *WechatPayNotifyData) GetTransactionID() string {
	return p.TransactionId
}
func (p *WechatPayNotifyData) GetTradeState() string {
	return p.TradeState
}
func (p *WechatPayNotifyData) GetAppID() string {
	return p.Appid
}
func (p *WechatPayNotifyData) GetMchID() string {
	return p.Mchid
}
func (p *WechatPayNotifyData) GetAmount() int64 {
	return int64(p.Amount.Total)
}
func (p *WechatPayNotifyData) GetPayerOpenID() string {
	return p.Payer.Openid
}
func (p *WechatPayNotifyData) GetNotifyRaw() string {
	rawBytes, _ := json.Marshal(p)
	notifyRaw := string(rawBytes)
	return notifyRaw
}

func (p *WechatPayNotifyData) GetPayChannel() int {
	return order.PayChannelWechat
}
