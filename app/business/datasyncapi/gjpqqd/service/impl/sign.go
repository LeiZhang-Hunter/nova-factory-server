// 管家婆全渠道 API 的 MD5 签名工具。
// 实现管家婆开放平台要求的签名算法：
// 将所有非 sign 参数按 key 字母序排列，
// 以 appSecret + 排序后的 keyvalue 对 + body + appSecret 拼接后计算 MD5。
package impl

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// GenerateMD5Sign 按管家婆开放接口的 MD5 签名规则生成签名
// 规则：
// 1. 排除 params 中 key 为 sign 的项
// 2. 剩余 key 按字母序升序排列
// 3. 拼接：appSecret + key1value1 + key2value2 + ... + body + appSecret
// 4. 计算 MD5 并转为大写 hex 字符串
func GenerateMD5Sign(params map[string]string, body string, appSecret string) (string, error) {
	if appSecret == "" {
		return "", fmt.Errorf("appSecret is required")
	}

	// 收集并排序参数 key
	keys := make([]string, 0, len(params))
	for key := range params {
		if strings.EqualFold(key, "sign") {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 按规则拼接签名字符串
	var builder strings.Builder
	builder.Grow(len(appSecret)*2 + len(body) + len(keys)*16)
	builder.WriteString(appSecret)
	for _, key := range keys {
		builder.WriteString(key)
		builder.WriteString(params[key])
	}
	builder.WriteString(body)
	builder.WriteString(appSecret)

	// MD5 哈希并转大写
	sum := md5.Sum([]byte(builder.String()))
	return strings.ToUpper(hex.EncodeToString(sum[:])), nil
}
