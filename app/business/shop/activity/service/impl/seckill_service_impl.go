package impl

import (
	"nova-factory-server/app/business/shop/activity/dao"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/activity/service"
	"nova-factory-server/app/utils/fileUtils"
	"strings"

	"github.com/gin-gonic/gin"
)

type ShopSeckillServiceImpl struct {
	dao dao.IShopSeckillDao
}

func NewShopSeckillService(dao dao.IShopSeckillDao) service.IShopSeckillService {
	return &ShopSeckillServiceImpl{dao: dao}
}

func (s *ShopSeckillServiceImpl) Set(c *gin.Context, req *models.SeckillSet) (*models.Seckill, error) {
	req.Image = strings.TrimSpace(req.Image)
	req.Images = strings.TrimSpace(req.Images)
	req.Title = strings.TrimSpace(req.Title)
	req.Info = strings.TrimSpace(req.Info)
	req.UnitName = strings.TrimSpace(req.UnitName)
	req.StartTime = strings.TrimSpace(req.StartTime)
	req.StopTime = strings.TrimSpace(req.StopTime)
	req.TimeID = strings.TrimSpace(req.TimeID)
	req.Logistics = strings.TrimSpace(req.Logistics)
	req.CustomForm = strings.TrimSpace(req.CustomForm)
	item, err := s.dao.Set(c, req)
	if err != nil || item == nil {
		return item, err
	}
	item.Image = fileUtils.BuildAbsoluteURL(c, item.Image)
	return item, nil
}

func (s *ShopSeckillServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopSeckillServiceImpl) GetByID(c *gin.Context, id int64) (*models.Seckill, error) {
	item, err := s.dao.GetByID(c, id)
	if err != nil || item == nil {
		return item, err
	}
	item.Image = fileUtils.BuildAbsoluteURL(c, item.Image)
	return item, nil
}

func (s *ShopSeckillServiceImpl) List(c *gin.Context, req *models.SeckillQuery) (*models.SeckillListData, error) {
	data, err := s.dao.List(c, req)
	if err != nil {
		return nil, err
	}
	for _, row := range data.Rows {
		if row == nil {
			continue
		}
		row.Image = fileUtils.BuildAbsoluteURL(c, row.Image)
	}
	return data, nil
}
