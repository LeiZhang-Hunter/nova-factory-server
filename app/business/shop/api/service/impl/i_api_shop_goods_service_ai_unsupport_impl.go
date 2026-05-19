//go:build !ai
// +build !ai

package impl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	discountservice "nova-factory-server/app/business/shop/discount/service"
	"nova-factory-server/app/business/shop/product/shopservice"
)

// IApiShopGoodsServiceImpl 商品服务实现
type IApiShopGoodsServiceImpl struct {
	dao              dao.IApiShopGoodsDao
	shopGoodsService shopservice.IShopGoodsService
	discountService  discountservice.IDiscountCalculateService
}

// NewIApiShopGoodsServiceImpl  创建商品服务
func NewIApiShopGoodsServiceImpl(dao dao.IApiShopGoodsDao, shopGoodsService shopservice.IShopGoodsService, discountService discountservice.IDiscountCalculateService) service.IApiShopGoodsService {
	return &IApiShopGoodsServiceImpl{
		dao:              dao,
		shopGoodsService: shopGoodsService,
		discountService:  discountService,
	}
}

// Search 按多个商品名称检索相似商品，并回填数据库中的最新商品数据
func (s *IApiShopGoodsServiceImpl) Search(c *gin.Context, req *models.GoodsSearchReq) (*models.GoodsSearchData, error) {

	return &models.GoodsSearchData{
		Rows:  make([]*models.GoodsSearchItem, 0),
		Total: 0,
	}, nil
}
