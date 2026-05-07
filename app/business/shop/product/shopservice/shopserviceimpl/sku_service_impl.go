package shopserviceimpl

import (
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/utils/fileUtils"

	"github.com/gin-gonic/gin"
)

type ShopSkuServiceImpl struct {
	dao shopdao.IShopSkuDao
}

func NewShopSkuService(dao shopdao.IShopSkuDao) shopservice.IShopSkuService {
	return &ShopSkuServiceImpl{dao: dao}
}

func (s *ShopSkuServiceImpl) Create(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error) {
	return s.dao.Create(c, req)
}

func (s *ShopSkuServiceImpl) Update(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error) {
	return s.dao.Update(c, req)
}

func (s *ShopSkuServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopSkuServiceImpl) GetByID(c *gin.Context, id int64) (*shopmodels.GoodsSku, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopSkuServiceImpl) List(c *gin.Context, req *shopmodels.GoodsSkuQuery) (*shopmodels.GoodsSkuListData, error) {
	ret, err := s.dao.List(c, req)
	if err != nil {
		return nil, err
	}

	if ret == nil || len(ret.Rows) == 0 {
		return ret, nil
	}

	for k, v := range ret.Rows {
		ret.Rows[k].ImageURL = fileUtils.BuildAbsoluteURL(c, v.ImageURL)
	}
	return ret, nil
}
