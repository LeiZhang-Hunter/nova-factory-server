package order

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	apimodels "nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/order/callback"
	orderDao "nova-factory-server/app/business/shop/order/dao"
	models2 "nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/business/shop/order/provider"
	service2 "nova-factory-server/app/business/shop/order/service"
	orderConstant "nova-factory-server/app/constant/order"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/observer"
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"nova-factory-server/app/business/shop/api/service"

	"github.com/gin-gonic/gin"
	gopayWechat "github.com/go-pay/gopay/wechat/v3"
)

// OrderNotify 支付回调控制器（通道无关，通过工厂分发）。
type OrderNotify struct {
	service        service.IApiShopOrderService
	orderDao       orderDao.IOrderDao
	orderRefundDao orderDao.IOrderRefundDao
	orderService   service2.IOrderService
	db             *gorm.DB
	cache          cache.Cache
}

// NewOrderNotify 创建支付回调控制器。
func NewOrderNotify(service service.IApiShopOrderService,
	orderDao orderDao.IOrderDao, orderService service2.IOrderService,
	aftersaleDao orderDao.IOrderRefundDao, db *gorm.DB, cache cache.Cache) *OrderNotify {
	return &OrderNotify{service: service, db: db,
		orderDao:       orderDao,
		orderRefundDao: aftersaleDao,
		cache:          cache,
		orderService:   orderService,
	}
}

// PublicRoutes 注册微信支付回调路由到 publicGroup（不鉴权）。
func (s *OrderNotify) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop")
	group.Any("/wechat/order/notify", s.HandleWechatNotify)
	group.Any("/wechat/order/refund/notify", s.HandleWechatRefundNotify)

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

	pm, err := provider.GetPaymentMethod(orderConstant.PayChannelWechat)
	if err != nil {
		writeWechatNotifyFail(c, "获取支付通道失败", zap.Error(err))
		return
	}
	nd, err := pm.ParsePayNotify(c)
	if err != nil {
		writeWechatNotifyFail(c, "解析回调失败", zap.Error(err))
		return
	}

	r, err := models2.GetPayNotifyResult[*apimodels.WechatPayNotifyData](nd)
	if err != nil {
		writeWechatNotifyFail(c, "通知类型断言失败", zap.Error(err))
		return
	}

	if r.TradeState != "SUCCESS" {
		zap.L().Warn("wechat pay notify ignored", zap.String("trade_state", r.TradeState), zap.String("out_trade_no", r.OutTradeNo), zap.String("transaction_id", r.TransactionId))
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "SUCCESS", Message: "非支付成功通知"})
		return
	}
	//data := models.NewOrderStatusData(r.OutTradeNo, order.ERPStatusPayed, order.REFUNDStatusNormal)
	//m.WithOrders(data)
	//m.WithMetadata(gin.H{"order": nd})
	//m.WithCtx(c)
	//m.WithDB(s.db)

	info, err := s.orderService.GetByTID(c, r.OutTradeNo)
	if err != nil {
		zap.L().Error("获取订单失败", zap.Error(err), zap.String("tid", r.OutTradeNo))
		writeWechatNotifyFail(c, "获取订单失败", zap.Error(err), zap.String("tid", r.OutTradeNo))
		return
	}
	if info == nil {
		zap.L().Error("订单不存在", zap.String("tid", r.OutTradeNo))
		writeWechatNotifyFail(c, "订单不存在", zap.String("tid", r.OutTradeNo))
		return
	}
	info.Status = orderConstant.ERPStatusPayed
	now := time.Now()
	info.PayTime = &now
	// 序列化回调原文（排障用）
	var request models2.OrderSyncRequest
	request.Orders = []*models2.OrderSyncOrder{
		models2.ToOrderSyncOrder(info, r),
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

// HandleWechatRefundNotify  微信退款异步回调。
func (s *OrderNotify) HandleWechatRefundNotify(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		writeRefundNotifyFail(c, "读取请求体失败")
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	pm, err := provider.GetPaymentMethod(orderConstant.PayChannelWechat)
	if err != nil {
		writeRefundNotifyFail(c, "获取支付通道失败")
		return
	}
	nd, err := pm.ParseRefundNotify(c)
	if err != nil {
		writeRefundNotifyFail(c, "解析回调失败")
		return
	}

	r, err := models2.GetRefundNotifyResult[*apimodels.WechatRefundNotifyData](nd)
	if err != nil {
		writeRefundNotifyFail(c, "通知类型断言失败")
		return
	}

	if r.RefundStatus != "SUCCESS" {
		zap.L().Info("退款回调非成功状态",
			zap.String("out_refund_no", r.OutRefundNo),
			zap.String("refund_status", r.RefundStatus),
		)
		c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "SUCCESS", Message: "非成功状态"})
		return
	}

	refunds, err := s.orderRefundDao.GetByOutRefundNo(c, r.OutRefundNo)
	if err != nil || refunds == nil {
		zap.L().Error("售后单不存在", zap.String("out_refund_no", r.OutRefundNo))
		writeRefundNotifyFail(c, "售后单不存在")
		return
	}
	if refunds.Status == orderConstant.AftersaleStatusRefundSuccess {
		c.JSON(http.StatusOK, nil)
		return
	}

	updates := map[string]any{
		"third_refund_id":      r.RefundId,
		"third_transaction_id": r.OutTradeNo,
	}
	err = s.orderRefundDao.UpdateStatusWithTx(s.db, refunds.ID, orderConstant.AftersaleStatusRefundSuccess, updates)
	if err != nil {
		zap.L().Error("售后单更新失败", zap.String("out_refund_no", r.OutRefundNo), zap.Error(err))
		return
	}

	// 触发管家婆售后同步
	order, _ := s.orderDao.GetByID(c, uint64(refunds.OrderID))
	if order != nil {
		event := models2.NewAftersaleSyncEvent(refunds, order)
		cb := callback.NewAfterSaleSyncCallback(c, s.orderRefundDao, refunds.ID, event)
		event.WithCallback(cb)
		event.WithDB(s.db)

		if err := observer.GetNotifier().OnAfterSaleOrderChanged(event); err != nil {
			zap.L().Error("售后单同步触发失败",
				zap.String("out_refund_no", refunds.OutRefundNo),
				zap.Error(err),
			)
		}
	}
	c.JSON(http.StatusOK, nil)
	return
}

func writeRefundNotifyFail(c *gin.Context, message string) {
	payload, _ := json.Marshal(map[string]string{
		"message": message,
		"path":    c.Request.URL.Path,
	})
	zap.L().Error("refund notify failed", zap.String("detail", string(payload)))
	c.JSON(http.StatusOK, &gopayWechat.V3NotifyRsp{Code: "FAIL", Message: message})
}
