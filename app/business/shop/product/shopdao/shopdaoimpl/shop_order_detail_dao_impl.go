package shopdaoimpl

import (
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/commonStatus"
	objectutil "nova-factory-server/app/utils/json"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ShopOrderDetailDaoImpl 商城订单明细 DAO。
//
// 该 DAO 只负责 shop_order_detail 表的写入和删除。
// 它不持有 *gorm.DB，也不自己开启事务；所有方法都接收外部传入的 tx。
// 这样可以确保订单主表、明细表、账户表共用同一个事务，任一步失败都会由外层事务统一回滚。
type ShopOrderDetailDaoImpl struct {
	tableName string
}

// NewShopOrderDetailDaoImpl 创建订单明细 DAO。
func NewShopOrderDetailDaoImpl() shopdao.IShopOrderDetailDao {
	return &ShopOrderDetailDaoImpl{tableName: "shop_order_detail"}
}

// BatchCreate 批量创建订单明细。
//
// 参数说明：
// - tx：外层订单同步事务，必须由主订单 DAO 传入；
// - orderID：shop_order 主表 ID，用于填充 shop_order_detail.order_id；
// - order：当前同步的订单模型，明细数据来自 order.Details；
// - now：本次同步时间，用于补齐 CreateTime/UpdateTime。
//
// 写入前会对每条明细做规范化：
// 1. ID 置 0，避免上游传入的旧 ID 影响新增；
// 2. OrderID/Tid 以当前主表为准，保证子表归属正确；
// 3. OID 去掉首尾空格，避免同一明细因为空格导致唯一键异常；
// 4. CreateTime/UpdateTime 缺失时使用同步时间；
// 5. State 固定为 commonStatus.NORMAL。
func (d *ShopOrderDetailDaoImpl) BatchCreate(tx *gorm.DB, orderID uint64, order *shopmodels.Order, now *time.Time) error {
	if len(order.Details) == 0 {
		return nil
	}

	rows := make([]*shopmodels.OrderDetail, 0, len(order.Details))
	for _, detail := range order.Details {
		if detail == nil {
			continue
		}
		detail.ID = 0
		detail.OrderID = orderID
		detail.Tid = order.Tid
		detail.OID = strings.TrimSpace(detail.OID)
		detail.CreateTime = objectutil.FirstTime(detail.CreateTime, now)
		detail.UpdateTime = objectutil.FirstTime(detail.UpdateTime, now)
		detail.State = commonStatus.NORMAL
		rows = append(rows, detail)
	}
	if len(rows) == 0 {
		return nil
	}
	return tx.Table(d.tableName).Create(&rows).Error
}

// DeleteByOrderID 删除指定订单的旧明细。
//
// 当前同步策略是“主表更新 + 子表重建”：
// 已存在订单更新时，先删除该 order_id 下的旧明细，再插入本次事件带来的最新明细。
// 该方法必须使用外层 tx；如果后续插入新明细或账户失败，删除操作也会随事务回滚。
func (d *ShopOrderDetailDaoImpl) DeleteByOrderID(tx *gorm.DB, orderID uint64) error {
	return tx.Table(d.tableName).Where("order_id = ?", orderID).Delete(&shopmodels.OrderDetail{}).Error
}
