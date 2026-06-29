package provider

import (
	"nova-factory-server/app/business/shop/order/models"

	"github.com/gin-gonic/gin"
)

// PaymentMethod 统一支付方法接口，每种支付通道实现一个实例。
type PaymentMethod interface {
	Prepay(c *gin.Context, req *models.PrepayRequest) (models.PrepayResult, error)
	Refund(c *gin.Context, req *models.RefundRequest) (*models.RefundResult, error)

	// ParsePayNotify 解析并验证支付回调通知。
	// 负责：通道级验签、解密、appid/mchid 校验。对合法但 TradeState != SUCCESS 的通知正常返回数据。
	ParsePayNotify(c *gin.Context) (models.PayNotifyDataInterface, error)

	// ParseRefundNotify 解析并验证退款回调通知。同上，RefundStatus != SUCCESS 不视为解析错误。
	ParseRefundNotify(c *gin.Context) (models.RefundNotifyDataInterface, error)

	// QueryOrderByOutTradeNo 通过商户订单号查询微信侧订单状态。
	QueryOrderByOutTradeNo(c *gin.Context, outTradeNo string) (*models.QueryOrderResult, error)

	// QueryRefundByOutRefundNo 通过商户退款单号查询微信侧退款状态。
	QueryRefundByOutRefundNo(c *gin.Context, outRefundNo string) (*models.QueryRefundResult, error)

	Channel() int
}
