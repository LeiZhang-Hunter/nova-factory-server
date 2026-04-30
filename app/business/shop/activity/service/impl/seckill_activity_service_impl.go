package impl

import (
	"nova-factory-server/app/business/shop/activity/dao"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/activity/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type ShopSeckillActivityServiceImpl struct {
	dao dao.IShopSeckillActivityDao
}

func NewShopSeckillActivityService(dao dao.IShopSeckillActivityDao) service.IShopSeckillActivityService {
	return &ShopSeckillActivityServiceImpl{dao: dao}
}

func (s *ShopSeckillActivityServiceImpl) Set(c *gin.Context, req *models.SeckillActivitySet) (*models.SeckillActivity, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.TimeIDs = normalizeCSV(req.TimeIDs)
	return s.dao.Set(c, req)
}

func (s *ShopSeckillActivityServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopSeckillActivityServiceImpl) GetByID(c *gin.Context, id int64) (*models.SeckillActivity, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopSeckillActivityServiceImpl) List(c *gin.Context, req *models.SeckillActivityQuery) (*models.SeckillActivityListData, error) {
	return s.dao.List(c, req)
}

func normalizeCSV(raw string) string {
	parts := strings.Split(strings.TrimSpace(raw), ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		result = append(result, part)
	}
	return strings.Join(result, ",")
}
