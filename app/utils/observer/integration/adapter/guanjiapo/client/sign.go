package client

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// GenerateMD5Sign 按管家婆开放接口的 MD5 规则生成签名
func GenerateMD5Sign(params map[string]string, body string, appSecret string) (string, error) {
	if appSecret == "" {
		return "", fmt.Errorf("appSecret is required")
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		if strings.EqualFold(key, "sign") {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	builder.Grow(len(appSecret)*2 + len(body) + len(keys)*16)
	builder.WriteString(appSecret)
	for _, key := range keys {
		builder.WriteString(key)
		builder.WriteString(params[key])
	}
	builder.WriteString(body)
	builder.WriteString(appSecret)

	sum := md5.Sum([]byte(builder.String()))
	return strings.ToUpper(hex.EncodeToString(sum[:])), nil
}
