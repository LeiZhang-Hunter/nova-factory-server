package saleserviceimpl

import (
	"errors"
	"nova-factory-server/app/business/erp/core/integration"
	"nova-factory-server/app/business/erp/core/integration/api"
	"nova-factory-server/app/business/erp/setting/settingdao"
	"nova-factory-server/app/utils/order"
	"strings"

	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"
	"nova-factory-server/app/datasource/cache"

	"github.com/gin-gonic/gin"
)

// OrderServiceImpl 提供 ERP 订单的业务实现与同步能力。
type OrderServiceImpl struct {
	dao                  saledao.IOrderDao
	integrationConfigDao settingdao.IIntegrationConfigDao
	cache                cache.Cache
}

// NewOrderService 创建 ERP 订单服务。
func NewOrderService(dao saledao.IOrderDao, integrationConfigDao settingdao.IIntegrationConfigDao, cache cache.Cache) saleservice.IOrderService {
	return &OrderServiceImpl{
		dao:                  dao,
		cache:                cache,
		integrationConfigDao: integrationConfigDao,
	}
}

// Set 新增或修改 ERP 订单。
func (o *OrderServiceImpl) Set(c *gin.Context, req *salemodels.OrderSet) (*salemodels.Order, error) {
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
	return o.dao.Set(c, req)
}

// GetByID 查询 ERP 订单详情。
func (o *OrderServiceImpl) GetByID(c *gin.Context, id uint64) (*salemodels.Order, error) {
	if id == 0 {
		return nil, errors.New("id不能为空")
	}
	return o.dao.GetByID(c, id)
}

// List 分页查询 ERP 订单。
func (o *OrderServiceImpl) List(c *gin.Context, req *salemodels.OrderQuery) (*salemodels.OrderListData, error) {
	if req == nil {
		req = new(salemodels.OrderQuery)
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

func (o *OrderServiceImpl) CheckLoginState(c *gin.Context, req *salemodels.CheckLoginStateReq) (*salemodels.CheckLoginStateResp, error) {
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
	return &salemodels.CheckLoginStateResp{
		Online:   state.Online,
		Message:  state.Message,
		Type:     state.Type,
		CheckURL: state.CheckURL,
	}, nil
}

// SynchronizeSalesOrders 调用管家婆接口同步销售订单。
func (o *OrderServiceImpl) SynchronizeSalesOrders(c *gin.Context, req *api.OrderSyncRequest) (*api.OrderSyncResponse, error) {
	cfg, err := o.integrationConfigDao.GetEnabled(c)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("未找到管家婆启用配置")
	}
	client, err := integration.CreateByType(cfg.Type)
	if err != nil {
		return nil, err
	}
	return client.SynchronizeOrders(c, cfg, req, o.cache)
}
