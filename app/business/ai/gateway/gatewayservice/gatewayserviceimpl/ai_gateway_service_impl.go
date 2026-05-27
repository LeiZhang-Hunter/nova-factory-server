package gatewayserviceimpl

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	rediskey "nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AIGatewayServiceImpl struct {
	dao   gatewaydao.IAIGatewayDao
	cache cache.Cache
}

func NewAIGatewayService(dao gatewaydao.IAIGatewayDao, cache cache.Cache) gatewayservice.IAIGatewayService {
	return &AIGatewayServiceImpl{dao: dao, cache: cache}
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
	if req == nil {
		req = new(gatewaymodels.AIGatewayQuery)
	}
	data, err := a.dao.List(c, req)
	if err != nil || data == nil || len(data.Rows) == 0 {
		return data, err
	}
	a.fillGatewayActive(c, data.Rows)
	return data, nil
}

// RefreshAlive 刷新网关在线标记。
func (a *AIGatewayServiceImpl) RefreshAlive(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("id不能为空")
	}
	a.cache.Set(ctx, rediskey.MakeAIGatewayAliveCacheKey(id), strconv.FormatInt(time.Now().Unix(), 10), 2*time.Minute)
	return nil
}

func (a *AIGatewayServiceImpl) fillGatewayActive(ctx context.Context, rows []*gatewaymodels.AIGateway) {
	keys := make([]string, 0, len(rows))
	gateways := make([]*gatewaymodels.AIGateway, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		row.Active = 0
		keys = append(keys, rediskey.MakeAIGatewayAliveCacheKey(row.ID))
		gateways = append(gateways, row)
	}
	if len(keys) == 0 {
		return
	}

	// 批量读取 Redis 在线标记，避免列表场景逐条查询带来的额外开销。
	values, err := a.cache.MGet(ctx, keys).Result()
	if err != nil {
		zap.L().Warn("batch get gateway alive failed", zap.Error(err))
		return
	}
	for i, value := range values {
		if value == nil {
			continue
		}
		gateways[i].Active = 1
	}
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
