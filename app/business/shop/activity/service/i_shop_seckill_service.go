package service

import (
	"nova-factory-server/app/business/shop/activity/models"

	"github.com/gin-gonic/gin"
)

type IShopSeckillService interface {
	Set(c *gin.Context, req *models.SeckillSet) (*models.Seckill, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.Seckill, error)
	List(c *gin.Context, req *models.SeckillQuery) (*models.SeckillListData, error)
}
