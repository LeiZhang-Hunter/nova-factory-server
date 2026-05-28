package aidatasetserviceimpl

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
)

type IConfigLoaderServiceImpl struct {
	historyDao gatewaydao.IAIAgentConfigPublishHistoryDao
	agentDao   gatewaydao.IAIAgentDao
}

func NewIConfigLoaderServiceImpl(
	historyDao gatewaydao.IAIAgentConfigPublishHistoryDao,
	agentDao gatewaydao.IAIAgentDao,
) aidatasetservice.IConfigLoaderService {
	return &IConfigLoaderServiceImpl{
		historyDao: historyDao,
		agentDao:   agentDao,
	}
}

func (i *IConfigLoaderServiceImpl) All(ctx context.Context) ([]*gatewaymodels.AIAgentConfigPublishHistory, error) {
	list, err := i.agentDao.GetEnable(ctx)
	if err != nil {
		return nil, err
	}

	var agentIdMap = make(map[int64]string)
	for _, v := range list {
		if v.ConfigVersion == "" {
			continue
		}
		agentIdMap[v.ID] = v.ConfigVersion
	}

	configs, err := i.historyDao.GetConfigsByAgentIdAndVersion(ctx, agentIdMap)
	return configs, err
}

func (i *IConfigLoaderServiceImpl) GetByAgentIdAndVersion(ctx context.Context, agentId uint64, version string) (*gatewaymodels.AIAgentConfigPublishHistory, error) {
	return i.historyDao.GetByAgentIDAndVersion(&gin.Context{}, int64(agentId), version)
}
