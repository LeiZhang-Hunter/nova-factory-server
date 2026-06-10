package order

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nova-factory-server/app/datasource/objectFile"
	"os"
	"time"

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
	group.Any("/order/notify", s.HandleWechatNotify)
}

var notifyConfigKeys = []string{
	"wechat_mini_program_app_id",
	"wechat_pay_mch_id",
	"wechat_pay_api_v3_key",
	"wechat_pay_serial_no",
	"wechat_pay_platform_public_key_path",
	"wechat_pay_platform_public_key_id",
	"wechat_pay_private_key_path",
}

// HandleWechatNotify 微信支付异步回调。
func (s *OrderNotify) HandleWechatNotify(c *gin.Context) {
	// 读取原始请求体
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "读取请求体失败"})
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if err := writeWechatNotifyRequestFile(c, bodyBytes); err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "写入请求数据失败"})
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
	privateKeyPath := cfgMap["wechat_pay_private_key_path"]
	platformPublicKeyPath := cfgMap["wechat_pay_platform_public_key_path"]
	platformPublicKeyID := cfgMap["wechat_pay_platform_public_key_id"]
	//NewClientV3 初始化微信客户端 v3
	// mchid：商户ID 或者服务商模式的 sp_mchid
	// serialNo：商户证书的证书序列号
	// apiV3Key：apiV3Key，商户平台获取
	// privateKey：私钥 apiclient_key.pem 读取后的内容
	file := objectFile.NewConfig()
	privateKeyData, err := file.ReadPrivateFile(c, privateKeyPath)
	if err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "私钥读取失败"})
		return
	}
	client, err := gopayWechat.NewClientV3(mchId, serialNo, apiV3Key, string(privateKeyData))
	if err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "初始化微信客户端失败"})
		return
	}
	platformPublicKeyData, err := file.ReadPrivateFile(c, platformPublicKeyPath)
	if err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "微信支付公钥读取失败"})
		return
	}
	if err := client.AutoVerifySignByPublicKey(platformPublicKeyData, platformPublicKeyID); err != nil {
		_ = writeWechatNotifyErrorFile("auto_verify_sign_by_public_key", err)
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "验签初始化失败"})
		return
	}

	// 解析回调
	notifyReq, err := gopayWechat.V3ParseNotify(c.Request)
	if err != nil {
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: "解析失败"})
		return
	}
	// 获取微信平台证书
	certMap := client.WxPublicKeyMap()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		_ = writeWechatNotifyErrorFile("verify_sign", err)
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

func writeWechatNotifyRequestFile(c *gin.Context, bodyBytes []byte) error {
	record := struct {
		Time          string              `json:"time"`
		Method        string              `json:"method"`
		Scheme        string              `json:"scheme"`
		Host          string              `json:"host"`
		Path          string              `json:"path"`
		RawQuery      string              `json:"rawQuery"`
		RequestURI    string              `json:"requestUri"`
		Proto         string              `json:"proto"`
		RemoteAddr    string              `json:"remoteAddr"`
		ClientIP      string              `json:"clientIp"`
		ContentLength int64               `json:"contentLength"`
		Headers       map[string][]string `json:"headers"`
		Query         map[string][]string `json:"query"`
		Body          string              `json:"body"`
		BodyBase64    string              `json:"bodyBase64"`
	}{
		Time:          time.Now().Format(time.RFC3339Nano),
		Method:        c.Request.Method,
		Scheme:        c.Request.URL.Scheme,
		Host:          c.Request.Host,
		Path:          c.Request.URL.Path,
		RawQuery:      c.Request.URL.RawQuery,
		RequestURI:    c.Request.RequestURI,
		Proto:         c.Request.Proto,
		RemoteAddr:    c.Request.RemoteAddr,
		ClientIP:      c.ClientIP(),
		ContentLength: c.Request.ContentLength,
		Headers:       c.Request.Header,
		Query:         c.Request.URL.Query(),
		Body:          string(bodyBytes),
		BodyBase64:    base64.StdEncoding.EncodeToString(bodyBytes),
	}

	payload, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	payload = append(payload, '\n')

	file, err := os.OpenFile("pay.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(payload); err != nil {
		return err
	}
	return nil
}

func writeWechatNotifyErrorFile(stage string, err error) error {
	record := struct {
		Time  string `json:"time"`
		Stage string `json:"stage"`
		Error string `json:"error"`
	}{
		Time:  time.Now().Format(time.RFC3339Nano),
		Stage: stage,
		Error: err.Error(),
	}

	payload, marshalErr := json.MarshalIndent(record, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	payload = append(payload, '\n')

	file, openErr := os.OpenFile("pay.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if openErr != nil {
		return openErr
	}
	defer file.Close()

	_, writeErr := file.Write(payload)
	return writeErr
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
