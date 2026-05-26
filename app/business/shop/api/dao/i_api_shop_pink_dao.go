package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopPinkDao 拼团记录数据访问接口
type IApiShopPinkDao interface {
	List(c *gin.Context, query *models.PinkQuery) (*models.PinkListData, error)
	GetByID(c *gin.Context, id int64) (*models.Pink, error)
	CountMembers(c *gin.Context, pinkID int64) (int64, error)
}
