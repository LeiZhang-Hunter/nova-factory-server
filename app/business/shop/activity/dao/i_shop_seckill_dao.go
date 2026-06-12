package dao

import (
	"nova-factory-server/app/business/shop/activity/models"

	"github.com/gin-gonic/gin"
)

type IShopSeckillDao interface {
	Set(c *gin.Context, req *models.SeckillSet) (*models.Seckill, error)
	BatchCreate(c *gin.Context, reqs []*models.SeckillSet, batchSize int) error
	BatchUpdate(c *gin.Context, reqs []*models.SeckillSet, batchSize int) error
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.Seckill, error)
	List(c *gin.Context, req *models.SeckillQuery) (*models.SeckillListData, error)
	DeductStock(c *gin.Context, id int64, quantity int64) error
	RestoreStock(c *gin.Context, id int64, quantity int64) error
}
