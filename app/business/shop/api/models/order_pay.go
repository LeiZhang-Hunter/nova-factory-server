package models

type PayNotifyData struct {
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
}
