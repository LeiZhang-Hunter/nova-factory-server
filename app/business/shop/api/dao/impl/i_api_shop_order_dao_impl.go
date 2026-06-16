package impl

import (
	"errors"
	"fmt"
	"nova-factory-server/app/business/datasyncapi/gjpqqd/models"
	apimodels "nova-factory-server/app/business/shop/api/models"
	shopordermodels "nova-factory-server/app/business/shop/order/models"
	shopusermodels "nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/utils/stringUtils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/constant/commonStatus"
	orderConstant "nova-factory-server/app/constant/order"
)

type IApiShopOrderDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewIApiShopOrderDaoImpl 初始化商城 API 订单 DAO。
func NewIApiShopOrderDaoImpl(db *gorm.DB) dao.IApiShopOrderDao {
	return &IApiShopOrderDaoImpl{
		db:    db,
		table: "shop_order",
	}
}

// ListShopOrders 分页查询当前商城用户在 shop_order 表中的订单列表。
func (i *IApiShopOrderDaoImpl) ListShopOrders(c *gin.Context, shopUser *shopusermodels.User, query *apimodels.OrderQuery) (*apimodels.OrderListData, error) {
	if i == nil || i.db == nil {
		return nil, errors.New("数据库连接不存在")
	}
	if shopUser == nil {
		return nil, errors.New("商城用户不存在")
	}
	if query == nil {
		query = &models.OrderQuery{}
	}

	db := i.db.WithContext(c).
		Table(i.table).
		Where("state = ?", commonStatus.NORMAL).
		Where("buyer_nick = ?", buildOrderBuyerNick(shopUser))
	if query.Status != nil {
		db = db.Where("status = ?", orderConstant.ShopStatusToOrderStatus(*query.Status))
	}
	if strings.TrimSpace(query.OrderNo) != "" {
		db = db.Where("tid LIKE ?", "%"+strings.TrimSpace(query.OrderNo)+"%")
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	size := query.Size
	if size < 1 {
		size = 10
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*shopordermodels.Order, 0)
	if err := db.Order("id DESC").
		Offset(int((page - 1) * size)).
		Limit(int(size)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	if err := i.attachShopOrderDetails(c, rows); err != nil {
		return nil, err
	}

	data := &apimodels.OrderListData{
		Rows:  make([]*apimodels.OrderVO, 0, len(rows)),
		Total: total,
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		data.Rows = append(data.Rows, toShopOrderVO(row))
	}
	return data, nil
}

// DeleteByIDs 删除 ERP 订单。
func (i *IApiShopOrderDaoImpl) DeleteByIDs(c *gin.Context, ids []uint64) error {
	return nil
}

// UpdateERPOrderStatus 更新当前商城用户的 shop_order 订单状态。
func (i *IApiShopOrderDaoImpl) UpdateERPOrderStatus(c *gin.Context, id int64, shopUser *shopusermodels.User, status int32) (int64, error) {
	result := i.db.Table(i.table).
		Where("id = ?", id).
		Where("buyer_nick = ?", buildOrderBuyerNick(shopUser)).
		Updates(map[string]interface{}{
			"status":      orderConstant.ShopStatusToOrderStatus(status),
			"update_time": gorm.Expr("NOW()"),
		})
	return result.RowsAffected, result.Error
}

// CancelERPOrder 将当前商城用户的 shop_order 订单标记为已取消。
func (i *IApiShopOrderDaoImpl) CancelERPOrder(c *gin.Context, id int64, shopUser *models.User, reason string) (int64, error) {
	result := i.db.Table(i.table).
		Where("id = ?", id).
		Where("buyer_nick = ?", buildOrderBuyerNick(shopUser)).
		Updates(map[string]interface{}{
			"status":      orderConstant.ShopStatusToOrderStatus(orderConstant.OrderStatusCancelled),
			"seller_memo": strings.TrimSpace(reason),
			"update_time": gorm.Expr("NOW()"),
		})
	return result.RowsAffected, result.Error
}

// GetERPOrderStatistics 统计当前商城用户的 shop_order 订单状态数量。
func (i *IApiShopOrderDaoImpl) GetERPOrderStatistics(c *gin.Context, shopUser *models.User) (*models.OrderStatistics, error) {
	stats := &models.OrderStatistics{}
	baseQuery := i.db.Table(i.table).Where("buyer_nick = ?", buildOrderBuyerNick(shopUser))
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", orderConstant.ShopStatusToOrderStatus(orderConstant.OrderStatusPending)).
		Count(&stats.PendingPay).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", orderConstant.ShopStatusToOrderStatus(orderConstant.OrderStatusPaid)).
		Count(&stats.PendingSend).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", orderConstant.ShopStatusToOrderStatus(orderConstant.OrderStatusShipped)).
		Count(&stats.PendingReceive).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", orderConstant.ShopStatusToOrderStatus(orderConstant.OrderStatusCompleted)).
		Count(&stats.Completed).Error; err != nil {
		return nil, err
	}
	if err := baseQuery.Session(&gorm.Session{}).
		Where("status = ?", orderConstant.ShopStatusToOrderStatus(orderConstant.OrderStatusCancelled)).
		Count(&stats.Cancelled).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

// MarkOrderPaidWithTx 在事务内锁定并标记 shop_order 为已支付。
func (i *IApiShopOrderDaoImpl) MarkOrderPaidWithTx(c *gin.Context, tx *gorm.DB, id uint64, payTime *time.Time, transactionId, notifyRaw, mchId, appid, payerOpenid string) error {
	if tx == nil {
		return errors.New("订单支付更新失败：事务不能为空")
	}
	var orderModel shopordermodels.Order
	if err := tx.WithContext(c).
		Table(i.table).
		Where("id = ?", id).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&orderModel).Error; err != nil {
		return fmt.Errorf("订单不存在: %v", err)
	}
	return tx.WithContext(c).
		Table(i.table).
		Where("id = ?", id).
		Where("status = ?", orderConstant.ERPStatusNoPay).
		Updates(map[string]interface{}{
			"status":         orderConstant.ERPStatusPayed,
			"pay_time":       payTime,
			"transaction_id": transactionId,
			"notify_raw":     notifyRaw,
			"mch_id":         mchId,
			"appid":          appid,
			"payer_openid":   payerOpenid,
			"pay_channel":    1,
		}).Error
}

func (o *IApiShopOrderDaoImpl) GetByTid(c *gin.Context, tid string) (*models.Order, error) {
	var row *models.Order
	if err := o.db.WithContext(c).Table(o.table).
		Where("tid = ?", strings.TrimSpace(tid)).
		Where("state = ?", commonStatus.NORMAL).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if err := o.attachChildren(c, []*models.Order{&item}); err != nil {
		return nil, err
	}
	return &item, nil
}

func (i *IApiShopOrderDaoImpl) attachShopOrderDetails(c *gin.Context, orders []*shopordermodels.Order) error {
	if len(orders) == 0 {
		return nil
	}
	orderIDs := make([]uint64, 0, len(orders))
	for _, order := range orders {
		if order == nil {
			continue
		}
		orderIDs = append(orderIDs, order.ID)
	}
	if len(orderIDs) == 0 {
		return nil
	}
	details := make([]*shopordermodels.OrderDetail, 0)
	if err := i.db.WithContext(c).
		Table("shop_order_detail").
		Where("order_id IN ?", orderIDs).
		Where("state = ?", commonStatus.NORMAL).
		Order("id ASC").
		Find(&details).Error; err != nil {
		return err
	}
	detailMap := make(map[uint64][]*shopordermodels.OrderDetail)
	for _, detail := range details {
		if detail == nil {
			continue
		}
		detailMap[detail.OrderID] = append(detailMap[detail.OrderID], detail)
	}
	for _, order := range orders {
		if order == nil {
			continue
		}
		order.Details = detailMap[order.ID]
	}
	return nil
}

func toShopOrderVO(order *shopordermodels.Order) *models.OrderVO {
	if order == nil {
		return nil
	}
	items := make([]*models.OrderItem, 0, len(order.Details))
	for _, detail := range order.Details {
		if detail == nil {
			continue
		}
		items = append(items, toShopOrderItem(detail))
	}
	return &models.OrderVO{
		Order: *toShopOrder(order),
		Items: items,
	}
}

func toShopOrder(order *shopordermodels.Order) *models.Order {
	if order == nil {
		return nil
	}
	return &models.Order{
		ID:                    int64(order.ID),
		OrderNo:               order.Tid,
		TotalAmount:           order.Total,
		PayAmount:             order.Total - order.Privilege + order.PostFee,
		FreightAmount:         order.PostFee,
		DiscountAmount:        order.Privilege,
		Status:                orderConstant.OrderStatusToShopStatus(order.Status),
		PayTime:               order.PayTime,
		ReceiverName:          order.ReceiverName,
		ReceiverPhone:         stringUtils.FirstNonEmpty(order.ReceiverMobile, order.ReceiverPhone),
		ReceiverProvince:      stringUtils.FirstNonEmpty(order.ReceiverProvinceName, order.ReceiverProvince),
		ReceiverCity:          stringUtils.FirstNonEmpty(order.ReceiverCityName, order.ReceiverCity),
		ReceiverDistrict:      stringUtils.FirstNonEmpty(order.ReceiverDistrictName, order.ReceiverDistrict),
		ReceiverDetailAddress: order.ReceiverAddress,
		Remark:                stringUtils.FirstNonEmpty(order.SellerMemo, order.BuyerMessage),
		DeptID:                order.DeptID,
		State:                 order.State,
		BaseEntity:            order.BaseEntity,
	}
}

func toShopOrderItem(detail *shopordermodels.OrderDetail) *models.OrderItem {
	if detail == nil {
		return nil
	}
	price := detail.Payment
	if detail.Num > 0 {
		price = detail.Payment / detail.Num
	}
	return &models.OrderItem{
		ID:          int64(detail.ID),
		OrderID:     int64(detail.OrderID),
		OrderNo:     detail.Tid,
		GoodsID:     stringUtils.FirstNonEmpty(detail.EShopGoodsID, fmt.Sprintf("%d", detail.NumIID)),
		SkuID:       stringUtils.FirstNonEmpty(detail.EShopSkuID, fmt.Sprintf("%d", detail.SkuID)),
		GoodsName:   detail.EShopGoodsName,
		SkuName:     detail.EShopSkuName,
		ImageURL:    detail.PicPath,
		Price:       price,
		Quantity:    int32(detail.Num),
		TotalAmount: detail.Payment,
		DeptID:      detail.DeptID,
		State:       detail.State,
		BaseEntity:  detail.BaseEntity,
	}
}

func buildOrderBuyerNick(shopUser *models.User) string {
	if shopUser == nil {
		return ""
	}
	return fmt.Sprintf("shop-user-%s", shopUser.UserID)
}
