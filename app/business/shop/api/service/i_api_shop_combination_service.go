package service

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopCombinationService 拼团服务接口
type IApiShopCombinationService interface {
	// List 获取拼团商品列表
	List(c *gin.Context, query *models.CombinationQuery) (*models.CombinationListData, error)
	// GetByID 获取拼团商品详情
	GetByID(c *gin.Context, id int64) (*models.Combination, error)
	// GetPinkList 获取正在进行中的团列表
	GetPinkList(c *gin.Context, cid int64) (*models.PinkListData, error)
}

// IApiShopPinkService 拼团记录服务接口
type IApiShopPinkService interface {
	// Create 创建拼团记录（开团或参团）
	Create(c *gin.Context, userID int64, combinationID int64, pinkID int64, orderID int64) (*models.Pink, error)
	// GetDetail 获取团详情（含当前人数）
	GetDetail(c *gin.Context, pinkID int64) (*models.Pink, error)
	// GetPinkMemberCount 获取团内人数
	GetPinkMemberCount(c *gin.Context, pinkID int64) (int64, error)
	// ListMyPinks 获取用户的拼团记录
	ListMyPinks(c *gin.Context, userID int64) (*models.PinkListData, error)
}
