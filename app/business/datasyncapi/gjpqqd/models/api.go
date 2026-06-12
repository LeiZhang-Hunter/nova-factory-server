// 管家婆全渠道 API 通用请求与响应数据结构。
// 定义 API 调用的标准请求参数（method、签名、分页等）及标准错误响应体。
package models

// ApiReq 全渠道 API 统一请求参数
// method 为 API 调用方法名，sign 为 MD5 签名，pageno/pagesize 为分页参数
type ApiReq struct {
	Method      string `json:"method" form:"method"`
	AppKey      string `json:"app_key" form:"app_key"`
	AccessToken string `json:"access_token" form:"access_token"`
	Sign        string `json:"sign" form:"sign"`
	PageNo      int    `json:"pageno" form:"pageno"`
	PageSize    int    `json:"pagesize" form:"pagesize"`
}
