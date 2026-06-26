package models

type PayNotifyDataInterface interface {
	GetOutTradeNo() string
	GetTransactionID() string
	GetTradeState() string
	GetAppID() string
	GetMchID() string
	GetAmount() int64
	GetPayerOpenID() string
	GetNotifyRaw() string
	GetPayChannel() int
}
