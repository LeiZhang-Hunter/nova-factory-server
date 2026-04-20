package product

import "github.com/gin-gonic/gin"

type Product struct{}

func NewProduct() *Product {
	return &Product{}
}

func (p *Product) PublicRoutes(router *gin.RouterGroup) {

}

func (p *Product) PrivateRoutes(router *gin.RouterGroup) {

}
