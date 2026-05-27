package gatewayserviceimpl

import (
	"errors"
	"nova-factory-server/app/constant/aiagent"
	"strconv"
	"strings"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	v1 "nova-factory-server/app/utils/grpc/confighotload/v1"
	"nova-factory-server/app/utils/uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIAgentConfigPublishHistoryServiceImpl 提供智能体配置发布历史的业务实现。
type AIAgentConfigPublishHistoryServiceImpl struct {
	dao              gatewaydao.IAIAgentConfigPublishHistoryDao
	agentDao         gatewaydao.IAIAgentDao
	orchestrationDao gatewaydao.IAIAgentOrchestrationDao
	gatewaySvc       gatewayservice.IAIGatewayService
	configPusher     gatewayservice.IAIAgentConfigPublisher
	db               *gorm.DB
}

// NewAIAgentConfigPublishHistoryService 创建智能体配置发布历史服务。
func NewAIAgentConfigPublishHistoryService(
	dao gatewaydao.IAIAgentConfigPublishHistoryDao,
	agentDao gatewaydao.IAIAgentDao,
	gatewaySvc gatewayservice.IAIGatewayService,
	configPusher gatewayservice.IAIAgentConfigPublisher,
	orchestrationDao gatewaydao.IAIAgentOrchestrationDao,
	db *gorm.DB,
) gatewayservice.IAIAgentConfigPublishHistoryService {
	return &AIAgentConfigPublishHistoryServiceImpl{
		dao:              dao,
		agentDao:         agentDao,
		gatewaySvc:       gatewaySvc,
		configPusher:     configPusher,
		orchestrationDao: orchestrationDao,
		db:               db,
	}
}

// Set 保存智能体配置发布历史。
func (a *AIAgentConfigPublishHistoryServiceImpl) Set(c *gin.Context, req *gatewaymodels.AIAgentConfigPublishHistoryUpsert) (*gatewaymodels.AIAgentConfigPublishHistory, error) {
	if req.Action != string(aiagent.ConfigPublishType) && req.Action != string(aiagent.ConfigRemoveType) {
		return nil, errors.New("下发类型错误")
	}
	if err := a.validateUpsert(c, req); err != nil {
		return nil, err
	}
	exists, err := a.dao.GetByAgentIDAndVersion(c, req.AgentID, req.Version)
	if err != nil {
		return nil, err
	}
	isUpdate := req.ID > 0
	if exists != nil {
		req.ID = exists.ID
		isUpdate = true
	}
	if isUpdate {
		current, err := a.dao.GetByID(c, req.ID)
		if err != nil {
			return nil, err
		}
		if current == nil {
			return nil, errors.New("配置发布历史不存在")
		}
	}
	if a.db == nil {
		return nil, errors.New("数据库连接不存在")
	}
	var item *gatewaymodels.AIAgentConfigPublishHistory
	err = a.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// 读取编排内容
		orchestrationConfig, err := a.orchestrationDao.GetByAgentID(c, req.AgentID)
		if err != nil {
			return err
		}

		if orchestrationConfig == nil {
			return nil
		}

		req.ConfigSnapshot = orchestrationConfig.Config
		req.ConfigMd5 = orchestrationConfig.ConfigMd5

		if isUpdate {
			item, err = a.dao.UpdateWithTx(c, tx, req)
		} else {
			item, err = a.dao.CreateWithTx(c, tx, req)
		}
		if err != nil {
			return err
		}
		err = a.agentDao.UpdateConfigVersion(c, req.AgentID, req.Version)
		if err != nil {
			return err
		}
		return a.publishConfig(c, req.AgentID, aiagent.ConfigUpdate(req.Action))
	})
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Info 查询智能体配置发布历史详情。
func (a *AIAgentConfigPublishHistoryServiceImpl) Info(c *gin.Context, id int64) (*gatewaymodels.AIAgentConfigPublishHistory, error) {
	if id == 0 {
		return nil, errors.New("id不能为空")
	}
	return a.dao.GetByID(c, id)
}

// List 查询智能体配置发布历史列表。
func (a *AIAgentConfigPublishHistoryServiceImpl) List(c *gin.Context, req *gatewaymodels.AIAgentConfigPublishHistoryQuery) (*gatewaymodels.AIAgentConfigPublishHistoryListData, error) {
	if req == nil {
		req = new(gatewaymodels.AIAgentConfigPublishHistoryQuery)
	}
	req.Version = strings.TrimSpace(req.Version)
	return a.dao.List(c, req)
}

// Remove 删除智能体配置发布历史。
func (a *AIAgentConfigPublishHistoryServiceImpl) Remove(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的发布历史")
	}
	agentIDs := make(map[int64]struct{})
	for _, id := range ids {
		if id == 0 {
			return errors.New("id不能为空")
		}
		item, err := a.dao.GetByID(c, id)
		if err != nil {
			return err
		}
		if item != nil && item.AgentID > 0 {
			agentIDs[item.AgentID] = struct{}{}
		}
	}
	if a.db == nil {
		return errors.New("数据库连接不存在")
	}
	return a.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := a.dao.DeleteByIDsWithTx(c, tx, ids); err != nil {
			return err
		}
		for agentID := range agentIDs {
			if err := a.publishConfig(c, agentID, aiagent.ConfigRemoveType); err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *AIAgentConfigPublishHistoryServiceImpl) validateUpsert(c *gin.Context, req *gatewaymodels.AIAgentConfigPublishHistoryUpsert) error {
	if req == nil {
		return errors.New("参数不能为空")
	}

	req.Version = strings.TrimSpace(req.Version)
	req.ConfigSnapshot = strings.TrimSpace(req.ConfigSnapshot)
	req.PublishDescription = strings.TrimSpace(req.PublishDescription)
	if req.AgentID == 0 {
		return errors.New("agentId不能为空")
	}
	if req.Version == "" {
		return errors.New("version不能为空")
	}
	if req.ConfigSnapshot == "" {
		return errors.New("configSnapshot不能为空")
	}

	agent, err := a.agentDao.GetByID(c, req.AgentID)
	if err != nil {
		return err
	}
	if agent == nil {
		return errors.New("智能体不存在")
	}
	return nil
}

func (a *AIAgentConfigPublishHistoryServiceImpl) publishConfig(c *gin.Context, agentID int64, updateType aiagent.ConfigUpdate) error {
	gatewayInfo, err := a.getPublishGateway(c)
	if err != nil {
		return err
	}

	reply, err := a.configPusher.BroadcastByGrpcClient(c, &v1.AgentBroadcastRequest{
		RequestId: uuid.MakeUuid(),
		GatewayId: gatewayInfo.ID,
		AgentId:   strconv.FormatInt(agentID, 10),
		Action:    string(updateType),
	})
	if err != nil {
		return err
	}
	if reply == nil {
		return errors.New("发布配置失败")
	}
	if reply.GetCode() != 0 {
		return errors.New(reply.GetMsg())
	}
	return nil
}

func (a *AIAgentConfigPublishHistoryServiceImpl) getPublishGateway(c *gin.Context) (*gatewaymodels.AIGateway, error) {
	enabled := true
	list, err := a.gatewaySvc.List(c, &gatewaymodels.AIGatewayQuery{
		Enabled: &enabled,
		Page:    1,
		Size:    1,
	})
	if err != nil {
		return nil, err
	}
	if list == nil || len(list.Rows) == 0 || list.Rows[0] == nil {
		return nil, errors.New("未找到可用网关，无法发布配置")
	}
	return list.Rows[0], nil
}

func (a *AIAgentConfigPublishHistoryServiceImpl) GetByIds() {

}
