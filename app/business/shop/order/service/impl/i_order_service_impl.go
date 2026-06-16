package impl

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"nova-factory-server/app/business/erp/setting/settingdao"
	"nova-factory-server/app/business/shop/order/dao"
	"nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/business/shop/order/service"
	"nova-factory-server/app/constant/commonStatus"
	orderConstant "nova-factory-server/app/constant/order"
	"nova-factory-server/app/datasource/cache"
	objectutil "nova-factory-server/app/utils/json"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/result"
	"nova-factory-server/app/utils/order"
	timeutil "nova-factory-server/app/utils/time"
)

// OrderServiceImpl 提供 ERP 订单的业务实现与同步能力。
type OrderServiceImpl struct {
	orderDao             dao.IOrderDao
	detailDao            dao.IOrderDetailDao
	accountDao           dao.IOrderAccountDao
	integrationConfigDao settingdao.IIntegrationConfigDao
	cache                cache.Cache
	host                 string
}

// NewOrderService 创建 ERP 订单服务。
func NewOrderService(
	orderDao dao.IOrderDao,
	detailDao dao.IOrderDetailDao,
	accountDao dao.IOrderAccountDao,
	integrationConfigDao settingdao.IIntegrationConfigDao,
	cache cache.Cache,
) service.IOrderService {
	host := viper.GetString("host")
	return &OrderServiceImpl{
		orderDao:             orderDao,
		detailDao:            detailDao,
		accountDao:           accountDao,
		cache:                cache,
		host:                 host,
		integrationConfigDao: integrationConfigDao,
	}
}

// Set 新增或修改 ERP 订单。
func (o *OrderServiceImpl) Set(c *gin.Context, req *models.OrderSet) (*models.Order, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	req.Tid = strings.TrimSpace(req.Tid)
	req.ReceiverName = strings.TrimSpace(req.ReceiverName)

	req.ReceiverProvince = strings.TrimSpace(req.ReceiverProvince)
	req.ReceiverProvinceName = strings.TrimSpace(req.ReceiverProvinceName)

	req.ReceiverCity = strings.TrimSpace(req.ReceiverCity)
	req.ReceiverCityName = strings.TrimSpace(req.ReceiverCityName)

	req.ReceiverDistrict = strings.TrimSpace(req.ReceiverDistrict)
	req.ReceiverDistrictName = strings.TrimSpace(req.ReceiverDistrictName)

	req.ReceiverStreet = strings.TrimSpace(req.ReceiverStreet)
	req.ReceiverStreetName = strings.TrimSpace(req.ReceiverStreetName)

	req.ReceiverAddress = strings.TrimSpace(req.ReceiverAddress)
	req.ReceiverMobile = strings.TrimSpace(req.ReceiverMobile)
	req.Status = strings.TrimSpace(req.Status)
	req.OrderType = strings.TrimSpace(req.OrderType)
	if req.Tid == "" {
		if req.ID > 0 {
			return nil, errors.New("tid不能为空")
		}
		req.Tid = order.GenerateOrderNo()
	}
	if req.ReceiverName == "" {
		return nil, errors.New("收货人名称不能为空")
	}
	if req.ReceiverProvince == "" {
		return nil, errors.New("收货省不能为空")
	}
	if req.ReceiverCity == "" {
		return nil, errors.New("收货市不能为空")
	}
	if req.ReceiverDistrict == "" {
		return nil, errors.New("收货区不能为空")
	}
	if req.ReceiverAddress == "" {
		return nil, errors.New("收货地址不能为空")
	}
	if req.ReceiverMobile == "" {
		return nil, errors.New("收货人手机号不能为空")
	}
	if req.Status == "" {
		return nil, errors.New("status不能为空")
	}
	if req.OrderType == "" {
		return nil, errors.New("订单type不能为空")
	}
	if len(req.Details) == 0 {
		return nil, errors.New("details不能为空")
	}
	detailOIDMap := make(map[string]struct{}, len(req.Details))
	for _, detail := range req.Details {
		if detail == nil {
			continue
		}
		detail.OID = strings.TrimSpace(detail.OID)
		detail.EShopGoodsName = strings.TrimSpace(detail.EShopGoodsName)
		if detail.OID == "" {
			return nil, errors.New("订单明细oid不能为空")
		}
		if _, exists := detailOIDMap[detail.OID]; exists {
			return nil, errors.New("订单明细oid重复: " + detail.OID)
		}
		detailOIDMap[detail.OID] = struct{}{}
		if detail.EShopGoodsName == "" {
			return nil, errors.New("订单明细商品名称不能为空")
		}
		if detail.Num <= 0 {
			return nil, errors.New("订单明细数量必须大于0")
		}
	}
	for _, account := range req.Accounts {
		if account == nil {
			continue
		}
		account.FinanceCode = strings.TrimSpace(account.FinanceCode)
		if account.FinanceCode == "" {
			return nil, errors.New("账户编码不能为空")
		}
	}
	return o.orderDao.Set(c, req)
}

// GetByID 查询 ERP 订单详情。
func (o *OrderServiceImpl) GetByID(c *gin.Context, id uint64) (*models.Order, error) {
	if id == 0 {
		return nil, errors.New("id不能为空")
	}
	return o.orderDao.GetByID(c, id)
}

// List 分页查询 ERP 订单。
func (o *OrderServiceImpl) List(c *gin.Context, req *models.OrderQuery) (*models.OrderListData, error) {
	if req == nil {
		req = new(models.OrderQuery)
	}
	req.Tid = strings.TrimSpace(req.Tid)
	req.Status = strings.TrimSpace(req.Status)
	req.ReceiverName = strings.TrimSpace(req.ReceiverName)
	return o.orderDao.List(c, req)
}

// DeleteByIDs 删除 ERP 订单。
func (o *OrderServiceImpl) DeleteByIDs(c *gin.Context, ids []uint64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的订单")
	}
	return o.orderDao.DeleteByIDs(c, ids)
}

// SynchronizeSalesOrders 调用集成客户端接口同步销售订单。
func (o *OrderServiceImpl) SynchronizeSalesOrders(c *gin.Context, req *models.OrderSyncRequest) (result.OrderSyncResponse, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	cfg, err := o.integrationConfigDao.GetEnabled(c)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("未找到管家婆启用配置")
	}
	cfgService, err := cfg.Service()
	if err != nil {
		return nil, err
	}
	if cfgService == nil {
		return nil, errors.New("没有配置集成商")
	}
	syncer := cfgService.OrderSyncer()
	if syncer == nil {
		return nil, errors.New("集成商未实现订单同步能力")
	}

	if err := o.fillOrderSyncRequestFromDB(c, req); err != nil {
		return nil, err
	}
	req.WithConfig(cfg)
	req.WithCache(o.cache)
	return syncer.SyncOrders(c, req)
}

// Sync 同步销售订单事件数据。
//
// 该方法是 ERP 观察者模式下的订单同步入口，职责与 shop 侧的订单同步 service 类似：
// observer 只负责分发 event.OrderEvent，具体的事件数据转换和订单保存放在 service 层。
//
// ERP 订单同步是增量同步：只处理本次事件携带的订单，不删除其他未出现在事件里的订单。
// 方法签名不返回 error 是为了匹配观察者接口；内部错误会完整记录日志，便于排查。
func (o *OrderServiceImpl) Sync(event event.OrderEvent) {
	if event == nil {
		return
	}
	if err := o.syncOrders(models.ToOrder(event)); err != nil {
		zap.L().Error("ERP销售订单同步失败", zap.Error(err))
	}
}

// syncOrders 增量同步订单集合。
//
// 所有订单共用同一个事务：任意订单主表、明细表或账户表写入失败，都会让整个批次回滚。
// 如果后续业务要求“单个订单失败不影响其它订单”，可以把 Transaction 移到循环内部。
func (o *OrderServiceImpl) syncOrders(orders []*models.Order) error {
	if len(orders) == 0 {
		return nil
	}
	if o.orderDao == nil {
		return errors.New("ERP订单主表DAO不能为空")
	}
	if o.detailDao == nil {
		return errors.New("ERP订单明细DAO不能为空")
	}
	if o.accountDao == nil {
		return errors.New("ERP订单账户DAO不能为空")
	}

	return o.orderDao.Transaction(func(tx *gorm.DB) error {
		for _, item := range orders {
			if item == nil {
				continue
			}
			if err := o.syncOne(tx, item); err != nil {
				zap.L().Error("ERP销售订单同步失败", zap.String("tid", item.Tid), zap.Error(err))
				return err
			}
		}
		return nil
	})
}

// syncOne 在外层事务中同步单个 ERP 销售订单。
//
// 处理顺序：
// 1. 校验并标准化 tid；
// 2. 查询 erp_order 是否存在有效记录；
// 3. 准备主表数据，包括 details_json/accounts_json 快照、创建/更新时间、state；
// 4. 不存在则插入主表；
// 5. 已存在则更新主表，删除旧明细和旧账户；
// 6. 插入本次事件携带的新明细和新账户。
//
// 该函数不创建事务，必须使用 syncOrders 传入的 tx。
func (o *OrderServiceImpl) syncOne(tx *gorm.DB, order *models.Order) error {
	if order == nil {
		return nil
	}
	tid := strings.TrimSpace(order.Tid)
	if tid == "" {
		err := errors.New("订单tid不能为空")
		zap.L().Error("ERP销售订单同步失败：订单tid为空", zap.Error(err))
		return err
	}
	order.Tid = tid

	exists, err := o.orderDao.GetByTidTx(tx, tid)
	if err != nil {
		zap.L().Error("ERP销售订单同步失败：查询已存在订单失败", zap.String("tid", tid), zap.Error(err))
		return err
	}

	now := time.Now()
	if err := prepareOrderForSave(order, &now); err != nil {
		zap.L().Error("ERP销售订单同步失败：准备订单数据失败", zap.String("tid", tid), zap.Error(err))
		return err
	}

	var orderID uint64
	if exists == nil {
		if err := o.orderDao.Create(tx, order); err != nil {
			zap.L().Error("ERP销售订单同步失败：创建订单主表失败", zap.String("tid", tid), zap.Error(err))
			return err
		}
		orderID = order.ID
	} else {
		order.ID = exists.ID
		orderID = exists.ID
		if err := o.orderDao.UpdateByID(tx, exists.ID, buildERPOrderUpdateMap(order, exists)); err != nil {
			zap.L().Error("ERP销售订单同步失败：更新订单主表失败", zap.String("tid", tid), zap.Uint64("order_id", exists.ID), zap.Error(err))
			return err
		}
		if err := o.detailDao.DeleteByOrderID(tx, exists.ID); err != nil {
			zap.L().Error("ERP销售订单同步失败：删除旧订单明细失败", zap.String("tid", tid), zap.Uint64("order_id", exists.ID), zap.Error(err))
			return err
		}
		if err := o.accountDao.DeleteByOrderID(tx, exists.ID); err != nil {
			zap.L().Error("ERP销售订单同步失败：删除旧订单账户失败", zap.String("tid", tid), zap.Uint64("order_id", exists.ID), zap.Error(err))
			return err
		}
	}

	if err := o.detailDao.BatchCreateByOrder(tx, orderID, order, &now); err != nil {
		zap.L().Error("ERP销售订单同步失败：创建订单明细失败", zap.String("tid", tid), zap.Uint64("order_id", orderID), zap.Int("details", len(order.Details)), zap.Error(err))
		return err
	}
	if err := o.accountDao.BatchCreateByOrder(tx, orderID, order, &now); err != nil {
		zap.L().Error("ERP销售订单同步失败：创建订单账户失败", zap.String("tid", tid), zap.Uint64("order_id", orderID), zap.Int("accounts", len(order.Accounts)), zap.Error(err))
		return err
	}
	return nil
}

// prepareOrderForSave 准备订单主表写入前的数据。
//
// Details 和 Accounts 会序列化成主表 JSON 快照字段，便于只查主表时保留同步载荷。
// 同时补齐 create_time/update_time/state。
func prepareOrderForSave(order *models.Order, now *time.Time) error {
	detailsJSON, err := objectutil.MarshalJSON(order.Details)
	if err != nil {
		zap.L().Error("ERP销售订单同步失败：订单明细JSON序列化失败", zap.String("tid", order.Tid), zap.Int("details", len(order.Details)), zap.Error(err))
		return fmt.Errorf("订单明细JSON序列化失败: %w", err)
	}
	accountsJSON, err := objectutil.MarshalJSON(order.Accounts)
	if err != nil {
		zap.L().Error("ERP销售订单同步失败：订单账户JSON序列化失败", zap.String("tid", order.Tid), zap.Int("accounts", len(order.Accounts)), zap.Error(err))
		return fmt.Errorf("订单账户JSON序列化失败: %w", err)
	}

	order.DetailsJSON = detailsJSON
	order.AccountsJSON = accountsJSON
	order.CreateTime = timeutil.FirstTime(order.CreateTime, now)
	order.UpdateTime = timeutil.FirstTime(order.UpdateTime, now)
	order.State = commonStatus.NORMAL
	return nil
}

// buildERPOrderUpdateMap 构建 ERP 订单主表更新字段。
//
// 使用 map 是为了允许零值覆盖。status 字段额外经过 shouldUpdateOrderStatus 校验，
// 避免空状态、未知状态、终态覆盖和乱序事件导致的状态回退。
func buildERPOrderUpdateMap(order *models.Order, current *models.Order) map[string]any {
	updates := map[string]any{
		"weight":                 order.Weight,
		"size":                   order.Size,
		"buyer_nick":             order.BuyerNick,
		"buyer_message":          order.BuyerMessage,
		"seller_memo":            order.SellerMemo,
		"total":                  order.Total,
		"privilege":              order.Privilege,
		"post_fee":               order.PostFee,
		"receiver_name":          order.ReceiverName,
		"receiver_province":      order.ReceiverProvince,
		"receiver_province_name": order.ReceiverProvinceName,
		"receiver_city":          order.ReceiverCity,
		"receiver_city_name":     order.ReceiverCityName,
		"receiver_district":      order.ReceiverDistrict,
		"receiver_district_name": order.ReceiverDistrictName,
		"receiver_street":        order.ReceiverStreet,
		"receiver_street_name":   order.ReceiverStreetName,
		"receiver_address":       order.ReceiverAddress,
		"receiver_phone":         order.ReceiverPhone,
		"receiver_mobile":        order.ReceiverMobile,
		"receiver_zip":           order.ReceiverZip,
		"order_type":             order.Type,
		"invoice_name":           order.InvoiceName,
		"seller_flag":            order.SellerFlag,
		"pay_time":               order.PayTime,
		"logist_b_type_code":     order.LogistBTypeCode,
		"logist_bill_code":       order.LogistBillCode,
		"b_type_code":            order.BTypeCode,
		"details_json":           order.DetailsJSON,
		"accounts_json":          order.AccountsJSON,
		"bill_code":              order.BillCode,
		"sync_message":           order.SyncMessage,
		"sync_status":            order.SyncStatus,
		"sync_time":              order.SyncTime,
		"dept_id":                order.DeptID,
		"update_by":              order.UpdateBy,
		"update_time":            order.UpdateTime,
	}

	if shouldUpdateOrderStatus(current.Status, order.Status) {
		updates["status"] = strings.TrimSpace(order.Status)
	} else {
		zap.L().Debug("ERP销售订单同步跳过状态更新",
			zap.String("tid", order.Tid),
			zap.String("current_status", current.Status),
			zap.String("incoming_status", order.Status),
		)
	}
	return updates
}

// shouldUpdateOrderStatus 判断本次同步状态是否允许覆盖数据库当前状态。
func shouldUpdateOrderStatus(current, incoming string) bool {
	current = strings.TrimSpace(current)
	incoming = strings.TrimSpace(incoming)
	if incoming == "" {
		return false
	}

	incomingRank, incomingKnown := orderStatusRank(incoming)
	if !incomingKnown {
		return false
	}
	if current == "" {
		return true
	}
	if isFinalOrderStatus(current) {
		return false
	}
	currentRank, currentKnown := orderStatusRank(current)
	if !currentKnown {
		return true
	}
	return incomingRank >= currentRank
}

// isFinalOrderStatus 判断订单是否已经进入本地终态。
func isFinalOrderStatus(status string) bool {
	switch strings.TrimSpace(status) {
	case orderConstant.ERPStatusTradeSuccess,
		orderConstant.ERPStatusTradeClosed,
		orderConstant.ERPStatusAftersale:
		return true
	default:
		return false
	}
}

// orderStatusRank 返回订单状态推进优先级。
func orderStatusRank(status string) (int, bool) {
	switch strings.TrimSpace(status) {
	case orderConstant.ERPStatusNoPay:
		return 1, true
	case orderConstant.ERPStatusPayed:
		return 2, true
	case orderConstant.ERPStatusPartSend:
		return 3, true
	case orderConstant.ERPStatusSended:
		return 4, true
	case orderConstant.ERPStatusTradeSuccess,
		orderConstant.ERPStatusTradeClosed,
		orderConstant.ERPStatusAftersale:
		return 5, true
	default:
		return 0, false
	}
}
