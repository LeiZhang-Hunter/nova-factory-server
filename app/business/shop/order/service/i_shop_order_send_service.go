package service

import (
	"nova-factory-server/app/utils/observer/integration/event"
)

// IShopOrderSendService 订单发货业务接口。
//
// 发货属于业务编排逻辑，service 层负责：
// - 在一个事务里完成发货主表和明细表的组合写入；
// - 参数校验与数据规范化。
//
// DAO 层只提供单表数据操作能力，service 层负责组合多个 DAO 完成一次完整写入。
type IShopOrderSendService interface {
	// Set 存储订单发货数据（主表 + 明细）。
	//
	// 实现层会在一个事务中先写入 shop_order_send 主表，
	// 再批量写入 shop_order_send_detail 明细。
	// 任一步失败时返回 error 并触发事务回滚。
	Set(sendEvent event.OrderSendEvent) error
}
