package service

import "context"

// IOrderTimeoutService 订单超时自动取消服务。
//
// Consumer 与 HTTP 兜底接口共享同一套取消逻辑（ProcessExpiredOrders），
// 仅入口不同：Consumer 为常驻 goroutine，HTTP 接口由外部定时任务驱动。
type IOrderTimeoutService interface {
	// StartConsumer 启动后台 goroutine，每秒轮询 Redis 延迟队列并取消到期订单。
	// 通过 ctx 取消可优雅停止 Consumer。
	StartConsumer(ctx context.Context)

	// ProcessExpiredOrders 扫描延迟队列中所有到期订单并执行取消，返回被取消的订单数量。
	ProcessExpiredOrders(ctx context.Context) (int, error)
}
