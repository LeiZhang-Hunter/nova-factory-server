package auth

import (
	"nova-factory-server/app/business/shop/api/models"
	shopAuthService "nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Auth 商城用户鉴权控制器
type Auth struct {
	authService       shopAuthService.IApiShopAuthService
	wechatAuthService shopAuthService.IApiShopWechatAuthService
}

// NewAuth 创建商城用户鉴权控制器。
func NewAuth(authService shopAuthService.IApiShopAuthService, wechatAuthService shopAuthService.IApiShopWechatAuthService) *Auth {
	return &Auth{
		authService:       authService,
		wechatAuthService: wechatAuthService,
	}
}

// PublicRoutes 注册商城用户公开鉴权路由。
func (s *Auth) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/auth")
	group.POST("/wechat-login", s.WechatLogin)
	group.POST("/account-login", s.AccountLogin)
	group.POST("/refresh", s.RefreshToken)
}

// PrivateRoutes 注册商城用户受保护鉴权路由。
func (s *Auth) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop")
	group.GET("/getInfo", s.GetInfo)
	group.POST("/logout", s.Logout)
}

// WechatLogin 微信小程序授权登录
// @Summary 微信小程序授权登录
// @Description 微信小程序授权登录，code换openid并生成JWT token
// @Tags 商城/App鉴权
// @Accept application/json
// @Param body body models.WechatLoginReq true "微信登录参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "登录成功"
// @Router /api/v1/app/shop/auth/wechat-login [post]
func (s *Auth) WechatLogin(c *gin.Context) {
	req := new(models.WechatLoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.wechatAuthService.WechatLogin(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// AccountLogin 商城小程序账号密码登录
// @Summary 商城小程序账号密码登录
// @Description 使用 shop_user 账号密码登录并生成商城小程序会话 token
// @Tags 商城/App鉴权
// @Accept application/json
// @Param body body models.AccountLoginReq true "账号密码登录参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "登录成功"
// @Router /api/v1/app/shop/auth/account-login [post]
func (s *Auth) AccountLogin(c *gin.Context) {
	req := new(models.AccountLoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.wechatAuthService.AccountLogin(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// RefreshToken 刷新JWT Token
// @Summary 刷新JWT Token
// @Description 用有效JWT换取新Token
// @Tags 商城/App鉴权
// @Accept application/json
// @Param body body models.RefreshTokenReq true "刷新Token参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "刷新成功"
// @Router /api/v1/app/shop/auth/refresh [post]
func (s *Auth) RefreshToken(c *gin.Context) {
	req := new(models.RefreshTokenReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.wechatAuthService.RefreshToken(c, req)
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
// @Router /api/v1/app/shop/getInfo [get]
func (s *Auth) GetInfo(c *gin.Context) {
	data, err := s.authService.GetInfo(c)
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
// @Router /api/v1/app/shop/logout [post]
func (s *Auth) Logout(c *gin.Context) {
	if err := s.authService.Logout(c); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
