package controller

import (
	"net/http"
	"nova-factory-server/app/business/datasyncapi/gjpqqd/callback"
	"nova-factory-server/app/business/datasyncapi/gjpqqd/models"
	"nova-factory-server/app/business/datasyncapi/gjpqqd/service"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/observer/integration/observer"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// API 全渠道 API 调用控制器，负责鉴权、签名校验及方法分发
type API struct {
	service service.GjpQqdService
	db      *gorm.DB
}

// NewAPI 创建 API 控制器实例
func NewAPI(service service.GjpQqdService, db *gorm.DB) *API {
	return &API{
		service: service,
		db:      db,
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

	config, err := q.service.GetConfig(c)
	if err != nil || config == nil {
		c.JSON(http.StatusOK, qqdError("后台没有配置参数"))
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
	if !q.service.ValidSign(formValues(c), body, req.Sign, config) {
		c.JSON(http.StatusOK, qqdError(models.ErrInvalidSign.Error()))
		return
	}

	// 按 method 分发到不同业务处理逻辑
	switch strings.ToLower(req.Method) {
	case "selfmall.product.list.get":
		productReq := new(models.ProductListRequest)
		if err := c.ShouldBind(productReq); err != nil {
			zap.L().Error("bind qqd product list request failed", zap.Error(err))
			c.JSON(http.StatusOK, qqdError(err.Error()))
			return
		}
		productReq.PageSize = req.PageSize
		productReq.PageNo = req.PageNo
		response := q.service.ProductList(c, productReq)
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
func (q *API) productStockUpdate(c *gin.Context) {
	req := new(models.ProductStockUpdateRequest)
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, qqdError(err.Error()))
		return
	}

	if req.ProductID != "" {
		req.ProductID = "0" + req.ProductID
	}

	stockReq := req.ToStockSyncReq()
	stockReq.WithDB(q.db)
	if err := observer.GetNotifier().OnStockChanged(stockReq); err != nil {
		zap.L().Error("stock changed notify failed", zap.Error(err))
		c.JSON(http.StatusOK, qqdError("stock sync failed: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.ErrorResponse{
		Iserror:  false,
		Errormsg: "ok",
	})
}

// productAdd 处理 selfmall.product.add 请求
// 将商品同步请求绑定为 GoodsSyncReq，通过 Observer 观察者模式
// 分发给所有注册的集成系统（如管家婆、金蝶）执行商品同步
func (q *API) productAdd(c *gin.Context) {
	goodsInfos := new(models.GoodsSyncReq)
	err := c.ShouldBind(goodsInfos)
	if err != nil {
		zap.L().Error("should bind error", zap.Error(err))
		c.JSON(http.StatusOK, models.ErrorResponse{
			Iserror:  true,
			Errormsg: err.Error(),
		})
		return
	}
	if len(goodsInfos.GoodsInfos) == 0 {
		c.JSON(http.StatusOK, models.ErrorResponse{
			Iserror:  true,
			Errormsg: "goodsinfo is required",
		})
		return
	}
	result := make([]gin.H, 0, len(goodsInfos.GoodsInfos))
	for _, goods := range goodsInfos.GetProducts() {
		result = append(result, gin.H{
			"goodsid":  goods.GetGoodsId(),
			"errormsg": "",
		})
	}

	call := callback.NewGoodsCallback()
	goodsInfos.WIthCallback(call)
	goodsInfos.WithDB(q.db)
	err = observer.GetNotifier().OnProductChanged(goodsInfos)
	if err != nil {
		c.JSON(http.StatusOK, models.ErrorResponse{
			Iserror:  true,
			Errormsg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"iserror":  false,
		"errormsg": "ok",
		"result":   result,
		"total":    len(result),
	})
}
