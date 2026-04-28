package dao

import (
	"nova-factory-server/app/business/shop/activity/models"

	"github.com/gin-gonic/gin"
)

type IShopPinkDao interface {
	GetByID(c *gin.Context, id int64) (*models.Pink, error)
	List(c *gin.Context, req *models.PinkQuery) (*models.PinkListData, error)
}
