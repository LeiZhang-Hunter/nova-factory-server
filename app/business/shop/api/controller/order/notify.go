package order

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/order/callback"
	orderDao "nova-factory-server/app/business/shop/order/dao"
	models2 "nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/constant/order"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/datasource/objectFile"
	"nova-factory-server/app/utils/observer/integration/observer"
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/service"

	"github.com/gin-gonic/gin"
	gopayWechat "github.com/go-pay/gopay/wechat/v3"
)

// OrderNotify 微信支付回调控制器
type OrderNotify struct {
	service   service.IApiShopOrderService
	configDao dao.IApiShopSysConfigDao
	orderDao  orderDao.IOrderDao
	db        *gorm.DB
	cache     cache.Cache
}

// NewOrderNotify 创建微信支付回调控制器。
func NewOrderNotify(service service.IApiShopOrderService, configDao dao.IApiShopSysConfigDao,
	orderDao orderDao.IOrderDao, db *gorm.DB, cache cache.Cache) *OrderNotify {
	return &OrderNotify{service: service,
		configDao: configDao, db: db,
		orderDao: orderDao,
		cache:    cache,
	}
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
	// 读取微信配置
	config, err := s.configDao.GetWechatPayConfig(c)
	if err != nil {
		zap.L().Error("读取微信配置失败", zap.Error(err))
		writeWechatNotifyFail(c, "读取请求体失败", zap.Error(err))
		return
	}
	if config.AppId == "" || config.MchId == "" || config.ApiV3Key == "" || config.SerialNo == "" || config.PrivateKeyPath == "" || config.NotifyUrl == "" {
		writeWechatNotifyFail(c, "微信支付配置不完整，请在后台管理配置微信支付参数")
		return
	}
	file := objectFile.NewConfig()
	privateKeyData, err := file.ReadPrivateFile(c, config.PrivateKeyPath)
	if err != nil {
		writeWechatNotifyFail(c, "私钥读取失败", zap.Error(err), zap.String("private_key_path", config.PrivateKeyPath))
		return
	}
	client, err := gopayWechat.NewClientV3(config.MchId, config.SerialNo, config.ApiV3Key, string(privateKeyData))
	if err != nil {
		writeWechatNotifyFail(c, "初始化微信客户端失败", zap.Error(err), zap.String("mch_id", config.MchId), zap.String("serial_no", config.SerialNo))
		return
	}
	platformPublicKeyData, err := file.ReadPrivateFile(c, config.PlatformPublicKeyPath)
	if err != nil {
		writeWechatNotifyFail(c, "微信支付公钥读取失败", zap.Error(err), zap.String("platform_public_key_path", config.PlatformPublicKeyPath))
		return
	}
	if err := client.AutoVerifySignByPublicKey(platformPublicKeyData, config.PlatformPublicKeyId); err != nil {
		if fileErr := writeWechatNotifyErrorFile("auto_verify_sign_by_public_key", err); fileErr != nil {
			zap.L().Error("write wechat notify error file failed", zap.Error(fileErr))
		}
		writeWechatNotifyFail(c, "验签初始化失败", zap.Error(err), zap.String("platform_public_key_id", config.PlatformPublicKeyId))
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
	var nd models.WechatPayNotifyData
	if err := notifyReq.DecryptCipherTextToStruct(config.ApiV3Key, &nd); err != nil {
		writeWechatNotifyFail(c, "解密失败", zap.Error(err))
		return
	}

	// 校验
	if nd.TradeState != "SUCCESS" {
		zap.L().Warn("wechat pay notify ignored", zap.String("trade_state", nd.TradeState), zap.String("out_trade_no", nd.OutTradeNo), zap.String("transaction_id", nd.TransactionID))
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "SUCCESS", Message: "非支付成功通知"})
		return
	}
	if nd.AppID != config.AppId {
		writeWechatNotifyFail(c, "appid不匹配", zap.String("expected_appid", config.AppId), zap.String("actual_appid", nd.AppID), zap.String("out_trade_no", nd.OutTradeNo))
		return
	}
	if nd.MchID != config.MchId {
		writeWechatNotifyFail(c, "mchid不匹配", zap.String("expected_mch_id", config.MchId), zap.String("actual_mch_id", nd.MchID), zap.String("out_trade_no", nd.OutTradeNo))
		return
	}

	//data := models.NewOrderStatusData(nd.OutTradeNo, order.ERPStatusPayed, order.REFUNDStatusNormal)
	//m.WithOrders(data)
	//m.WithMetadata(gin.H{"order": nd})
	//m.WithCtx(c)
	//m.WithDB(s.db)

	info, err := s.orderDao.GetByTid(c, nd.OutTradeNo)
	if err != nil {
		zap.L().Error("获取订单失败", zap.Error(err), zap.String("tid", nd.OutTradeNo))
		writeWechatNotifyFail(c, "获取订单失败", zap.Error(err), zap.String("tid", nd.OutTradeNo))
		return
	}
	if info == nil {
		zap.L().Error("订单不存在", zap.String("tid", nd.OutTradeNo))
		writeWechatNotifyFail(c, "订单不存在", zap.String("tid", nd.OutTradeNo))
		return
	}
	info.Status = order.ERPStatusPayed
	now := time.Now()
	info.PayTime = &now
	// 序列化回调原文（排障用）
	var request models2.OrderSyncRequest
	request.Orders = []*models2.OrderSyncOrder{
		models2.ToOrderSyncOrder(info, &nd),
	}
	request.WithCallback(callback.NewOrderSyncRequestCallback(c, &request))
	request.WithDB(s.db)
	request.WithCache(s.cache)
	//m.WithCallback(callback.NewShopApiCallback(m))
	err = observer.GetNotifier().OnOrderChanged(&request)
	if err != nil {
		writeWechatNotifyFail(c, "订单同步观察者失败", zap.Error(err))
		return
	}
	// 更新订单状态
	//if err := s.service.HandleWechatNotify(c, nd.OutTradeNo, nd.TransactionID, notifyRaw, nd.MchID, nd.AppID, nd.Payer.Openid, nd.Amount.Total); err != nil {
	//	writeWechatNotifyFail(c, err.Error(), zap.Error(err), zap.String("out_trade_no", nd.OutTradeNo), zap.String("transaction_id", nd.TransactionID), zap.Int64("notify_total", nd.Amount.Total))
	//	return
	//}

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
