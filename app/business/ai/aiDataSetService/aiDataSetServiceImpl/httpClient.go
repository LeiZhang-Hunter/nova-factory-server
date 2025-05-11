package aiDataSetServiceImpl

import (
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func NewHttpClient() *http.Client {
	var config RagFlowConfig
	err := viper.UnmarshalKey("dataSet", &config)
	if err != nil {
		panic(err)
	}
	// 创建自定义的Transport
	transport := &http.Transport{
		MaxIdleConns:        100,              // 整个客户端的最大空闲连接数
		MaxIdleConnsPerHost: 60,               // 每个主机的最大空闲连接数
		MaxConnsPerHost:     50,               // 每个主机的最大连接数
		IdleConnTimeout:     60 * time.Second, // 空闲连接的超时时间
		DisableKeepAlives:   false,            // 不禁用连接保持活动
		ForceAttemptHTTP2:   true,             // 尝试使用HTTP/2
	}

	// 创建HTTP客户端并设置自定义的Transport
	client := &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second, // 请求的总超时时间
	}
	return client
}
