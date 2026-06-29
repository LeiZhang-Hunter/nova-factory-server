package models

import (
	"time"
)

// WechatRefundNotifyData 微信退款通知数据结构
type WechatRefundNotifyData struct {
	Mchid               string    `json:"mchid"`
	TransactionId       string    `json:"transaction_id"`
	OutTradeNo          string    `json:"out_trade_no"`
	RefundId            string    `json:"refund_id"`
	OutRefundNo         string    `json:"out_refund_no"`
	RefundStatus        string    `json:"refund_status"`
	SuccessTime         time.Time `json:"success_time"`
	UserReceivedAccount string    `json:"user_received_account"`
	Amount              struct {
		Total       int `json:"total"`
		Refund      int `json:"refund"`
		PayerTotal  int `json:"payer_total"`
		PayerRefund int `json:"payer_refund"`
	} `json:"amount"`
}

func (d *WechatRefundNotifyData) GetOutTradeNo() string    { return d.OutTradeNo }
func (d *WechatRefundNotifyData) GetOutRefundNo() string   { return d.OutRefundNo }
func (d *WechatRefundNotifyData) GetRefundID() string      { return d.RefundId }
func (d *WechatRefundNotifyData) GetRefundStatus() string  { return d.RefundStatus }
func (d *WechatRefundNotifyData) GetTransactionID() string { return d.TransactionId }

// RefundApplyReq 售后申请请求
type RefundApplyReq struct {
	OrderID int64  `json:"order_id,string" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
}

// RefundApplyResp 售后申请响应
type RefundApplyResp struct {
	ID          int64  `json:"id,string"`
	OutRefundNo string `json:"out_refund_no"`
	Status      int32  `json:"status"`
	StatusText  string `json:"status_text"`
}

//// RefundDetailResp 售后单详情
//type RefundDetailResp struct {
//	ID                 uint64  `json:"id,string"`
//	OrderID            uint64  `json:"order_id,string"`
//	Tid                string  `json:"tid"`
//	PayChannel         int     `json:"pay_channel"`
//	OutRefundNo        string  `json:"out_refund_no"`
//	RefundAmount       float64 `json:"refund_amount"`
//	TotalAmount        float64 `json:"total_amount"`
//	Reason             string  `json:"reason"`
//	Status             int32   `json:"status"`
//	StatusText         string  `json:"status_text"`
//	AuditRemark        string  `json:"audit_remark"`
//	ThirdTransactionID string  `json:"third_transaction_id"`
//	ThirdRefundID      string  `json:"third_refund_id"`
//	SyncStatus         int32   `json:"sync_status"`
//	SyncMessage        string  `json:"sync_message"`
//	CreateTime         string  `json:"create_time"`
//}
//
//// RefundListResp 售后单列表
//type RefundListResp struct {
//	Rows  []*RefundDetailResp `json:"rows"`
//	Total int64               `json:"total"`
//}
