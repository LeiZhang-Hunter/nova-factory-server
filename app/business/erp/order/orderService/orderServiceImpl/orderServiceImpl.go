package orderServiceImpl

import (
	"errors"

	"nova-factory-server/app/business/erp/core/integration"
	"nova-factory-server/app/business/erp/order/orderDao"
	"nova-factory-server/app/business/erp/order/orderModels"
	"nova-factory-server/app/business/erp/order/orderService"

	"github.com/gin-gonic/gin"
)

type OrderServiceImpl struct {
	dao orderDao.IOrderDao
}

func NewOrderService(dao orderDao.IOrderDao) orderService.IOrderService {
	return &OrderServiceImpl{dao: dao}
}

func (o *OrderServiceImpl) CheckLoginState(c *gin.Context, req *orderModels.CheckLoginStateReq) (*orderModels.CheckLoginStateResp, error) {
	cfg, err := o.dao.GetEnabledGJPCfg(c)
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
	return &orderModels.CheckLoginStateResp{
		Online:   state.Online,
		Message:  state.Message,
		Type:     state.Type,
		CheckURL: state.CheckURL,
	}, nil
}
