package alertServiceImpl

import (
	"errors"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"nova-factory-server/app/business/ai/agent/aiDataSetService"
	"nova-factory-server/app/business/iot/alert/alertService"
	"sync"

	"github.com/gin-gonic/gin"
)

type RunnerServiceImpl struct {
	chatService  aiDataSetService.IChartService
	agentSession map[string]string
	mtx          sync.RWMutex
}

func NewRunnerServiceImpl(chatService aiDataSetService.IChartService) alertService.RunnerService {
	return &RunnerServiceImpl{
		chatService:  chatService,
		agentSession: make(map[string]string),
	}
}

func (r *RunnerServiceImpl) Load(ctx *gin.Context, agentId string) (string, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	sessionId, ok := r.agentSession[agentId]
	if ok {
		return sessionId, nil
	}
	agent, err := r.chatService.AgentSessionCreate(ctx, &aiDataSetModels.SessionAgentCreate{
		AgentId: agentId,
	})
	if err != nil {
		return "", err
	}

	if agent == nil {
		return "", errors.New("agent not found")
	}
	r.agentSession[agent.Data.AgentId] = agent.Data.Id
	return agent.Data.Id, nil
}
