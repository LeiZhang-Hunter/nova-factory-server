package masterservice

import "github.com/gin-gonic/gin"

type IStockService interface {
	Create(c *gin.Context) error
	Update(c *gin.Context) error
	UpdateStock(c *gin.Context) error
	GetStock(c *gin.Context) error
}
