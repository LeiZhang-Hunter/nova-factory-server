package models

// RefundRequest 退款请求参数
type RefundRequest struct {
	Tid           string  `json:"tid"`
	TransactionID string  `json:"transaction_id"`
	OutRefundNo   string  `json:"out_refund_no"`
	Reason        string  `json:"reason"`
	RefundAmount  float64 `json:"refund_amount"`
	TotalAmount   float64 `json:"total_amount"`
	NotifyURL     string  `json:"notify_url"`
	FundsAccount  string  `json:"funds_account"`
	MchID         string  `json:"mch_id"`
	AppID         string  `json:"appid"`
}

// RefundResult 退款结果
type RefundResult struct {
	ThirdRefundID      string `json:"third_refund_id"`
	ThirdTransactionID string `json:"third_transaction_id"`
	Status             string `json:"status"`
}
