package impl

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	activityService "nova-factory-server/app/business/shop/activity/service"
	"nova-factory-server/app/business/shop/home/dao"
	"nova-factory-server/app/business/shop/home/models"
	"nova-factory-server/app/business/shop/home/service"
	"nova-factory-server/app/constant/shop"
	"nova-factory-server/app/utils/fileUtils"

	"github.com/gin-gonic/gin"
)

// ShopHomeModuleItemServiceImpl 提供首页模块明细相关业务能力。
type ShopHomeModuleItemServiceImpl struct {
	dao                  dao.IShopHomeModuleItemDao
	moduleDao            dao.IShopHomeModuleDao
	seckillService       activityService.IShopSeckillActivityService
	seckillConfigService activityService.IShopSeckillConfigService
}

// NewShopHomeModuleItemService 创建首页模块明细服务。
func NewShopHomeModuleItemService(
	dao dao.IShopHomeModuleItemDao,
	moduleDao dao.IShopHomeModuleDao,
	seckillService activityService.IShopSeckillActivityService,
	seckillConfigService activityService.IShopSeckillConfigService,
) service.IShopHomeModuleItemService {
	return &ShopHomeModuleItemServiceImpl{
		dao:                  dao,
		moduleDao:            moduleDao,
		seckillService:       seckillService,
		seckillConfigService: seckillConfigService,
	}
}

// Set 新增或修改首页模块明细。
func (s *ShopHomeModuleItemServiceImpl) Set(c *gin.Context, req *models.HomeModuleItemSet) (*models.HomeModuleItem, error) {
	if req == nil {
		return nil, fmt.Errorf("参数不能为空")
	}
	req.BusinessType = strings.TrimSpace(req.BusinessType)
	req.ItemName = strings.TrimSpace(req.ItemName)
	req.ItemSubTitle = strings.TrimSpace(req.ItemSubTitle)
	req.ItemImage = strings.TrimSpace(req.ItemImage)
	req.ExtJSON = strings.TrimSpace(req.ExtJSON)

	if req.ModuleID <= 0 {
		return nil, fmt.Errorf("模块ID不能为空")
	}
	module, err := s.moduleDao.GetByID(c, req.ModuleID)
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, fmt.Errorf("首页模块不存在")
	}
	if req.BusinessType == "" {
		return nil, fmt.Errorf("业务类型不能为空")
	}
	if req.LinkID <= 0 {
		return nil, fmt.Errorf("关联ID不能为空")
	}
	if req.Sort < 0 {
		return nil, fmt.Errorf("排序不能小于0")
	}
	return s.dao.Set(c, req)
}

// DeleteByIDs 删除首页模块明细。
func (s *ShopHomeModuleItemServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

// GetByID 根据主键获取首页模块明细详情。
func (s *ShopHomeModuleItemServiceImpl) GetByID(c *gin.Context, id int64) (*models.HomeModuleItem, error) {
	data, err := s.dao.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if data != nil {
		if err := s.decorateItems(c, []*models.HomeModuleItem{data}); err != nil {
			return nil, err
		}
	}
	return data, nil
}

// ListByModuleIDs 按模块ID批量查询首页模块明细。
func (s *ShopHomeModuleItemServiceImpl) ListByModuleIDs(c *gin.Context, moduleIDs []int64) ([]*models.HomeModuleItem, error) {

	rows, err := s.dao.ListByModuleIDs(c, moduleIDs)
	if err != nil {
		return nil, err
	}
	if err := s.decorateItems(c, rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// List 查询首页模块明细列表。
func (s *ShopHomeModuleItemServiceImpl) List(c *gin.Context, req *models.HomeModuleItemQuery) (*models.HomeModuleItemListData, error) {
	data, err := s.dao.List(c, req)
	if err != nil {
		return nil, err
	}
	if err := s.decorateItems(c, data.Rows); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *ShopHomeModuleItemServiceImpl) fillSeckillItemImages(c *gin.Context, rows []*models.HomeModuleItem) error {
	seckillIDs := make([]int64, 0)
	for _, row := range rows {
		if row == nil || row.BusinessType != shop.SeckillActivity || row.LinkID <= 0 {
			continue
		}
		seckillIDs = append(seckillIDs, row.LinkID)
	}
	if len(seckillIDs) == 0 {
		return nil
	}

	activities, err := s.seckillService.GetByIDs(c, seckillIDs)
	if err != nil {
		return err
	}
	configIDs := make([]int64, 0, len(activities))
	activityConfigMap := make(map[int64][]int64, len(activities))
	configSeen := make(map[int64]struct{}, len(activities))
	for _, activity := range activities {
		if activity == nil {
			continue
		}
		timeIDs := parseSeckillConfigIDs(activity.TimeIDs)
		if len(timeIDs) == 0 {
			continue
		}
		activityConfigMap[activity.ID] = timeIDs
		for _, configID := range timeIDs {
			if _, ok := configSeen[configID]; ok {
				continue
			}
			configSeen[configID] = struct{}{}
			configIDs = append(configIDs, configID)
		}
	}
	if len(configIDs) == 0 {
		return nil
	}

	configs, err := s.seckillConfigService.GetByIDs(c, configIDs)
	if err != nil {
		return err
	}
	now := time.Now()
	configImageMap := make(map[int64]string, len(configs))
	for _, config := range configs {
		if config == nil {
			continue
		}
		if !isSeckillConfigActive(now, config.BeginClock, config.ContinueClock) {
			continue
		}
		configImageMap[config.ID] = config.Images
	}
	for _, row := range rows {
		if row == nil || row.BusinessType != shop.SeckillActivity {
			continue
		}
		configIDs, ok := activityConfigMap[row.LinkID]
		if !ok {
			continue
		}
		for _, configID := range configIDs {
			if image, ok := configImageMap[configID]; ok && image != "" {
				row.ItemImage = image
				break
			}
		}
	}
	return nil
}

func parseSeckillConfigIDs(raw string) []int64 {
	parts := strings.Split(strings.TrimSpace(raw), ",")
	result := make([]int64, 0, len(parts))
	seen := make(map[int64]struct{}, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil || id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

func isSeckillConfigActive(now time.Time, beginClock, continueClock int64) bool {
	if continueClock <= 0 {
		return false
	}
	currentMinute := int64(now.Hour()*60 + now.Minute())
	startMinute := beginClock * 60
	endMinute := (beginClock + continueClock) * 60
	return currentMinute >= startMinute && currentMinute < endMinute
}

func (s *ShopHomeModuleItemServiceImpl) decorateItems(c *gin.Context, rows []*models.HomeModuleItem) error {
	if err := s.fillSeckillItemImages(c, rows); err != nil {
		return err
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		row.ItemImage = fileUtils.BuildAbsoluteURL(c, row.ItemImage)
	}
	return nil
}
