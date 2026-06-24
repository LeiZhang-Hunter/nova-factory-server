package impl

import (
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"

	"github.com/gin-gonic/gin"
)

// IApiShopFavoriteServiceImpl 用户商品收藏服务实现
type IApiShopFavoriteServiceImpl struct {
	favoriteDao dao.IApiShopFavoriteDao
	goodsDao    dao.IApiShopGoodsDao
}

// NewIApiShopFavoriteServiceImpl 创建用户商品收藏服务
func NewIApiShopFavoriteServiceImpl(favoriteDao dao.IApiShopFavoriteDao, goodsDao dao.IApiShopGoodsDao) service.IApiShopFavoriteService {
	return &IApiShopFavoriteServiceImpl{
		favoriteDao: favoriteDao,
		goodsDao:    goodsDao,
	}
}

// AddFavorite 添加收藏
func (s *IApiShopFavoriteServiceImpl) AddFavorite(c *gin.Context, userId int64, goodsId int64) error {
	return s.favoriteDao.Add(c, userId, goodsId)
}

// RemoveFavorite 移除收藏
func (s *IApiShopFavoriteServiceImpl) RemoveFavorite(c *gin.Context, userId int64, goodsId int64) error {
	return s.favoriteDao.Remove(c, userId, goodsId)
}

// ListFavorites 获取收藏列表
func (s *IApiShopFavoriteServiceImpl) ListFavorites(c *gin.Context, userId int64, page int64, size int64, goodsName string) (*models.GoodsListData, error) {
	favorites, total, err := s.favoriteDao.ListByUserID(c, userId, page, size, goodsName)
	if err != nil {
		return nil, err
	}

	goodsList := make([]*models.Goods, 0, len(favorites))
	for _, fav := range favorites {
		goods, err := s.goodsDao.GetByGoodsID(c, fav.GoodsID)
		if err != nil {
			continue
		}
		if goods != nil {
			goodsList = append(goodsList, goods)
		}
	}

	return &models.GoodsListData{
		Rows:  goodsList,
		Total: total,
	}, nil
}

// CheckFavoriteStatus 检查收藏状态
func (s *IApiShopFavoriteServiceImpl) CheckFavoriteStatus(c *gin.Context, userId int64, goodsId string) (bool, error) {
	favorite, err := s.favoriteDao.GetByUserAndGoods(c, userId, goodsId)
	if err != nil {
		return false, err
	}
	return favorite != nil, nil
}
