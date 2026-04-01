package client

import "errors"

var (
	// ErrNoAvailableEndpoint 表示没有可用上游地址。
	ErrNoAvailableEndpoint = errors.New("no available endpoint")
	// ErrInvalidMethod 表示请求方法为空或不合法。
	ErrInvalidMethod = errors.New("invalid method")
	// ErrInvalidBaseURL 表示配置中的地址不是合法 URL。
	ErrInvalidBaseURL = errors.New("invalid base url")
)
