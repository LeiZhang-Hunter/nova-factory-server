package models

import (
	"encoding/json"
	"nova-factory-server/app/constant/order"
)

type WechatPayNotifyData struct {
	OutTradeNo    string `json:"out_trade_no"`
	TransactionID string `json:"transaction_id"`
	TradeState    string `json:"trade_state"`
	AppID         string `json:"appid"`
	MchID         string `json:"mchid"`
	Amount        struct {
		Total int64 `json:"total"`
	} `json:"amount"`
	Payer struct {
		Openid string `json:"openid"`
	} `json:"payer"`
	PayChannel int
}

func (p *WechatPayNotifyData) GetOutTradeNo() string {
	return p.OutTradeNo
}
func (p *WechatPayNotifyData) GetTransactionID() string {
	return p.TransactionID
}
func (p *WechatPayNotifyData) GetTradeState() string {
	return p.TradeState
}
func (p *WechatPayNotifyData) GetAppID() string {
	return p.AppID
}
func (p *WechatPayNotifyData) GetMchID() string {
	return p.MchID
}
func (p *WechatPayNotifyData) GetAmount() int64 {
	return p.Amount.Total
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
