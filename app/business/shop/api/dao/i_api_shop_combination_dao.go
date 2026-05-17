package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IApiShopCombinationDao 拼团商品数据访问接口
type IApiShopCombinationDao interface {
	List(c *gin.Context, query *models.CombinationQuery) (*models.CombinationListData, error)
	GetByID(c *gin.Context, id int64) (*models.Combination, error)
}
