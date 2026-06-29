package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	configDao "nova-factory-server/app/business/shop/config/dao"
	"nova-factory-server/app/business/shop/logistics/client/api"
	"nova-factory-server/app/business/shop/logistics/dao"
	"nova-factory-server/app/business/shop/logistics/models"
	"nova-factory-server/app/business/shop/logistics/service"
	redisk "nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"

	"github.com/gin-gonic/gin"
)

const (
	// TrackingCacheTTL 物流轨迹短期缓存时间（10分钟）
	TrackingCacheTTL = 10 * time.Minute
)

// stateDescMap 统一状态码 → 中文描述
var stateDescMap = map[string]string{
	"0": "暂无轨迹",
	"1": "已揽收",
	"2": "在途中",
	"3": "已签收",
	"4": "问题件",
}

// TrackingServiceImpl 物流轨迹查询服务实现
type TrackingServiceImpl struct {
	dao                    dao.ITrackingDao
	cache                  cache.Cache
	logisticsDao           configDao.ILogisticsCompanyDao
	shopLogisticsConfigDao configDao.IShopLogisticsConfigDao
}

// NewTrackingService 创建物流轨迹查询服务
func NewTrackingService(
	cache cache.Cache,
	dao dao.ITrackingDao,
	logisticsDao configDao.ILogisticsCompanyDao,
	shopLogisticsConfigDao configDao.IShopLogisticsConfigDao,
) service.ITrackingService {
	return &TrackingServiceImpl{
		dao:                    dao,
		cache:                  cache,
		logisticsDao:           logisticsDao,
		shopLogisticsConfigDao: shopLogisticsConfigDao,
	}
}

// Query 即时查询物流轨迹（缓存优先策略：DB已签收记录 → Redis短期缓存 → 第三方API）
func (s *TrackingServiceImpl) Query(c *gin.Context, outsid, companyCode string) (*models.TrackingQueryResponse, error) {
	// 1. 查 DB 已签收记录（永久缓存，直接返回）
	signedResp, err := s.dao.GetSignedRecord(c, outsid, companyCode)
	if err != nil {
		return nil, fmt.Errorf("查询已签收记录失败: %w", err)
	}
	if signedResp != nil {
		return signedResp, nil
	}

	// 2. 查 Redis 短期缓存
	cacheKey := redisk.MakeLogisticsTrackingCacheKey(outsid, companyCode)
	if cached, err := s.cache.Get(c, cacheKey); err == nil && cached != "" {
		var resp models.TrackingQueryResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			resp.FromCache = true
			return &resp, nil
		}
	}

	shopLogisticsCfg, err := s.shopLogisticsConfigDao.GetEnabled(c)
	if err != nil {
		return nil, err
	}
	if shopLogisticsCfg == nil {
		return nil, errors.New("没有快递配置")
	}

	shopLogisticsClient, err := shopLogisticsCfg.Service()
	if err != nil {
		return nil, err
	}
	if shopLogisticsClient == nil {
		return nil, errors.New("没有快递配置")
	}

	// 3. 调用快递鸟客户端即时查询
	result, err := shopLogisticsClient.Query(companyCode, outsid)
	if err != nil {
		return nil, fmt.Errorf("物流查询失败: %w", err)
	}
	if !result.Success() {
		reason := result.Reason()
		if reason == "" {
			reason = "未知错误"
		}
		// 失败时也缓存结果，避免频繁无效调用（缓存 30 分钟）
		failJSON, _ := json.Marshal(map[string]string{"error": reason})
		s.cache.Set(c, cacheKey, string(failJSON), 30*time.Minute)
		return nil, fmt.Errorf("物流查询失败: %s", reason)
	}

	// 5. 统一模型 → 响应 DTO
	resp := s.convertResult(result)
	resp.FromCache = false

	// 6. 缓存策略
	respJSON, _ := json.Marshal(resp)
	respStr := string(respJSON)
	if resp.IsSigned {
		// 已签收：写入 DB 永久缓存
		_ = s.dao.SaveSignedRecord(c, &models.TrackingRecordSet{
			Outsid:      outsid,
			CompanyCode: companyCode,
			TraceJSON:   respStr,
			SignedTime:  s.extractSignedTime(result),
			OriginInfo:  s.extractOriginInfo(result),
			DestInfo:    result.Location(),
		})
	} else {
		// 未签收：写入 Redis 短期缓存（10分钟）
		s.cache.Set(c, cacheKey, respStr, TrackingCacheTTL)
	}

	return resp, nil
}

func (s *TrackingServiceImpl) convertResult(result api.ExpressQueryResult) *models.TrackingQueryResponse {
	traceList := result.Traces()
	traces := make([]*models.TrackingTraceNode, 0, len(traceList))
	for _, t := range traceList {
		traces = append(traces, &models.TrackingTraceNode{
			AcceptTime:    t.AcceptTime(),
			AcceptStation: t.AcceptStation(),
			Location:      t.Location(),
			Action:        t.Action(),
		})
	}

	state := result.State()
	return &models.TrackingQueryResponse{
		Outsid:      result.LogisticCode(),
		CompanyCode: result.ShipperCode(),
		State:       state,
		StateDesc:   stateDescMap[state],
		IsSigned:    state == SignedState,
		Traces:      traces,
		Location:    result.Location(),
		StateName:   result.GetStateName(),
	}
}

const SignedState = "3"

// extractSignedTime 提取签收时间（取最后一条轨迹的时间）
func (s *TrackingServiceImpl) extractSignedTime(result api.ExpressQueryResult) string {
	traces := result.Traces()
	if result.State() != SignedState || len(traces) == 0 {
		return ""
	}
	return traces[len(traces)-1].AcceptTime()
}

// extractOriginInfo 提取始发地（取第一条轨迹所在城市）
func (s *TrackingServiceImpl) extractOriginInfo(result api.ExpressQueryResult) string {
	traces := result.Traces()
	if len(traces) > 0 {
		return traces[0].Location()
	}
	return ""
}
