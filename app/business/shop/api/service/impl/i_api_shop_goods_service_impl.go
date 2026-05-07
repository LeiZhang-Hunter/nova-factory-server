package impl

import (
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"

	"github.com/gin-gonic/gin"
)

// IApiShopGoodsServiceImpl 商品服务实现
type IApiShopGoodsServiceImpl struct {
	dao dao.IApiShopGoodsDao
}

// NewIApiShopGoodsServiceImpl  创建商品服务
func NewIApiShopGoodsServiceImpl(dao dao.IApiShopGoodsDao) service.IApiShopGoodsService {
	return &IApiShopGoodsServiceImpl{
		dao: dao,
	}
}

// GetByID 获取商品详情
func (s *IApiShopGoodsServiceImpl) GetByID(c *gin.Context, id int64) (*models.Goods, error) {
	return s.dao.GetByID(c, id)
}

// List 获取商品列表
func (s *IApiShopGoodsServiceImpl) List(c *gin.Context, query *models.GoodsQuery) (*models.GoodsListData, error) {
	return s.dao.List(c, query)
}
