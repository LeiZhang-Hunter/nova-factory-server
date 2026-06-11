package qqd

import (
	"net/http"
	"strings"

	qqdservice "nova-factory-server/app/business/erp_api/service/qqd"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type API struct {
	baseController
}

func NewAPI(service qqdservice.Service) *API {
	return &API{baseController: baseController{service: service}}
}

func (q *API) API(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		zap.L().Error("parse qqd api form failed", zap.Error(err))
	}
	method := qqdParam(c, "method")
	appKey := qqdParam(c, "app_key")
	accessToken := qqdParam(c, "access_token")

	if appKey == "" {
		appKey = qqdParam(c, "appkey")
	}
	if accessToken == "" {
		accessToken = qqdParam(c, "token")
	}
	if method == "" {
		c.JSON(http.StatusOK, qqdError("method is required"))
		return
	}
	if !q.service.ValidAccessToken(c, accessToken, appKey) {
		c.JSON(http.StatusOK, qqdError(qqdservice.ErrInvalidAccessToken.Error()))
		return
	}

	body, err := readAndRestoreRequestBody(c)
	if err != nil {
		zap.L().Error("read qqd sign body failed", zap.Error(err))
		c.JSON(http.StatusOK, qqdError(qqdservice.ErrInvalidSign.Error()))
		return
	}
	if !q.service.ValidSign(formValues(c), body, qqdParam(c, "sign")) {
		c.JSON(http.StatusOK, qqdError(qqdservice.ErrInvalidSign.Error()))
		return
	}

	switch strings.ToLower(method) {
	case "selfmall.product.list.get":
		response, err := q.service.ProductList(c, qqdservice.ProductListRequest{
			PageNo:   parseIntParam(c, "pageno"),
			PageSize: parseIntParam(c, "pagesize"),
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
		c.JSON(http.StatusOK, qqdError("method not implemented: "+method))
	default:
		c.JSON(http.StatusOK, qqdError("unsupported method: "+method))
	}
}

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

func (q *API) productAdd(c *gin.Context) {
	goodsInfos, err := parseProductAddGoodsInfos(c)
	if err != nil {
		c.JSON(http.StatusOK, qqdError(err.Error()))
		return
	}

	storedGoods, err := q.service.AddProducts(c, goodsInfos)
	if err != nil {
		zap.L().Error("save qqd product add payload failed", zap.Error(err))
		c.JSON(http.StatusOK, qqdError("save product failed: "+err.Error()))
		return
	}

	result := make([]gin.H, 0, len(goodsInfos))
	for _, goods := range goodsInfos {
		result = append(result, gin.H{
			"goodsid":  goods["goodsid"],
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
