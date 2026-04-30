package impl

import (
	"fmt"
	"nova-factory-server/app/utils/fileUtils"
	"strings"

	"nova-factory-server/app/business/shop/activity/dao"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/activity/service"

	"github.com/gin-gonic/gin"
)

// ShopSeckillConfigServiceImpl 提供商品秒杀配置相关业务能力。
type ShopSeckillConfigServiceImpl struct {
	dao dao.IShopSeckillConfigDao
}

// NewShopSeckillConfigService 创建商品秒杀配置服务。
func NewShopSeckillConfigService(dao dao.IShopSeckillConfigDao) service.IShopSeckillConfigService {
	return &ShopSeckillConfigServiceImpl{dao: dao}
}

// Set 新增或修改商品秒杀配置。
func (s *ShopSeckillConfigServiceImpl) Set(c *gin.Context, req *models.SeckillConfigSet) (*models.SeckillConfig, error) {
	if req == nil {
		return nil, fmt.Errorf("参数不能为空")
	}
	req.Images = strings.TrimSpace(req.Images)
	if req.BeginClock < 0 {
		return nil, fmt.Errorf("开启时间不能小于0")
	}
	if req.ContinueClock < 0 {
		return nil, fmt.Errorf("持续时间不能小于0")
	}
	if req.Sort < 0 {
		return nil, fmt.Errorf("排序不能小于0")
	}
	return s.dao.Set(c, req)
}

// DeleteByIDs 删除商品秒杀配置。
func (s *ShopSeckillConfigServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

// GetByID 根据主键查询商品秒杀配置。
func (s *ShopSeckillConfigServiceImpl) GetByID(c *gin.Context, id int64) (*models.SeckillConfig, error) {
	return s.dao.GetByID(c, id)
}

// List 查询商品秒杀配置列表。
func (s *ShopSeckillConfigServiceImpl) List(c *gin.Context, req *models.SeckillConfigQuery) (*models.SeckillConfigListData, error) {
	data, err := s.dao.List(c, req)
	if err != nil {
		return data, err
	}

	for k, v := range data.Rows {
		data.Rows[k].Images = fileUtils.BuildAbsoluteURL(c, v.Images)
	}

	return data, nil
}
