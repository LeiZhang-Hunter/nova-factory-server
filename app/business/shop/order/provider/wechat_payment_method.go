package provider

import (
	"fmt"
	"math"
	"time"

	apiDao "nova-factory-server/app/business/shop/api/dao"
	apimodels "nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/order/models"
	orderConstant "nova-factory-server/app/constant/order"
	"nova-factory-server/app/datasource/objectFile"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"go.uber.org/zap"
)

// WechatPaymentMethod 微信 V3 支付实现（预支付 + 退款）。
type WechatPaymentMethod struct {
	configDao apiDao.IApiShopSysConfigDao
}

// NewWechatPaymentMethod 创建微信支付方法。
func NewWechatPaymentMethod(configDao apiDao.IApiShopSysConfigDao) PaymentMethod {
	p := &WechatPaymentMethod{
		configDao: configDao,
	}
	RegisterPaymentMethod(p)
	return p
}

func (p *WechatPaymentMethod) Channel() int {
	return orderConstant.PayChannelWechat
}

// Prepay 微信 V3 JSAPI 预下单。
func (p *WechatPaymentMethod) Prepay(c *gin.Context, req *models.PrepayRequest) (models.PrepayResult, error) {
	cfg, err := p.configDao.GetWechatPayConfig(c)
	if err != nil {
		return nil, fmt.Errorf("读取微信配置失败: %v", err)
	}
	if cfg.AppId == "" || cfg.MchId == "" || cfg.ApiV3Key == "" || cfg.SerialNo == "" || cfg.PrivateKeyPath == "" || cfg.NotifyUrl == "" {
		return nil, fmt.Errorf("微信支付配置不完整，请在后台管理配置微信支付参数")
	}

	file := objectFile.NewConfig()
	privateKeyData, err := file.ReadPrivateFile(c, cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取微信支付私钥文件失败: %v", err)
	}

	client, err := wechat.NewClientV3(cfg.MchId, cfg.SerialNo, cfg.ApiV3Key, string(privateKeyData))
	if err != nil {
		return nil, fmt.Errorf("初始化微信支付客户端失败: %v", err)
	}

	total := int64(math.Round(req.TotalAmount * 100))
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)

	bm := make(gopay.BodyMap)
	bm.Set("appid", cfg.AppId).
		Set("mchid", cfg.MchId).
		Set("description", req.Description).
		Set("out_trade_no", req.Tid).
		Set("time_expire", expire).
		Set("notify_url", cfg.NotifyUrl+"/api/v1/app/shop/wechat/order/notify").
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", total).
				Set("currency", "CNY")
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", req.Openid)
		})

	wxRsp, err := client.V3TransactionJsapi(c.Request.Context(), bm)
	if err != nil {
		return nil, fmt.Errorf("微信预下单失败: %v", err)
	}
	if wxRsp.Code != 0 {
		return nil, fmt.Errorf("微信预下单失败: %s", wxRsp.Error)
	}

	paySign, err := client.PaySignOfApplet(cfg.AppId, wxRsp.Response.PrepayId)
	if err != nil {
		return nil, fmt.Errorf("生成支付签名失败: %v", err)
	}

	return &models.WechatPrepayResult{
		PrepayID:  wxRsp.Response.PrepayId,
		AppId:     paySign.AppId,
		TimeStamp: paySign.TimeStamp,
		NonceStr:  paySign.NonceStr,
		Package:   "prepay_id=" + wxRsp.Response.PrepayId,
		SignType:  "RSA",
		PaySign:   paySign.PaySign,
	}, nil
}

// Refund 微信 V3 退款。
func (p *WechatPaymentMethod) Refund(c *gin.Context, req *models.RefundRequest) (*models.RefundResult, error) {
	cfg, err := p.configDao.GetWechatPayConfig(c)
	if err != nil {
		return nil, fmt.Errorf("读取微信配置失败: %v", err)
	}
	if cfg.AppId == "" || cfg.MchId == "" || cfg.ApiV3Key == "" || cfg.SerialNo == "" || cfg.PrivateKeyPath == "" {
		return nil, fmt.Errorf("微信支付配置不完整")
	}

	file := objectFile.NewConfig()
	privateKeyData, err := file.ReadPrivateFile(c, cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取微信支付私钥失败: %v", err)
	}

	client, err := wechat.NewClientV3(cfg.MchId, cfg.SerialNo, cfg.ApiV3Key, string(privateKeyData))
	if err != nil {
		return nil, fmt.Errorf("初始化微信支付客户端失败: %v", err)
	}

	notifyURL := req.NotifyURL
	if notifyURL == "" {
		notifyURL = cfg.NotifyUrl + "/api/v1/app/shop/wechat/order/refund/notify"
	}

	refundCents := int64(req.RefundAmount * 100)
	totalCents := int64(req.TotalAmount * 100)

	bm := make(gopay.BodyMap)
	bm.Set("transaction_id", req.TransactionID).
		Set("out_trade_no", req.Tid).
		Set("out_refund_no", req.OutRefundNo).
		Set("reason", req.Reason).
		Set("notify_url", notifyURL).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("refund", refundCents).
				Set("total", totalCents).
				Set("currency", "CNY")
		})

	wxRsp, err := client.V3Refund(c, bm)
	if err != nil {
		return nil, fmt.Errorf("微信退款申请失败: %v", err)
	}
	if wxRsp.Code != 0 {
		return nil, fmt.Errorf("微信退款申请失败: %s", wxRsp.Error)
	}

	zap.L().Info("微信退款申请成功",
		zap.String("out_refund_no", req.OutRefundNo),
		zap.String("out_trade_no", req.Tid),
	)

	return &models.RefundResult{
		ThirdRefundID:      wxRsp.Response.RefundId,
		ThirdTransactionID: req.TransactionID,
		Status:             wxRsp.Response.Status,
	}, nil
}

// ParsePayNotify 解析并验证微信支付回调通知。
// 负责：验签、解密、appid/mchid 校验。TradeState != SUCCESS 不视为错误。
func (p *WechatPaymentMethod) ParsePayNotify(c *gin.Context) (models.PayNotifyDataInterface, error) {
	cfg, err := p.configDao.GetWechatPayConfig(c)
	if err != nil {
		return nil, fmt.Errorf("读取微信配置失败: %v", err)
	}
	if cfg.AppId == "" || cfg.MchId == "" || cfg.ApiV3Key == "" || cfg.SerialNo == "" || cfg.PrivateKeyPath == "" {
		return nil, fmt.Errorf("微信支付配置不完整")
	}

	file := objectFile.NewConfig()
	privateKeyData, err := file.ReadPrivateFile(c, cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取微信支付私钥失败: %v", err)
	}

	client, err := wechat.NewClientV3(cfg.MchId, cfg.SerialNo, cfg.ApiV3Key, string(privateKeyData))
	if err != nil {
		return nil, fmt.Errorf("初始化微信支付客户端失败: %v", err)
	}

	platformPublicKeyData, err := file.ReadPrivateFile(c, cfg.PlatformPublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("微信支付公钥读取失败: %v", err)
	}
	if err := client.AutoVerifySignByPublicKey(platformPublicKeyData, cfg.PlatformPublicKeyId); err != nil {
		return nil, fmt.Errorf("验签初始化失败: %v", err)
	}

	notifyReq, err := wechat.V3ParseNotify(c.Request)
	if err != nil {
		return nil, fmt.Errorf("解析回调通知失败: %v", err)
	}
	certMap := client.WxPublicKeyMap()
	if err := notifyReq.VerifySignByPKMap(certMap); err != nil {
		return nil, fmt.Errorf("验签失败: %v", err)
	}

	var nd apimodels.WechatPayNotifyData
	if err := notifyReq.DecryptCipherTextToStruct(cfg.ApiV3Key, &nd); err != nil {
		return nil, fmt.Errorf("解密回调内容失败: %v", err)
	}

	if nd.GetAppID() != cfg.AppId {
		return nil, fmt.Errorf("appid不匹配: expected %s, got %s", cfg.AppId, nd.GetAppID())
	}
	if nd.GetMchID() != cfg.MchId {
		return nil, fmt.Errorf("mchid不匹配: expected %s, got %s", cfg.MchId, nd.GetMchID())
	}

	return &nd, nil
}

// ParseRefundNotify 解析并验证微信退款回调通知。
func (p *WechatPaymentMethod) ParseRefundNotify(c *gin.Context) (models.RefundNotifyDataInterface, error) {
	cfg, err := p.configDao.GetWechatPayConfig(c)
	if err != nil {
		return nil, fmt.Errorf("读取微信配置失败: %v", err)
	}
	if cfg.AppId == "" || cfg.MchId == "" || cfg.ApiV3Key == "" || cfg.SerialNo == "" || cfg.PrivateKeyPath == "" {
		return nil, fmt.Errorf("微信支付配置不完整")
	}

	file := objectFile.NewConfig()
	privateKeyData, err := file.ReadPrivateFile(c, cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取微信支付私钥失败: %v", err)
	}

	client, err := wechat.NewClientV3(cfg.MchId, cfg.SerialNo, cfg.ApiV3Key, string(privateKeyData))
	if err != nil {
		return nil, fmt.Errorf("初始化微信支付客户端失败: %v", err)
	}

	platformPublicKeyData, err := file.ReadPrivateFile(c, cfg.PlatformPublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("微信支付公钥读取失败: %v", err)
	}
	if err := client.AutoVerifySignByPublicKey(platformPublicKeyData, cfg.PlatformPublicKeyId); err != nil {
		return nil, fmt.Errorf("验签初始化失败: %v", err)
	}

	notifyReq, err := wechat.V3ParseNotify(c.Request)
	if err != nil {
		return nil, fmt.Errorf("解析回调通知失败: %v", err)
	}
	certMap := client.WxPublicKeyMap()
	if err := notifyReq.VerifySignByPKMap(certMap); err != nil {
		return nil, fmt.Errorf("验签失败: %v", err)
	}

	var nd apimodels.WechatRefundNotifyData
	if err := notifyReq.DecryptCipherTextToStruct(cfg.ApiV3Key, &nd); err != nil {
		return nil, fmt.Errorf("解密回调内容失败: %v", err)
	}

	return &nd, nil
}

// QueryOrderByOutTradeNo 通过商户订单号查询微信侧订单状态。
func (p *WechatPaymentMethod) QueryOrderByOutTradeNo(c *gin.Context, outTradeNo string) (*models.QueryOrderResult, error) {
	cfg, err := p.configDao.GetWechatPayConfig(c)
	if err != nil {
		return nil, fmt.Errorf("读取微信配置失败: %v", err)
	}
	if cfg.MchId == "" || cfg.ApiV3Key == "" || cfg.SerialNo == "" || cfg.PrivateKeyPath == "" {
		return nil, fmt.Errorf("微信支付配置不完整")
	}

	file := objectFile.NewConfig()
	privateKeyData, err := file.ReadPrivateFile(c, cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取微信支付私钥失败: %v", err)
	}

	client, err := wechat.NewClientV3(cfg.MchId, cfg.SerialNo, cfg.ApiV3Key, string(privateKeyData))
	if err != nil {
		return nil, fmt.Errorf("初始化微信支付客户端失败: %v", err)
	}

	wxRsp, err := client.V3TransactionQueryOrder(c.Request.Context(), wechat.OutTradeNo, outTradeNo)
	if err != nil {
		return nil, fmt.Errorf("微信查询订单失败: %v", err)
	}
	if wxRsp.Code != 0 {
		return nil, fmt.Errorf("微信查询订单失败: %s", wxRsp.Error)
	}

	r := &models.QueryOrderResult{
		Appid:          wxRsp.Response.Appid,
		Mchid:          wxRsp.Response.Mchid,
		OutTradeNo:     wxRsp.Response.OutTradeNo,
		TransactionId:  wxRsp.Response.TransactionId,
		TradeType:      wxRsp.Response.TradeType,
		TradeState:     wxRsp.Response.TradeState,
		TradeStateDesc: wxRsp.Response.TradeStateDesc,
		BankType:       wxRsp.Response.BankType,
		Attach:         wxRsp.Response.Attach,
		SuccessTime:    wxRsp.Response.SuccessTime,
	}
	if wxRsp.Response.Payer != nil {
		r.PayerOpenid = wxRsp.Response.Payer.Openid
	}
	if wxRsp.Response.Amount != nil {
		r.AmountTotal = wxRsp.Response.Amount.Total
		r.AmountPayerTotal = wxRsp.Response.Amount.PayerTotal
		r.AmountCurrency = wxRsp.Response.Amount.Currency
	}
	return r, nil
}

// QueryRefundByOutRefundNo 通过商户退款单号查询微信侧退款状态。
func (p *WechatPaymentMethod) QueryRefundByOutRefundNo(c *gin.Context, outRefundNo string) (*models.QueryRefundResult, error) {
	cfg, err := p.configDao.GetWechatPayConfig(c)
	if err != nil {
		return nil, fmt.Errorf("读取微信配置失败: %v", err)
	}
	if cfg.MchId == "" || cfg.ApiV3Key == "" || cfg.SerialNo == "" || cfg.PrivateKeyPath == "" {
		return nil, fmt.Errorf("微信支付配置不完整")
	}

	file := objectFile.NewConfig()
	privateKeyData, err := file.ReadPrivateFile(c, cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取微信支付私钥失败: %v", err)
	}

	client, err := wechat.NewClientV3(cfg.MchId, cfg.SerialNo, cfg.ApiV3Key, string(privateKeyData))
	if err != nil {
		return nil, fmt.Errorf("初始化微信支付客户端失败: %v", err)
	}

	wxRsp, err := client.V3RefundQuery(c.Request.Context(), outRefundNo, nil)
	if err != nil {
		return nil, fmt.Errorf("微信查询退款失败: %v", err)
	}
	if wxRsp.Code != 0 {
		return nil, fmt.Errorf("微信查询退款失败: %s", wxRsp.Error)
	}

	r := &models.QueryRefundResult{
		RefundId:            wxRsp.Response.RefundId,
		OutRefundNo:         wxRsp.Response.OutRefundNo,
		TransactionId:       wxRsp.Response.TransactionId,
		OutTradeNo:          wxRsp.Response.OutTradeNo,
		Channel:             wxRsp.Response.Channel,
		UserReceivedAccount: wxRsp.Response.UserReceivedAccount,
		SuccessTime:         wxRsp.Response.SuccessTime,
		CreateTime:          wxRsp.Response.CreateTime,
		Status:              wxRsp.Response.Status,
		FundsAccount:        wxRsp.Response.FundsAccount,
	}
	if wxRsp.Response.Amount != nil {
		r.AmountTotal = wxRsp.Response.Amount.Total
		r.AmountRefund = wxRsp.Response.Amount.Refund
		r.AmountPayerTotal = wxRsp.Response.Amount.PayerTotal
		r.AmountPayerRefund = wxRsp.Response.Amount.PayerRefund
		r.AmountCurrency = wxRsp.Response.Amount.Currency
	}
	return r, nil
}
