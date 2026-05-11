package dao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/shop/api/models"
)

type IApiShopCartDao interface {
	Save(c *gin.Context, req *models.CartSetData) (*models.CartDto, error)
	List(c *gin.Context, userID int64) ([]*models.CartDto, error)
	ListByIDs(c *gin.Context, userID int64, ids []int64) ([]*models.CartDto, error)
	Remove(c *gin.Context, ids []string) error
}
