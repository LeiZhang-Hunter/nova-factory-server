package order

// DelayCancelKey 订单超时自动取消延迟队列（Redis Sorted Set）。
// Score 为订单应取消的 Unix 时间戳（秒），Member 为订单 ID 字符串。
const DelayCancelKey string = "shop:order:delay_cancel"
