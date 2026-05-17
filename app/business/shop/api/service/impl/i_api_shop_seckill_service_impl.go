package impl

import (
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	myTime "nova-factory-server/app/utils/time"
	"time"

	"github.com/gin-gonic/gin"
)

// IApiShopSeckillServiceImpl 秒杀服务实现
type IApiShopSeckillServiceImpl struct {
	seckillDao       dao.IApiShopSeckillDao
	seckillConfigDao dao.IApiShopSeckillConfigDao
}

// NewIApiShopSeckillServiceImpl 创建秒杀服务
func NewIApiShopSeckillServiceImpl(
	seckillDao dao.IApiShopSeckillDao,
	seckillConfigDao dao.IApiShopSeckillConfigDao,
) service.IApiShopSeckillService {
	return &IApiShopSeckillServiceImpl{
		seckillDao:       seckillDao,
		seckillConfigDao: seckillConfigDao,
	}
}

// ListConfigs 获取秒杀时间段列表（含当前状态）
func (s *IApiShopSeckillServiceImpl) ListConfigs(c *gin.Context) ([]*models.SeckillConfig, error) {
	data, err := s.seckillConfigDao.List(c, &models.SeckillConfigQuery{
		Status: func() *bool { v := true; return &v }(),
		Page:   1,
		Size:   100,
	})
	if err != nil {
		return nil, err
	}
	// 计算每个时段的当前状态
	for _, cfg := range data.Rows {
		s.computeConfigStatus(cfg)
	}
	return data.Rows, nil
}

// GetCurrentConfig 获取当前秒杀时段配置
func (s *IApiShopSeckillServiceImpl) GetCurrentConfig(c *gin.Context) (*models.SeckillConfig, error) {
	data, err := s.seckillConfigDao.List(c, &models.SeckillConfigQuery{
		Status: func() *bool { v := true; return &v }(),
		Page:   1,
		Size:   100,
	})
	if err != nil {
		return nil, err
	}
	now := time.Now()
	currentMinutes := int64(now.Hour()*60 + now.Minute())

	for _, cfg := range data.Rows {
		beginMinutes := cfg.BeginClock * 60
		endMinutes := (cfg.BeginClock + cfg.ContinueClock) * 60
		if currentMinutes >= beginMinutes && currentMinutes < endMinutes {
			s.computeConfigStatus(cfg)
			return cfg, nil
		}
	}
	return nil, nil
}

// ListGoods 获取秒杀商品列表（按日期范围过滤）
func (s *IApiShopSeckillServiceImpl) ListGoods(c *gin.Context, query *models.SeckillQuery) (*models.SeckillListData, error) {
	// 只查询显示中的、未删除的
	query.IsShow = func() *int32 { v := int32(1); return &v }()
	query.Status = func() *int32 { v := int32(1); return &v }()

	data, err := s.seckillDao.List(c, query)
	if err != nil || data == nil {
		return data, err
	}

	// 按日期范围过滤（start_time <= now <= stop_time）
	now := time.Now()
	filteredRows := make([]*models.Seckill, 0, len(data.Rows))
	for _, goods := range data.Rows {
		if goods == nil {
			continue
		}
		startTime, err1 := myTime.ParseDateTime(goods.StartTime)
		stopTime, err2 := myTime.ParseDateTime(goods.StopTime)
		if err1 != nil || err2 != nil {
			continue
		}
		if now.After(*startTime) && now.Before(stopTime.AddDate(0, 0, 1)) {
			filteredRows = append(filteredRows, goods)
		}
	}

	data.Rows = filteredRows
	data.Total = int64(len(filteredRows))
	return data, nil
}

// GetGoodsDetail 获取秒杀商品详情
func (s *IApiShopSeckillServiceImpl) GetGoodsDetail(c *gin.Context, id int64) (*models.Seckill, error) {
	goods, err := s.seckillDao.GetByID(c, id)
	if err != nil || goods == nil {
		return goods, err
	}
	// 计算库存百分比
	if goods.QuotaShow > 0 {
		goods.Stock = (goods.QuotaShow - goods.Quota) * 100 / goods.QuotaShow
	}
	return goods, nil
}

// computeConfigStatus 计算秒杀时段状态
// status: true=抢购中 false=即将开始/已结束
func (s *IApiShopSeckillServiceImpl) computeConfigStatus(cfg *models.SeckillConfig) {
	now := time.Now()
	currentMinutes := int64(now.Hour()*60 + now.Minute())
	beginMinutes := cfg.BeginClock * 60
	endMinutes := (cfg.BeginClock + cfg.ContinueClock) * 60

	if currentMinutes >= beginMinutes && currentMinutes < endMinutes {
		cfg.Status = true
	} else {
		cfg.Status = false
	}
}
