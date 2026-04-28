package dao

import (
	"nova-factory-server/app/business/shop/activity/models"

	"github.com/gin-gonic/gin"
)

type IShopCombinationDao interface {
	Set(c *gin.Context, req *models.CombinationSet) (*models.Combination, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.Combination, error)
	List(c *gin.Context, req *models.CombinationQuery) (*models.CombinationListData, error)
}
