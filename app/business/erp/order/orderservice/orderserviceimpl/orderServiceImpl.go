package orderserviceimpl

import (
	"errors"
	"nova-factory-server/app/business/erp/setting/settingdao"

	"nova-factory-server/app/business/erp/core/integration"
	"nova-factory-server/app/business/erp/core/integration/grasp"
	"nova-factory-server/app/business/erp/order/orderdao"
	"nova-factory-server/app/business/erp/order/ordermodels"
	"nova-factory-server/app/business/erp/order/orderservice"
	"nova-factory-server/app/datasource/cache"

	"github.com/gin-gonic/gin"
)

type OrderServiceImpl struct {
	dao                  orderdao.IOrderDao
	integrationConfigDao settingdao.IIntegrationConfigDao
	cache                cache.Cache
}

func NewOrderService(dao orderdao.IOrderDao, integrationConfigDao settingdao.IIntegrationConfigDao, cache cache.Cache) orderservice.IOrderService {
	return &OrderServiceImpl{
		dao:                  dao,
		cache:                cache,
		integrationConfigDao: integrationConfigDao,
	}
}

func (o *OrderServiceImpl) CheckLoginState(c *gin.Context, req *ordermodels.CheckLoginStateReq) (*ordermodels.CheckLoginStateResp, error) {
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
