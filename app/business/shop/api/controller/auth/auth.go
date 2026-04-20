package auth

import (
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Auth 商城用户鉴权控制器
type Auth struct {
	service service.IShopAuthService
}

// NewAuth 创建商城用户鉴权控制器。
func NewAuth(service service.IShopAuthService) *Auth {
	return &Auth{service: service}
}

// PublicRoutes 注册商城用户公开鉴权路由。
func (s *Auth) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/app")
	group.POST("/login", s.Login)
}

// PrivateRoutes 注册商城用户受保护鉴权路由。
func (s *Auth) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/app")
	group.GET("/getInfo", s.GetInfo)
	group.POST("/logout", s.Logout)
}

// Login 商城用户登录
// @Summary 商城用户登录
// @Description 使用 shop_user 账号登录商城
// @Tags 商城/App鉴权
// @Accept application/json
// @Param body body models.ShopLoginReq true "商城用户登录参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "登录成功"
// @Router /shop/app/login [post]
func (s *Auth) Login(c *gin.Context) {
	req := new(models.ShopLoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Login(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetInfo 获取当前商城登录用户信息
// @Summary 获取当前商城登录用户信息
// @Description 获取当前商城登录用户信息及权限
// @Tags 商城/App鉴权
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/app/getInfo [get]
func (s *Auth) GetInfo(c *gin.Context) {
	data, err := s.service.GetInfo(c)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Logout 退出当前商城登录会话
// @Summary 退出当前商城登录会话
// @Description 退出当前商城登录会话
// @Tags 商城/App鉴权
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "退出成功"
// @Router /shop/app/logout [post]
func (s *Auth) Logout(c *gin.Context) {
	if err := s.service.Logout(c); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
