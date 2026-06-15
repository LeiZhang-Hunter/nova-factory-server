// 管家婆全渠道控制器辅助工具函数。
// 提供统一错误响应构造、参数解析、请求体读取与商品数据解析等功能。
package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
)

// qqdError 构造管家婆 API 风格的标准错误响应
// 格式: {"iserror": true, "errormsg": "错误信息"}
func qqdError(message string) gin.H {
	return gin.H{
		"iserror":  true,
		"errormsg": message,
	}
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
