package product

import "github.com/gin-gonic/gin"

type Category struct{}

func NewCategory() *Category {
	return &Category{}
}

func (c *Category) PublicRoutes(router *gin.RouterGroup) {

}

func (c *Category) PrivateRoutes(router *gin.RouterGroup) {

}
