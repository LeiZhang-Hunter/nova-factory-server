package service

import (
	"nova-factory-server/app/business/shop/config/models"

	"github.com/gin-gonic/gin"
)

// IShopSysConfigService 商城系统配置服务接口
type IShopSysConfigService interface {
	GetWechatConfig(c *gin.Context) (*models.WechatConfigResp, error)
	UpdateWechatConfig(c *gin.Context, req *models.WechatConfigReq) error
}
