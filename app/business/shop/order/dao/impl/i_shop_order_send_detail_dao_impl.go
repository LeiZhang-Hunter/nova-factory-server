package impl

import (
	shopdao "nova-factory-server/app/business/shop/order/dao"
	shopmodels "nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/constant/commonStatus"

	"gorm.io/gorm"
)

// ShopOrderSendDetailDaoImpl 订单发货明细 DAO。
//
// 该 DAO 只负责 shop_order_send_detail 表的批量写入。
// 它不持有 *gorm.DB，也不自己开启事务；所有方法都接收外部传入的 tx。
// 这样可以确保发货主表、明细表共用同一个事务，任一步失败都会由外层事务统一回滚。
type ShopOrderSendDetailDaoImpl struct {
	tableName string
}

// NewShopOrderSendDetailDaoImpl 创建发货明细 DAO。
func NewShopOrderSendDetailDaoImpl() shopdao.IShopOrderSendDetailDao {
	return &ShopOrderSendDetailDaoImpl{tableName: "shop_order_send_detail"}
}

// BatchCreate 批量创建发货明细。
//
// 写入前对每条明细做规范化：
// 1. ID 置 0，避免上游传入的旧 ID 影响新增；
// 2. SendID 以当前发货主表 ID 为准；
// 3. State 固定为 commonStatus.NORMAL。
func (d *ShopOrderSendDetailDaoImpl) BatchCreate(tx *gorm.DB, sendID uint64, details []*shopmodels.OrderSendDetail) error {
	if len(details) == 0 {
		return nil
	}

	rows := make([]*shopmodels.OrderSendDetail, 0, len(details))
	for _, detail := range details {
		if detail == nil {
			continue
		}
		detail.ID = 0
		detail.SendID = sendID
		detail.State = commonStatus.NORMAL
		rows = append(rows, detail)
	}
	if len(rows) == 0 {
		return nil
	}
	return tx.Table(d.tableName).Create(&rows).Error
}
