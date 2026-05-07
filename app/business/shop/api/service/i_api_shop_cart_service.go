package service

import (
	"github.com/gin-gonic/gin"
	api "nova-factory-server/app/business/shop/api/models"
)

// IApiShopCartService api 商城服务
type IApiShopCartService interface {
	GenCart(c *gin.Context, req *api.CartSetDataReq) (*api.CartDto, error)
	List(c *gin.Context) ([]*api.CartDto, error)
}
