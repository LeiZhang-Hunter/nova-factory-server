package client

import (
	"net/http"
	"net/url"
	"time"
)

type Algorithm string

const (
	// AlgorithmRoundRobin 按顺序轮询选择地址。
	AlgorithmRoundRobin Algorithm = "round_robin"
	// AlgorithmRandom 随机选择地址。
	AlgorithmRandom Algorithm = "random"
	// AlgorithmWeightedRound 按权重轮询选择地址。
	AlgorithmWeightedRound Algorithm = "weighted_round"
)

// Endpoint 描述一个可被调度的上游地址及其认证信息。
type Endpoint struct {
	Name    string
	BaseURL string
	APIKey  string
	Enabled bool
	Weight  int
	Headers map[string]string
}

// Config 是客户端初始化配置。
type Config struct {
	Endpoints          []Endpoint
	Algorithm          Algorithm
	Timeout            time.Duration
	APIKeyHeader       string
	APIKeyPrefix       string
	AgentGatewayHeader string
	Transport          http.RoundTripper
}

// Request 是一次请求调用参数。
type Request struct {
	Method       string
	Path         string
	Query        url.Values
	Headers      map[string]string
	Body         []byte
	AgentGateway string
}
