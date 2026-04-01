package gatewayserviceimpl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"

	"github.com/gin-gonic/gin"
)

type AIGatewayServiceImpl struct {
	dao gatewaydao.IAIGatewayDao
}

func NewAIGatewayService(dao gatewaydao.IAIGatewayDao) gatewayservice.IAIGatewayService {
	return &AIGatewayServiceImpl{dao: dao}
}

func (a *AIGatewayServiceImpl) Create(c *gin.Context, req *gatewaymodels.AIGatewayUpsert) (*gatewaymodels.AIGateway, error) {
	if err := a.validateUpsert(req); err != nil {
		return nil, err
	}
	return a.dao.Create(c, req)
}

func (a *AIGatewayServiceImpl) Update(c *gin.Context, req *gatewaymodels.AIGatewayUpsert) (*gatewaymodels.AIGateway, error) {
	if req.ID == 0 {
		return nil, errors.New("id不能为空")
	}
	if err := a.validateUpsert(req); err != nil {
		return nil, err
	}
	return a.dao.Update(c, req)
}

func (a *AIGatewayServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return a.dao.DeleteByIDs(c, ids)
}

func (a *AIGatewayServiceImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.AIGateway, error) {
	return a.dao.GetByID(c, id)
}

func (a *AIGatewayServiceImpl) List(c *gin.Context, req *gatewaymodels.AIGatewayQuery) (*gatewaymodels.AIGatewayListData, error) {
	return a.dao.List(c, req)
}

func (a *AIGatewayServiceImpl) validateUpsert(req *gatewaymodels.AIGatewayUpsert) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("网关名称不能为空")
	}
	if strings.TrimSpace(req.BaseURL) == "" {
		return errors.New("API服务器地址不能为空")
	}
	if strings.TrimSpace(req.APIKey) == "" {
		return errors.New("API Key不能为空")
	}

	return nil
}
