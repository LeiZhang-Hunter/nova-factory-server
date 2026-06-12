package shopdaoimpl

import (
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/commonStatus"
	objectutil "nova-factory-server/app/utils/json"
	"time"

	"gorm.io/gorm"
)

// ShopOrderAccountDaoImpl 商城订单账户 DAO。
//
// 该 DAO 只负责 shop_order_account 表的写入和删除。
// 和订单明细 DAO 一样，它不单独开启事务，所有数据库操作都使用主订单 DAO 传入的 tx。
// 这样可以避免订单主表已更新但账户表没有同步成功的问题。
type ShopOrderAccountDaoImpl struct {
	tableName string
}

// NewShopOrderAccountDaoImpl 创建订单账户 DAO。
func NewShopOrderAccountDaoImpl() shopdao.IShopOrderAccountDao {
	return &ShopOrderAccountDaoImpl{tableName: "shop_order_account"}
}

// BatchCreate 批量创建订单账户记录。
//
// 参数说明：
// - tx：外层订单同步事务；
// - orderID：shop_order 主表 ID，用于填充 shop_order_account.order_id；
// - order：当前同步的订单模型，账户数据来自 order.Accounts；
// - now：本次同步时间，用于补齐 CreateTime/UpdateTime。
//
// 写入前会对每条账户记录做规范化：
// 1. ID 置 0，确保执行新增；
// 2. OrderID/Tid 以当前主表为准；
// 3. CreateTime/UpdateTime 缺失时使用同步时间；
// 4. State 固定为 commonStatus.NORMAL。
func (d *ShopOrderAccountDaoImpl) BatchCreate(tx *gorm.DB, orderID uint64, order *shopmodels.Order, now *time.Time) error {
	if len(order.Accounts) == 0 {
		return nil
	}

	rows := make([]*shopmodels.OrderAccount, 0, len(order.Accounts))
	for _, account := range order.Accounts {
		if account == nil {
			continue
		}
		account.ID = 0
		account.OrderID = orderID
		account.Tid = order.Tid
		account.CreateTime = objectutil.FirstTime(account.CreateTime, now)
		account.UpdateTime = objectutil.FirstTime(account.UpdateTime, now)
		account.State = commonStatus.NORMAL
		rows = append(rows, account)
	}
	if len(rows) == 0 {
		return nil
	}
	return tx.Table(d.tableName).Create(&rows).Error
}

// DeleteByOrderID 删除指定订单的旧账户记录。
//
// 已存在订单同步时，账户记录采用重建策略：先删旧账户，再插入本次同步的新账户。
// 该方法只使用传入的 tx，不自行提交；如果后续步骤失败，删除会随外层事务一起回滚。
func (d *ShopOrderAccountDaoImpl) DeleteByOrderID(tx *gorm.DB, orderID uint64) error {
	return tx.Table(d.tableName).Where("order_id = ?", orderID).Delete(&shopmodels.OrderAccount{}).Error
}
