package qqd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	qqdservice "nova-factory-server/app/business/erp_api/service/qqd"

	"github.com/gin-gonic/gin"
)

func qqdError(message string) gin.H {
	return gin.H{
		"iserror":  true,
		"errormsg": message,
	}
}

func qqdParam(c *gin.Context, key string) string {
	if value := c.PostForm(key); value != "" {
		return value
	}
	return c.Query(key)
}

func parseIntParam(c *gin.Context, key string) int {
	value := qqdParam(c, key)
	if value == "" {
		return 0
	}
	result, _ := strconv.Atoi(value)
	return result
}

func readAndRestoreRequestBody(c *gin.Context) (string, error) {
	if c.Request.Body == nil {
		return "", nil
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "", err
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	return string(body), nil
}

func formValues(c *gin.Context) map[string]string {
	params := make(map[string]string)
	for key, values := range c.Request.Form {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	return params
}

func parseProductAddGoodsInfos(c *gin.Context) ([]map[string]any, error) {
	if rawGoodsInfos := c.PostForm("goodsinfos"); rawGoodsInfos != "" {
		var goodsInfos []map[string]any
		if err := json.Unmarshal([]byte(rawGoodsInfos), &goodsInfos); err != nil {
			return nil, fmt.Errorf("invalid goodsinfos: %w", err)
		}
		if len(goodsInfos) == 0 {
			return nil, errors.New("goodsinfos is empty")
		}
		return goodsInfos, nil
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("read request body: %w", err)
	}
	if len(strings.TrimSpace(string(body))) == 0 {
		return nil, errors.New("goodsinfos is required")
	}

	var payload struct {
		GoodsInfos []map[string]any `json:"goodsinfos"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}
	if len(payload.GoodsInfos) == 0 {
		return nil, errors.New("goodsinfos is empty")
	}
	return payload.GoodsInfos, nil
}

func parseProductStockUpdateRequest(c *gin.Context) (qqdservice.ProductStockUpdateRequest, error) {
	request := qqdservice.ProductStockUpdateRequest{
		ProductID:  qqdParam(c, "productid"),
		ProductQty: qqdParam(c, "productqty"),
		Skus:       qqdParam(c, "skus"),
	}
	if request.ProductID != "" || request.ProductQty != "" || request.Skus != "" {
		return request, nil
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return request, fmt.Errorf("read request body: %w", err)
	}
	if len(strings.TrimSpace(string(body))) == 0 {
		return request, nil
	}

	var payload struct {
		ProductID  string `json:"productid"`
		ProductQty any    `json:"productqty"`
		Skus       string `json:"skus"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return request, fmt.Errorf("invalid request body: %w", err)
	}
	request.ProductID = payload.ProductID
	request.ProductQty = toString(payload.ProductQty)
	request.Skus = payload.Skus
	return request, nil
}

func toString(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}
