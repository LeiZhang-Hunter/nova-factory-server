// 管家婆全渠道控制器辅助工具函数。
// 提供统一错误响应构造、参数解析、请求体读取与商品数据解析等功能。
package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// qqdError 构造管家婆 API 风格的标准错误响应
// 格式: {"iserror": true, "errormsg": "错误信息"}
func qqdError(message string) gin.H {
	return gin.H{
		"iserror":  true,
		"errormsg": message,
	}
}

// qqdParam 从请求中获取参数值，优先 PostForm，其次 Query
func qqdParam(c *gin.Context, key string) string {
	if value := c.PostForm(key); value != "" {
		return value
	}
	return c.Query(key)
}

// readAndRestoreRequestBody 读取请求 Body 并重新填充，支持多次读取
// 用于在签名校验前读取原始 Body 内容
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

// formValues 将请求的 form 参数转换为 map（每个 key 取第一个值）
func formValues(c *gin.Context) map[string]string {
	params := make(map[string]string)
	for key, values := range c.Request.Form {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	return params
}

// parseProductAddGoodsInfos 解析商品新增请求中的 goodsinfos 数据
// 优先从 PostForm("goodsinfos") 解析 JSON 字符串，
// 若为空则从请求 Body 的 JSON 中按 goodsinfos 字段解析
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
