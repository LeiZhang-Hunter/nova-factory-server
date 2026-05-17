package impl

import (
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"

	"github.com/gin-gonic/gin"
)

// IApiShopCombinationServiceImpl 拼团服务实现
type IApiShopCombinationServiceImpl struct {
	combinationDao dao.IApiShopCombinationDao
	pinkDao        dao.IApiShopPinkDao
}

// NewIApiShopCombinationServiceImpl 创建拼团服务
func NewIApiShopCombinationServiceImpl(
	combinationDao dao.IApiShopCombinationDao,
	pinkDao dao.IApiShopPinkDao,
) service.IApiShopCombinationService {
	return &IApiShopCombinationServiceImpl{
		combinationDao: combinationDao,
		pinkDao:        pinkDao,
	}
}

// List 获取拼团商品列表
func (s *IApiShopCombinationServiceImpl) List(c *gin.Context, query *models.CombinationQuery) (*models.CombinationListData, error) {
	// 只查询显示中的、未删除的
	query.IsShow = func() *int32 { v := int32(1); return &v }()
	return s.combinationDao.List(c, query)
}

// GetByID 获取拼团商品详情（含进行中的团）
func (s *IApiShopCombinationServiceImpl) GetByID(c *gin.Context, id int64) (*models.Combination, error) {
	return s.combinationDao.GetByID(c, id)
}

// GetPinkList 获取正在进行中的团列表
func (s *IApiShopCombinationServiceImpl) GetPinkList(c *gin.Context, cid int64) (*models.PinkListData, error) {
	status := int32(1)
	isRefund := int32(0)
	return s.pinkDao.List(c, &models.PinkQuery{
		CID:      cid,
		Status:   &status,
		IsRefund: &isRefund,
		Page:     1,
		Size:     100,
	})
}

// IApiShopPinkServiceImpl 拼团记录服务实现
type IApiShopPinkServiceImpl struct {
	pinkDao dao.IApiShopPinkDao
}

// NewIApiShopPinkServiceImpl 创建拼团记录服务
func NewIApiShopPinkServiceImpl(pinkDao dao.IApiShopPinkDao) service.IApiShopPinkService {
	return &IApiShopPinkServiceImpl{
		pinkDao: pinkDao,
	}
}

// Create 创建拼团记录（开团或参团）
// 注意：实际开团由订单支付成功回调触发，这里只做参数校验
func (s *IApiShopPinkServiceImpl) Create(c *gin.Context, userID int64, combinationID int64, pinkID int64, orderID int64) (*models.Pink, error) {
	// 参数校验由 Controller 层处理
	// 实际拼团记录创建在订单支付回调中处理
	return nil, nil
}

// GetDetail 获取团详情（含当前人数）
func (s *IApiShopPinkServiceImpl) GetDetail(c *gin.Context, pinkID int64) (*models.Pink, error) {
	return s.pinkDao.GetByID(c, pinkID)
}

// GetPinkMemberCount 获取团内人数
func (s *IApiShopPinkServiceImpl) GetPinkMemberCount(c *gin.Context, pinkID int64) (int64, error) {
	return s.pinkDao.CountMembers(c, pinkID)
}

// ListMyPinks 获取用户的拼团记录
func (s *IApiShopPinkServiceImpl) ListMyPinks(c *gin.Context, userID int64) (*models.PinkListData, error) {
	return s.pinkDao.List(c, &models.PinkQuery{
		UID: userID,
	})
}
