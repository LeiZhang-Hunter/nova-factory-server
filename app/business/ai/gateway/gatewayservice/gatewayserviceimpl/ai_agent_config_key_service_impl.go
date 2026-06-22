package gatewayserviceimpl

import (
	"errors"
	"fmt"
	uuid2 "nova-factory-server/app/utils/uuid"
	"regexp"
	"strings"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AgentConfigKeyServiceImpl API Key Service 实现。
type AgentConfigKeyServiceImpl struct {
	dao gatewaydao.IAgentConfigKeyDao
}

// NewAgentConfigKeyService 创建 AgentConfigKeyServiceImpl。
func NewAgentConfigKeyService(dao gatewaydao.IAgentConfigKeyDao) gatewayservice.IAgentConfigKeyService {
	return &AgentConfigKeyServiceImpl{dao: dao}
}

func (a *AgentConfigKeyServiceImpl) Create(c *gin.Context, req *gatewaymodels.AgentConfigKeyUpsert) (*gatewaymodels.AgentConfigKey, error) {
	if err := a.validateUpsert(req); err != nil {
		return nil, err
	}
	existing, err := a.dao.GetByKey(c, req.Key)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("该 API Key 已存在")
	}
	return a.dao.Create(c, req)
}

func (a *AgentConfigKeyServiceImpl) Update(c *gin.Context, req *gatewaymodels.AgentConfigKeyUpsert) (*gatewaymodels.AgentConfigKey, error) {
	if req.ID == 0 {
		return nil, errors.New("id不能为空")
	}
	if err := a.validateUpsert(req); err != nil {
		return nil, err
	}
	return a.dao.Update(c, req)
}

func (a *AgentConfigKeyServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return a.dao.DeleteByIDs(c, ids)
}

func (a *AgentConfigKeyServiceImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.AgentConfigKey, error) {
	return a.dao.GetByID(c, id)
}

func (a *AgentConfigKeyServiceImpl) List(c *gin.Context, req *gatewaymodels.AgentConfigKeyQuery) (*gatewaymodels.AgentConfigKeyListData, error) {
	if req == nil {
		req = new(gatewaymodels.AgentConfigKeyQuery)
	}
	return a.dao.List(c, req)
}

func (a *AgentConfigKeyServiceImpl) Generate() string {
	uuidStr := strings.ReplaceAll(uuid.NewString(), "-", "")
	key := fmt.Sprintf("sk-%s", uuid2.MakeMd5([]byte(uuidStr)))
	return key
}

func (a *AgentConfigKeyServiceImpl) validateUpsert(req *gatewaymodels.AgentConfigKeyUpsert) error {
	if strings.TrimSpace(req.Key) == "" {
		return errors.New("密钥不能为空")
	}
	if !regexp.MustCompile(`^sk-[0-9a-f]{32}$`).MatchString(strings.TrimSpace(req.Key)) {
		return errors.New("密钥格式不正确，需为 sk- 前缀加 32 位小写 hex")
	}
	return nil
}
