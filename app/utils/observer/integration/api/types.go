// 定义集成 API 层的基础数据类型。
// LoginState 接口用于表示第三方系统的授权登录状态，
// 前端根据返回的 CheckURL 跳转至 OAuth 授权页面完成登录。
package api

// LoginState 登录状态接口，描述第三方集成系统的授权连接状态。
// 当第三方系统未授权时，前端根据 GetCheckURL 返回的地址跳转至 OAuth 页面。
type LoginState interface {
	// GetOnline 返回当前是否已在线/已授权
	GetOnline() bool
	// GetMessage 返回状态描述消息（如 "管家婆授权页面"）
	GetMessage() string
	// GetType 返回集成系统类型标识（如 "gjp_v1"）
	GetType() string
	// GetCheckURL 返回 OAuth 授权跳转地址，未授权时前端跳转到此 URL
	GetCheckURL() string
	// GetRaw 返回原始响应字符串，供调试使用
	GetRaw() string
}
