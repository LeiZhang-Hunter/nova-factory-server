package shopservice

import "nova-factory-server/app/utils/observer/integration/event"

// IShopOrderService 商城订单同步业务接口。
//
// 订单同步属于业务编排逻辑，应放在 service 层处理，而不是放在 DAO 层：
// - event.OrderEvent 到 shopmodels.Order 的转换；
// - 按 tid 判断订单新增或更新；
// - 订单状态是否允许覆盖的校验；
// - 主表、明细表、账户表多个 DAO 的调用顺序；
// - 事务边界和回滚策略。
//
// DAO 层只提供单表数据操作能力，service 层负责组合多个 DAO 完成一次完整同步。
type IShopOrderService interface {
	// SyncOrder 同步订单事件。
	//
	// 实现层会把事件转为商城订单模型，并在一个事务里完成主表和子表同步。
	// 任一 DAO 操作返回错误时，service 会返回该错误并触发事务回滚。
	SyncOrder(event event.OrderEvent) error
}
