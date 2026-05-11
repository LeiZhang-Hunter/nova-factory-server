package impl

import (
	"nova-factory-server/app/business/shop/config/dao"
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/business/shop/config/service"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	KeyWechatAppID          = "wechat_mini_program_app_id"
	KeyWechatAppSecret      = "wechat_mini_program_app_secret"
	KeyWechatToken          = "wechat_mini_program_token"
	KeyWechatEncodingAESKey = "wechat_mini_program_encoding_aes_key"
	KeyWechatPayMchID       = "wechat_pay_mch_id"
	KeyWechatPayMchKey      = "wechat_pay_mch_key"
	KeyWechatPayNotifyURL   = "wechat_pay_notify_url"
	KeyWechatPayCertPath    = "wechat_pay_cert_path"
)

type ShopSysConfigServiceImpl struct {
	dao dao.IShopSysConfigDao
}

func NewShopSysConfigService(dao dao.IShopSysConfigDao) service.IShopSysConfigService {
	return &ShopSysConfigServiceImpl{dao: dao}
}

func (s *ShopSysConfigServiceImpl) GetWechatConfig(c *gin.Context) (*models.WechatConfigResp, error) {
	resp := &models.WechatConfigResp{}

	if config, err := s.dao.GetByConfigKey(c, KeyWechatAppID); err == nil {
		resp.AppID = config.ConfigValue
	}
	if config, err := s.dao.GetByConfigKey(c, KeyWechatAppSecret); err == nil {
		resp.AppSecret = maskString(config.ConfigValue)
	}
	if config, err := s.dao.GetByConfigKey(c, KeyWechatToken); err == nil {
		resp.Token = config.ConfigValue
	}
	if config, err := s.dao.GetByConfigKey(c, KeyWechatEncodingAESKey); err == nil {
		resp.EncodingAESKey = config.ConfigValue
	}
	if config, err := s.dao.GetByConfigKey(c, KeyWechatPayMchID); err == nil {
		resp.MchID = config.ConfigValue
	}
	if config, err := s.dao.GetByConfigKey(c, KeyWechatPayMchKey); err == nil {
		resp.MchKey = maskString(config.ConfigValue)
	}
	if config, err := s.dao.GetByConfigKey(c, KeyWechatPayNotifyURL); err == nil {
		resp.NotifyURL = config.ConfigValue
	}
	if config, err := s.dao.GetByConfigKey(c, KeyWechatPayCertPath); err == nil {
		resp.CertPath = config.ConfigValue
	}

	return resp, nil
}

func (s *ShopSysConfigServiceImpl) UpdateWechatConfig(c *gin.Context, req *models.WechatConfigReq) error {
	if req.AppID != "" {
		if err := s.dao.UpdateByConfigKey(c, KeyWechatAppID, req.AppID); err != nil {
			return err
		}
	}
	if req.AppSecret != "" {
		if err := s.dao.UpdateByConfigKey(c, KeyWechatAppSecret, req.AppSecret); err != nil {
			return err
		}
	}
	if req.Token != "" {
		if err := s.dao.UpdateByConfigKey(c, KeyWechatToken, req.Token); err != nil {
			return err
		}
	}
	if req.EncodingAESKey != "" {
		if err := s.dao.UpdateByConfigKey(c, KeyWechatEncodingAESKey, req.EncodingAESKey); err != nil {
			return err
		}
	}
	if req.MchID != "" {
		if err := s.dao.UpdateByConfigKey(c, KeyWechatPayMchID, req.MchID); err != nil {
			return err
		}
	}
	if req.MchKey != "" {
		if err := s.dao.UpdateByConfigKey(c, KeyWechatPayMchKey, req.MchKey); err != nil {
			return err
		}
	}
	if req.NotifyURL != "" {
		if err := s.dao.UpdateByConfigKey(c, KeyWechatPayNotifyURL, req.NotifyURL); err != nil {
			return err
		}
	}
	if req.CertPath != "" {
		if err := s.dao.UpdateByConfigKey(c, KeyWechatPayCertPath, req.CertPath); err != nil {
			return err
		}
	}
	return nil
}

func maskString(s string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) <= 4 {
		return "****"
	}
	return strings.Repeat("*", len(s)-4) + s[len(s)-4:]
}
