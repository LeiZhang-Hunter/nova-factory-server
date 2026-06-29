package models

import (
	"fmt"

	orderConstant "nova-factory-server/app/constant/order"
)

// PrepayRequest 预下单请求，只保留业务输入。
// AppID/MchID/NotifyUrl/私钥等通道配置由 Provider 内部读取。
type PrepayRequest struct {
	Tid         string  // 订单号
	TotalAmount float64 // 金额（元）
	Description string  // 商品描述
	Openid      string  // 用户 openid（支付平台侧）
}

// PrepayResult 预下单结果接口，每种支付通道返回各自的实现。
type PrepayResult interface {
	Channel() int
}

// GetPrepayResult 泛型辅助函数，将 PrepayResult 断言为具体类型。
func GetPrepayResult[T PrepayResult](r PrepayResult) (T, error) {
	t, ok := r.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("不支持的支付结果类型: got %T, want %T", r, zero)
	}
	return t, nil
}

// WechatPrepayResult 微信 JSAPI 预下单结果
type WechatPrepayResult struct {
	PrepayID  string
	AppId     string
	TimeStamp string
	NonceStr  string
	Package   string
	SignType  string
	PaySign   string
}

// Channel 实现 PrepayResult 接口
func (r *WechatPrepayResult) Channel() int { return orderConstant.PayChannelWechat }

// QueryOrderResult 商户订单号查询订单结果（微信支付 /v3/pay/transactions/out-trade-no/{out_trade_no} 应答）
type QueryOrderResult struct {
	Appid            string `json:"appid"`
	Mchid            string `json:"mchid"`
	OutTradeNo       string `json:"out_trade_no"`
	TransactionId    string `json:"transaction_id"`
	TradeType        string `json:"trade_type"`
	TradeState       string `json:"trade_state"`
	TradeStateDesc   string `json:"trade_state_desc"`
	BankType         string `json:"bank_type"`
	Attach           string `json:"attach"`
	SuccessTime      string `json:"success_time"`
	PayerOpenid      string `json:"payer_openid"`
	AmountTotal      int    `json:"amount_total"`
	AmountPayerTotal int    `json:"amount_payer_total"`
	AmountCurrency   string `json:"amount_currency"`
}

// QueryRefundResult 商户退款单号查询退款结果（微信支付 /v3/refund/domestic/refunds/{out_refund_no} 应答）
type QueryRefundResult struct {
	RefundId            string `json:"refund_id"`
	OutRefundNo         string `json:"out_refund_no"`
	TransactionId       string `json:"transaction_id"`
	OutTradeNo          string `json:"out_trade_no"`
	Channel             string `json:"channel"`
	UserReceivedAccount string `json:"user_received_account"`
	SuccessTime         string `json:"success_time"`
	CreateTime          string `json:"create_time"`
	Status              string `json:"status"`
	FundsAccount        string `json:"funds_account"`
	AmountTotal         int    `json:"amount_total"`
	AmountRefund        int    `json:"amount_refund"`
	AmountPayerTotal    int    `json:"amount_payer_total"`
	AmountPayerRefund   int    `json:"amount_payer_refund"`
	AmountCurrency      string `json:"amount_currency"`
}
