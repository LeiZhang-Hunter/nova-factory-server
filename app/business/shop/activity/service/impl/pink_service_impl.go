package impl

import (
	"nova-factory-server/app/business/shop/activity/dao"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/activity/service"
	"nova-factory-server/app/utils/fileUtils"

	"github.com/gin-gonic/gin"
)

type ShopPinkServiceImpl struct {
	dao dao.IShopPinkDao
}

func NewShopPinkService(dao dao.IShopPinkDao) service.IShopPinkService {
	return &ShopPinkServiceImpl{dao: dao}
}

func (s *ShopPinkServiceImpl) GetByID(c *gin.Context, id int64) (*models.Pink, error) {
	item, err := s.dao.GetByID(c, id)
	if err != nil || item == nil {
		return item, err
	}
	item.Avatar = fileUtils.BuildAbsoluteURL(c, item.Avatar)
	item.CombinationImage = fileUtils.BuildAbsoluteURL(c, item.CombinationImage)
	return item, nil
}

func (s *ShopPinkServiceImpl) List(c *gin.Context, req *models.PinkQuery) (*models.PinkListData, error) {
	data, err := s.dao.List(c, req)
	if err != nil {
		return nil, err
	}
	for _, row := range data.Rows {
		if row == nil {
			continue
		}
		row.Avatar = fileUtils.BuildAbsoluteURL(c, row.Avatar)
		row.CombinationImage = fileUtils.BuildAbsoluteURL(c, row.CombinationImage)
	}
	return data, nil
}
