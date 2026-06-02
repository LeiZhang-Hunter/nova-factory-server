package order

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/service"

	"github.com/gin-gonic/gin"
	gopayWechat "github.com/go-pay/gopay/wechat/v3"
)

// OrderNotify 微信支付回调控制器
type OrderNotify struct {
	service   service.IApiShopOrderService
	configDao dao.IApiShopSysConfigDao
}

// NewOrderNotify 创建微信支付回调控制器。
func NewOrderNotify(service service.IApiShopOrderService, configDao dao.IApiShopSysConfigDao) *OrderNotify {
	return &OrderNotify{service: service, configDao: configDao}
}

// PublicRoutes 注册微信支付回调路由到 publicGroup（不鉴权）。
func (s *OrderNotify) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop")
	group.POST("/order/notify", s.HandleWechatNotify)
}

var notifyConfigKeys = []string{
	"wechat_mini_program_app_id",
	"wechat_pay_mch_id",
	"wechat_pay_api_v3_key",
	"wechat_pay_serial_no",
	"wechat_pay_platform_public_key_path",
}

// HandleWechatNotify 微信支付异步回调。
func (s *OrderNotify) HandleWechatNotify(c *gin.Context) {
	// 解析回调
	notifyReq, err := gopayWechat.V3ParseNotify(c.Request)
	if err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "解析失败"})
		return
	}

	// 读微信配置
	cfgMap, err := s.loadNotifyConfig(c)
	if err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "配置读取失败"})
		return
	}
	appId := cfgMap["wechat_mini_program_app_id"]
	mchId := cfgMap["wechat_pay_mch_id"]
	apiV3Key := cfgMap["wechat_pay_api_v3_key"]
	serialNo := cfgMap["wechat_pay_serial_no"]
	platformPublicKeyPath := cfgMap["wechat_pay_platform_public_key_path"]

	// 平台公钥验签
	platformKeyData, err := os.ReadFile(platformPublicKeyPath)
	if err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "公钥读取失败"})
		return
	}
	publicKey, err := parseRSAPublicKey(platformKeyData)
	if err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "公钥解析失败"})
		return
	}
	certMap := map[string]*rsa.PublicKey{"PUB_KEY_ID_" + serialNo: publicKey}
	if err := notifyReq.VerifySignByPKMap(certMap); err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "验签失败"})
		return
	}

	// 解密回调内容
	type notifyData struct {
		OutTradeNo    string `json:"out_trade_no"`
		TransactionID string `json:"transaction_id"`
		TradeState    string `json:"trade_state"`
		AppID         string `json:"appid"`
		MchID         string `json:"mchid"`
		Amount        struct {
			Total int64 `json:"total"`
		} `json:"amount"`
		Payer struct {
			Openid string `json:"openid"`
		} `json:"payer"`
	}
	var nd notifyData
	if err := notifyReq.DecryptCipherTextToStruct(apiV3Key, &nd); err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "解密失败"})
		return
	}

	// 校验
	if nd.TradeState != "SUCCESS" {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "SUCCESS", Message: "非支付成功通知"})
		return
	}
	if nd.AppID != appId {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "appid不匹配"})
		return
	}
	if nd.MchID != mchId {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "mchid不匹配"})
		return
	}

	// 序列化回调原文（排障用）
	rawBytes, _ := json.Marshal(nd)
	notifyRaw := string(rawBytes)

	// 更新订单状态
	if err := s.service.HandleWechatNotify(c, nd.OutTradeNo, nd.TransactionID, notifyRaw, nd.MchID, nd.AppID, nd.Payer.Openid, nd.Amount.Total); err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "SUCCESS", Message: "成功"})
}

func (s *OrderNotify) loadNotifyConfig(c *gin.Context) (map[string]string, error) {
	rows, err := s.configDao.GetByConfigKeys(c, notifyConfigKeys)
	if err != nil {
		return nil, fmt.Errorf("读取微信支付配置失败: %v", err)
	}
	cfgMap := make(map[string]string)
	for _, row := range rows {
		cfgMap[row.ConfigKey] = row.ConfigValue
	}
	return cfgMap, nil
}

func parseRSAPublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid public key pem")
	}
	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PublicKey); ok {
			return rsaKey, nil
		}
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}
