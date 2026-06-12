package controller

import (
	"net/http"
	"nova-factory-server/app/business/datasyncapi/gjpqqd/models"
	"nova-factory-server/app/business/datasyncapi/gjpqqd/service"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/observer/integration/observer"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// API 全渠道 API 调用控制器，负责鉴权、签名校验及方法分发
// 管家婆全渠道 API 调用控制器。
// 处理 POST /api 请求，通过 method 参数分发到不同子功能，
// 包括商品列表查询、商品新增、库存更新等。
// 调用前先校验 access_token 和签名。
type API struct {
	service service.GjpQqdService
}

// NewAPI 创建 API 控制器实例
func NewAPI(service service.GjpQqdService) *API {
	return &API{
		service: service,
	}
}

// API 处理全渠道 API 的统一入口
// 流程：解析请求参数 -> 校验 access_token -> 校验签名 -> 按 method 分发
func (q *API) API(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		zap.L().Error("parse qqd api form failed", zap.Error(err))
	}

	req := new(models.ApiReq)
	err := c.ShouldBind(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	if req.Method == "" {
		c.JSON(http.StatusOK, qqdError("method is required"))
		return
	}
	// 校验 access_token 是否有效
	if !q.service.ValidAccessToken(c, req.AccessToken, req.AppKey) {
		c.JSON(http.StatusOK, qqdError(models.ErrInvalidAccessToken.Error()))
		return
	}

	// 读取原始请求 Body 用于签名校验
	body, err := readAndRestoreRequestBody(c)
	if err != nil {
		zap.L().Error("read qqd sign body failed", zap.Error(err))
		c.JSON(http.StatusOK, qqdError(models.ErrInvalidSign.Error()))
		return
	}
	// MD5 签名校验
	if !q.service.ValidSign(formValues(c), body, req.Sign) {
		c.JSON(http.StatusOK, qqdError(models.ErrInvalidSign.Error()))
		return
	}

	// 按 method 分发到不同业务处理逻辑
	switch strings.ToLower(req.Method) {
	case "selfmall.product.list.get":
		response, err := q.service.ProductList(c, service.ProductListRequest{
			PageNo:   req.PageNo,
			PageSize: req.PageSize,
		})
		if err != nil {
			zap.L().Error("load qqd product list failed", zap.Error(err))
			c.JSON(http.StatusOK, qqdError("load product list failed: "+err.Error()))
			return
		}
		c.JSON(http.StatusOK, response)
	case "selfmall.product.add":
		q.productAdd(c)
	case "selfmall.productstock.list.update":
		q.productStockUpdate(c)
	case "selfmall.product.query",
		"selfmall.sellercats.list.get",
		"selfmall.order.ship",
		"selfmall.stock.update",
		"selfmall.sale.status.write":
		c.JSON(http.StatusOK, qqdError("method not implemented: "+req.Method))
	default:
		c.JSON(http.StatusOK, qqdError("unsupported method: "+req.Method))
	}
}

// productStockUpdate 处理 selfmall.productstock.list.update 请求
// 解析请求中的 productid、productqty 和 skus，调用服务层执行库存更新
func (q *API) productStockUpdate(c *gin.Context) {
	request, err := parseProductStockUpdateRequest(c)
	if err != nil {
		c.JSON(http.StatusOK, qqdError(err.Error()))
		return
	}

	response, err := q.service.ProductStockUpdate(c, request)
	if err != nil {
		c.JSON(http.StatusOK, qqdError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response)
}

// productAdd 处理 selfmall.product.add 请求
// 将商品同步请求绑定为 GoodsSyncReq，通过 Observer 观察者模式
// 分发给所有注册的集成系统（如管家婆、金蝶）执行商品同步
func (q *API) productAdd(c *gin.Context) {
	goodsInfos := new(models.GoodsSyncReq)
	err := c.ShouldBind(goodsInfos)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	if len(goodsInfos.GoodsInfos) == 0 {
		baizeContext.Waring(c, "goodsinfo is required")
		return
	}

	// 通过全局 Notifier 分发商品变更事件，所有已注册的 Observer 均会收到通知
	err = observer.GetNotifier().OnProductChanged(goodsInfos)

	if err != nil {
		zap.L().Error("save qqd product add payload failed", zap.Error(err))
		c.JSON(http.StatusOK, qqdError("save product failed: "+err.Error()))
		return
	}

	result := make([]gin.H, 0, len(goodsInfos.GoodsInfos))
	for _, goods := range goodsInfos.GoodsInfos {
		result = append(result, gin.H{
			"goodsid":  goods.Goodsid,
			"errormsg": "",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"iserror":  false,
		"errormsg": "ok",
		"result":   result,
		"total":    len(storedGoods),
	})
}
