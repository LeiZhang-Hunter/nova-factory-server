package dao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/shop/api/models"
)

type IApiShopCartDao interface {
	Save(c *gin.Context, req *models.CartSetData) (*models.CartDto, error)
}
