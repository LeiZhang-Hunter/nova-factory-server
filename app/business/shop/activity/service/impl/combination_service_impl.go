package impl

import (
	"nova-factory-server/app/business/shop/activity/dao"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/business/shop/activity/service"
	"nova-factory-server/app/utils/fileUtils"
	"strings"

	"github.com/gin-gonic/gin"
)

type ShopCombinationServiceImpl struct {
	dao dao.IShopCombinationDao
}

func NewShopCombinationService(dao dao.IShopCombinationDao) service.IShopCombinationService {
	return &ShopCombinationServiceImpl{dao: dao}
}

func (s *ShopCombinationServiceImpl) Set(c *gin.Context, req *models.CombinationSet) (*models.Combination, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Image = strings.TrimSpace(req.Image)
	req.Images = strings.TrimSpace(req.Images)
	req.Attr = strings.TrimSpace(req.Attr)
	req.Info = strings.TrimSpace(req.Info)
	req.UnitName = strings.TrimSpace(req.UnitName)
	item, err := s.dao.Set(c, req)
	if err != nil || item == nil {
		return item, err
	}
	item.Image = fileUtils.BuildAbsoluteURL(c, item.Image)
	return item, nil
}

func (s *ShopCombinationServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopCombinationServiceImpl) GetByID(c *gin.Context, id int64) (*models.Combination, error) {
	item, err := s.dao.GetByID(c, id)
	if err != nil || item == nil {
		return item, err
	}
	item.Image = fileUtils.BuildAbsoluteURL(c, item.Image)
	return item, nil
}

func (s *ShopCombinationServiceImpl) List(c *gin.Context, req *models.CombinationQuery) (*models.CombinationListData, error) {
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
