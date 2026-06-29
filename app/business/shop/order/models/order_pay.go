package models

import "fmt"

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

// RefundNotifyDataInterface 退款回调通知数据接口（通道级验签/解密之后的可信数据）
type RefundNotifyDataInterface interface {
	GetOutTradeNo() string
	GetOutRefundNo() string
	GetRefundID() string
	GetRefundStatus() string
	GetTransactionID() string
}

// GetPayNotifyResult 泛型辅助函数，将 PayNotifyDataInterface 断言为具体类型。
func GetPayNotifyResult[T PayNotifyDataInterface](r PayNotifyDataInterface) (T, error) {
	t, ok := r.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("不支持的支付通知类型: got %T, want %T", r, zero)
	}
	return t, nil
}

// GetRefundNotifyResult 泛型辅助函数，将 RefundNotifyDataInterface 断言为具体类型。
func GetRefundNotifyResult[T RefundNotifyDataInterface](r RefundNotifyDataInterface) (T, error) {
	t, ok := r.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("不支持的退款通知类型: got %T, want %T", r, zero)
	}
	return t, nil
}
