package impl

//
//const (
//	// TrackingCacheTTL 物流轨迹短期缓存时间（10分钟）
//	TrackingCacheTTL = 10 * time.Minute
//)
//
//// stateDescMap 统一状态码 → 中文描述
//var stateDescMap = map[string]string{
//	"0": "暂无轨迹",
//	"1": "已揽收",
//	"2": "在途中",
//	"3": "已签收",
//	"4": "问题件",
//}
//
//// TrackingServiceImpl 物流轨迹查询服务实现
//type TrackingServiceImpl struct {
//	router       *registry.ExpressRouter
//	dao          dao.ITrackingDao
//	cache        cache.Cache
//	logisticsDao configDao.ILogisticsCompanyDao
//}
//
//// NewTrackingService 创建物流轨迹查询服务，router 通过单例获取
//func NewTrackingService(
//	cache cache.Cache,
//	dao dao.ITrackingDao,
//	logisticsDao configDao.ILogisticsCompanyDao,
//) service.ITrackingService {
//	return &TrackingServiceImpl{
//		router:       registry.GetRouter(),
//		dao:          dao,
//		cache:        cache,
//		logisticsDao: logisticsDao,
//	}
//}
//
//// Query 即时查询物流轨迹（缓存优先策略：DB已签收记录 → Redis短期缓存 → 第三方API）
//func (s *TrackingServiceImpl) Query(c *gin.Context, outsid, companyCode string) (*models.TrackingQueryResponse, error) {
//	// 1. 查 DB 已签收记录（永久缓存，直接返回）
//	signedResp, err := s.dao.GetSignedRecord(c, outsid, companyCode)
//	if err != nil {
//		return nil, fmt.Errorf("查询已签收记录失败: %w", err)
//	}
//	if signedResp != nil {
//		return signedResp, nil
//	}
//
//	// 2. 查 Redis 短期缓存
//	cacheKey := redisk.MakeLogisticsTrackingCacheKey(outsid, companyCode)
//	if cached, err := s.cache.Get(c, cacheKey); err == nil && cached != "" {
//		var resp models.TrackingQueryResponse
//		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
//			resp.FromCache = true
//			return &resp, nil
//		}
//	}
//
//	// 3. 确定第三方渠道，默认快递鸟
//	channel := kdniao.Channel
//	if company, err := s.logisticsDao.GetByCode(c, companyCode); err == nil && company != nil && company.ThirdParty != "" {
//		channel = company.ThirdParty
//	}
//
//	// 4. 通过路由器调用对应第三方 API
//	result, err := s.router.Query(channel, companyCode, outsid)
//	if err != nil {
//		return nil, fmt.Errorf("物流查询失败: %w", err)
//	}
//	if !result.Success {
//		reason := result.ErrorMsg
//		if reason == "" {
//			reason = "未知错误"
//		}
//		// 失败时也缓存结果，避免频繁无效调用（缓存 30 分钟）
//		failJSON, _ := json.Marshal(map[string]string{"error": reason})
//		s.cache.Set(c, cacheKey, string(failJSON), 30*time.Minute)
//		return nil, fmt.Errorf("物流查询失败: %s", reason)
//	}
//
//	// 5. 统一模型 → 响应 DTO
//	resp := s.convertResult(outsid, companyCode, result)
//	resp.FromCache = false
//
//	// 6. 缓存策略
//	respJSON, _ := json.Marshal(resp)
//	respStr := string(respJSON)
//	if resp.IsSigned {
//		// 已签收：写入 DB 永久缓存
//		_ = s.dao.SaveSignedRecord(c, &models.TrackingRecordSet{
//			Outsid:      outsid,
//			CompanyCode: companyCode,
//			TraceJSON:   respStr,
//			SignedTime:  s.extractSignedTime(result),
//			OriginInfo:  s.extractOriginInfo(result),
//			DestInfo:    result.Location,
//		})
//	} else {
//		// 未签收：写入 Redis 短期缓存（10分钟）
//		s.cache.Set(c, cacheKey, respStr, TrackingCacheTTL)
//	}
//
//	return resp, nil
//}
//
//func (s *TrackingServiceImpl) convertResult(outsid, companyCode string, result *api.ExpressTrackResult) *models.TrackingQueryResponse {
//	traces := make([]*models.TrackingTraceNode, 0, len(result.Traces))
//	for _, t := range result.Traces {
//		traces = append(traces, &models.TrackingTraceNode{
//			AcceptTime:    t.AcceptTime,
//			AcceptStation: t.AcceptStation,
//			Location:      t.Location,
//			Action:        t.Action,
//		})
//	}
//
//	return &models.TrackingQueryResponse{
//		Outsid:      outsid,
//		CompanyCode: companyCode,
//		State:       result.State,
//		StateDesc:   stateDescMap[result.State],
//		IsSigned:    result.IsSigned,
//		Traces:      traces,
//	}
//}
//
//// extractSignedTime 提取签收时间（取最后一条轨迹的时间）
//func (s *TrackingServiceImpl) extractSignedTime(result *api.ExpressTrackResult) string {
//	if !result.IsSigned || len(result.Traces) == 0 {
//		return ""
//	}
//	return result.Traces[len(result.Traces)-1].AcceptTime
//}
//
//// extractOriginInfo 提取始发地（取第一条轨迹所在城市）
//func (s *TrackingServiceImpl) extractOriginInfo(result *api.ExpressTrackResult) string {
//	if len(result.Traces) > 0 {
//		return result.Traces[0].Location
//	}
//	return ""
//}
