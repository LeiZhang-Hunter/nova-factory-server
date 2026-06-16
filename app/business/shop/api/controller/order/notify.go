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

	"go.uber.org/zap"

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
		writeWechatNotifyFail(c, "读取请求体失败", zap.Error(err))
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if err := writeWechatNotifyRequestFile(c, bodyBytes); err != nil {
		writeWechatNotifyFail(c, "写入请求数据失败", zap.Error(err), zap.Int("body_bytes", len(bodyBytes)))
		return
	}
	// 读微信配置
	cfgMap, err := s.loadNotifyConfig(c)
	if err != nil {
		writeWechatNotifyFail(c, "配置读取失败", zap.Error(err))
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
		writeWechatNotifyFail(c, "私钥读取失败", zap.Error(err), zap.String("private_key_path", privateKeyPath))
		return
	}
	client, err := gopayWechat.NewClientV3(mchId, serialNo, apiV3Key, string(privateKeyData))
	if err != nil {
		writeWechatNotifyFail(c, "初始化微信客户端失败", zap.Error(err), zap.String("mch_id", mchId), zap.String("serial_no", serialNo))
		return
	}
	platformPublicKeyData, err := file.ReadPrivateFile(c, platformPublicKeyPath)
	if err != nil {
		writeWechatNotifyFail(c, "微信支付公钥读取失败", zap.Error(err), zap.String("platform_public_key_path", platformPublicKeyPath))
		return
	}
	if err := client.AutoVerifySignByPublicKey(platformPublicKeyData, platformPublicKeyID); err != nil {
		if fileErr := writeWechatNotifyErrorFile("auto_verify_sign_by_public_key", err); fileErr != nil {
			zap.L().Error("write wechat notify error file failed", zap.Error(fileErr))
		}
		writeWechatNotifyFail(c, "验签初始化失败", zap.Error(err), zap.String("platform_public_key_id", platformPublicKeyID))
		return
	}

	// 解析回调
	notifyReq, err := gopayWechat.V3ParseNotify(c.Request)
	if err != nil {
		writeWechatNotifyFail(c, "解析失败", zap.Error(err))
		return
	}
	// 获取微信平台证书
	certMap := client.WxPublicKeyMap()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		if fileErr := writeWechatNotifyErrorFile("verify_sign", err); fileErr != nil {
			zap.L().Error("write wechat notify error file failed", zap.Error(fileErr))
		}
		writeWechatNotifyFail(c, "验签失败", zap.Error(err))
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
		writeWechatNotifyFail(c, "解密失败", zap.Error(err))
		return
	}

	// 校验
	if nd.TradeState != "SUCCESS" {
		zap.L().Warn("wechat pay notify ignored", zap.String("trade_state", nd.TradeState), zap.String("out_trade_no", nd.OutTradeNo), zap.String("transaction_id", nd.TransactionID))
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "SUCCESS", Message: "非支付成功通知"})
		return
	}
	if nd.AppID != appId {
		writeWechatNotifyFail(c, "appid不匹配", zap.String("expected_appid", appId), zap.String("actual_appid", nd.AppID), zap.String("out_trade_no", nd.OutTradeNo))
		return
	}
	if nd.MchID != mchId {
		writeWechatNotifyFail(c, "mchid不匹配", zap.String("expected_mch_id", mchId), zap.String("actual_mch_id", nd.MchID), zap.String("out_trade_no", nd.OutTradeNo))
		return
	}

	// 序列化回调原文（排障用）
	rawBytes, _ := json.Marshal(nd)
	notifyRaw := string(rawBytes)

	// 更新订单状态
	if err := s.service.HandleWechatNotify(c, nd.OutTradeNo, nd.TransactionID, notifyRaw, nd.MchID, nd.AppID, nd.Payer.Openid, nd.Amount.Total); err != nil {
		writeWechatNotifyFail(c, err.Error(), zap.Error(err), zap.String("out_trade_no", nd.OutTradeNo), zap.String("transaction_id", nd.TransactionID), zap.Int64("notify_total", nd.Amount.Total))
		return
	}

	c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "SUCCESS", Message: "成功"})
}

func writeWechatNotifyFail(c *gin.Context, message string, fields ...zap.Field) {
	baseFields := []zap.Field{
		zap.String("message", message),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("client_ip", c.ClientIP()),
	}
	zap.L().Error("wechat pay notify failed", append(baseFields, fields...)...)
	c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: message})
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
