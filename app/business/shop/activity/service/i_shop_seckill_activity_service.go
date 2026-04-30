package service

import (
	"nova-factory-server/app/business/shop/activity/models"

	"github.com/gin-gonic/gin"
)

type IShopSeckillActivityService interface {
	Set(c *gin.Context, req *models.SeckillActivitySet) (*models.SeckillActivity, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.SeckillActivity, error)
	List(c *gin.Context, req *models.SeckillActivityQuery) (*models.SeckillActivityListData, error)
}
