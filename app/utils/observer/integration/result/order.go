// 定义订单同步的响应接口。
// 包含批量订单同步的整体响应（OrderSyncResponse）及
// 每笔订单的独立同步结果（OrderSyncResult），
// 用于对接第三方 ERP 的订单推送接口返回。
package result

// OrderSyncResponse 订单同步响应接口，表示一批订单推送至第三方系统后的整体返回。
// 包含整体状态码、消息及每笔订单的独立处理结果。
type OrderSyncResponse interface {
	// GetCode 返回整体业务状态码，0 表示成功
	GetCode() int64
	// GetMessage 返回整体业务消息
	GetMessage() string
	// GetOrders 返回每笔订单的同步结果列表
	GetOrders() []OrderSyncResult

	base
}

// OrderSyncResult 单笔订单同步结果接口，描述一条订单推送的处理状态。
type OrderSyncResult interface {
	// GetTid 网店订单号
	GetTid() string
	// GetBillCode ERP 中生成的单据编码
	GetBillCode() string
	// GetMessage 该笔订单的处理消息
	GetMessage() string

	base
}
