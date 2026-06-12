package object

import (
	"encoding/json"
)

// MarshalJSON 将任意值序列化为 JSON 字符串。
//
// 常用于把结构化快照写入数据库 text/json 字段。
func MarshalJSON(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
