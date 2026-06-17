package dao

import (
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/order/models"
)

// IShopOrderSendDetailDao 订单发货明细表数据访问接口。
//
// 只负责 shop_order_send_detail 单表操作。
// 调用方必须传入外层事务 tx，保证明细表操作能和发货主表一起提交或回滚。
type IShopOrderSendDetailDao interface {
	// BatchCreate 批量创建发货明细。
	//
	// 参数说明：
	// - tx：外层事务；
	// - sendID：发货主表 ID，用于填充明细 send_id；
	// - details：本次写入的发货明细列表。
	BatchCreate(tx *gorm.DB, sendID uint64, details []*models.OrderSendDetail) error
}
