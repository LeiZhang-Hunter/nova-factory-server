// 管家婆全渠道 API 标准响应结构。
// 定义通用错误响应体，所有 API 错误均按此格式返回。
package models

// ErrorResponse 全渠道 API 标准错误响应
// 管家婆 API 约定错误响应的 iserror 为 true，errormsg 包含错误描述
type ErrorResponse struct {
	Iserror  bool   `json:"iserror"`
	Errormsg string `json:"errormsg"`
}
