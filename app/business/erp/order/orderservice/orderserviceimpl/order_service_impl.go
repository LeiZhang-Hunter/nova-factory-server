package orderserviceimpl

import (
	"errors"
	"nova-factory-server/app/business/erp/setting/settingdao"
	"strings"

	"nova-factory-server/app/business/erp/core/integration"
	"nova-factory-server/app/business/erp/core/integration/grasp"
	"nova-factory-server/app/business/erp/order/orderdao"
	"nova-factory-server/app/business/erp/order/ordermodels"
	"nova-factory-server/app/business/erp/order/orderservice"
	"nova-factory-server/app/datasource/cache"

	"github.com/gin-gonic/gin"
)

// OrderServiceImpl 提供 ERP 订单的业务实现与同步能力。
type OrderServiceImpl struct {
	dao                  orderdao.IOrderDao
	integrationConfigDao settingdao.IIntegrationConfigDao
	cache                cache.Cache
}

// NewOrderService 创建 ERP 订单服务。
func NewOrderService(dao orderdao.IOrderDao, integrationConfigDao settingdao.IIntegrationConfigDao, cache cache.Cache) orderservice.IOrderService {
	return &OrderServiceImpl{
		dao:                  dao,
		cache:                cache,
		integrationConfigDao: integrationConfigDao,
	}
}

// Set 新增或修改 ERP 订单。
func (o *OrderServiceImpl) Set(c *gin.Context, req *ordermodels.OrderSet) (*ordermodels.Order, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	req.Tid = strings.TrimSpace(req.Tid)
	req.ReceiverName = strings.TrimSpace(req.ReceiverName)
	req.ReceiverState = strings.TrimSpace(req.ReceiverState)
	req.ReceiverCity = strings.TrimSpace(req.ReceiverCity)
	req.ReceiverDistrict = strings.TrimSpace(req.ReceiverDistrict)
	req.ReceiverAddress = strings.TrimSpace(req.ReceiverAddress)
	req.ReceiverMobile = strings.TrimSpace(req.ReceiverMobile)
	req.Status = strings.TrimSpace(req.Status)
	req.Type = strings.TrimSpace(req.Type)
	if req.Tid == "" {
		return nil, errors.New("tid不能为空")
	}
	if req.ReceiverName == "" || req.ReceiverState == "" || req.ReceiverCity == "" || req.ReceiverDistrict == "" || req.ReceiverAddress == "" || req.ReceiverMobile == "" {
		return nil, errors.New("收货信息不能为空")
	}
	if req.Status == "" {
		return nil, errors.New("status不能为空")
	}
	if req.Type == "" {
		return nil, errors.New("type不能为空")
	}
	if len(req.Details) == 0 {
		return nil, errors.New("details不能为空")
	}
	for _, detail := range req.Details {
		if detail == nil {
			continue
		}
		detail.OID = strings.TrimSpace(detail.OID)
		detail.EShopGoodsName = strings.TrimSpace(detail.EShopGoodsName)
		if detail.OID == "" {
			return nil, errors.New("订单明细oid不能为空")
		}
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
	return o.dao.Set(c, req)
}

// GetByID 查询 ERP 订单详情。
func (o *OrderServiceImpl) GetByID(c *gin.Context, id uint64) (*ordermodels.Order, error) {
	if id == 0 {
		return nil, errors.New("id不能为空")
	}
	return o.dao.GetByID(c, id)
}

// List 分页查询 ERP 订单。
func (o *OrderServiceImpl) List(c *gin.Context, req *ordermodels.OrderQuery) (*ordermodels.OrderListData, error) {
	if req == nil {
		req = new(ordermodels.OrderQuery)
	}
	req.Tid = strings.TrimSpace(req.Tid)
	req.Status = strings.TrimSpace(req.Status)
	req.ReceiverName = strings.TrimSpace(req.ReceiverName)
	return o.dao.List(c, req)
}

// DeleteByIDs 删除 ERP 订单。
func (o *OrderServiceImpl) DeleteByIDs(c *gin.Context, ids []uint64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的订单")
	}
	return o.dao.DeleteByIDs(c, ids)
}

func (o *OrderServiceImpl) CheckLoginState(c *gin.Context, req *ordermodels.CheckLoginStateReq) (*ordermodels.CheckLoginStateResp, error) {
	// CheckLoginState 检查当前 ERP 集成客户端的登录状态。
	cfg, err := o.integrationConfigDao.GetEnabled(c)
	if err != nil {
		return nil, err
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("未找到管家婆启用配置")
	}
	client, err := integration.CreateByType(cfg.Type)
	if err != nil {
		return nil, err
	}
	state, err := client.CheckLoginState(c, cfg, req.CheckURL, "")
	if err != nil {
		return nil, err
	}
	return &ordermodels.CheckLoginStateResp{
		Online:   state.Online,
		Message:  state.Message,
		Type:     state.Type,
		CheckURL: state.CheckURL,
	}, nil
}

// SynchronizeSalesOrders 调用管家婆接口同步销售订单。
func (o *OrderServiceImpl) SynchronizeSalesOrders(c *gin.Context, req *grasp.OrderSyncRequest) (*grasp.OrderSyncResponse, error) {
	cfg, err := o.integrationConfigDao.GetEnabled(c)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("未找到管家婆启用配置")
	}
	return grasp.New().SynchronizeOrders(c, cfg, req, o.cache)
}
