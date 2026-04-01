package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type endpointState struct {
	name    string
	baseURL string
	apiKey  string
	weight  int
	headers map[string]string
}

type Client struct {
	httpClient          *http.Client
	endpoints           []endpointState
	algorithm           Algorithm
	apiKeyHeader        string
	apiKeyPrefix        string
	agentGatewayHeader  string
	roundRobinCounter   atomic.Uint64
	weightedCursor      atomic.Uint64
	weightedRoundBucket []int
	randomMu            sync.Mutex
	random              *rand.Rand
}

// NewClient 创建一个支持多地址调度、鉴权头注入和网关头透传的 HTTP 客户端。
func NewClient(cfg Config) (*Client, error) {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	transport := cfg.Transport
	if transport == nil {
		transport = &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 100,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     90 * time.Second,
			DisableKeepAlives:   false,
			ForceAttemptHTTP2:   true,
		}
	}
	endpoints := make([]endpointState, 0, len(cfg.Endpoints))
	bucket := make([]int, 0)
	for _, ep := range cfg.Endpoints {
		if !ep.Enabled {
			continue
		}
		base := strings.TrimSpace(ep.BaseURL)
		if base == "" {
			continue
		}
		if _, err := url.ParseRequestURI(base); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidBaseURL, base)
		}
		w := ep.Weight
		if w <= 0 {
			w = 1
		}
		headers := make(map[string]string, len(ep.Headers))
		for k, v := range ep.Headers {
			headers[k] = v
		}
		endpoints = append(endpoints, endpointState{
			name:    ep.Name,
			baseURL: strings.TrimRight(base, "/"),
			apiKey:  ep.APIKey,
			weight:  w,
			headers: headers,
		})
		for i := 0; i < w; i++ {
			bucket = append(bucket, len(endpoints)-1)
		}
	}
	if len(endpoints) == 0 {
		return nil, ErrNoAvailableEndpoint
	}
	algo := cfg.Algorithm
	if algo == "" {
		algo = AlgorithmRoundRobin
	}
	apiKeyHeader := strings.TrimSpace(cfg.APIKeyHeader)
	if apiKeyHeader == "" {
		apiKeyHeader = "Authorization"
	}
	agentGatewayHeader := strings.TrimSpace(cfg.AgentGatewayHeader)
	if agentGatewayHeader == "" {
		agentGatewayHeader = "X-Agent-Gateway"
	}
	return &Client{
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		endpoints:           endpoints,
		algorithm:           algo,
		apiKeyHeader:        apiKeyHeader,
		apiKeyPrefix:        cfg.APIKeyPrefix,
		agentGatewayHeader:  agentGatewayHeader,
		weightedRoundBucket: bucket,
		random:              rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Do 发起请求并按调度策略自动选择上游。
// 返回 HTTP 状态码和消息文本，网络错误或重试失败时返回 error。
func (c *Client) Do(ctx context.Context, req Request) (int, string, error) {
	resp, _, err := c.doRaw(ctx, req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}
	return resp.StatusCode, extractMessage(resp.StatusCode, data), nil
}

// DoRaw 发起请求并返回上游原始响应，适用于 SSE 等流式透传场景。
func (c *Client) DoRaw(ctx context.Context, req Request) (*http.Response, string, error) {
	return c.doRaw(ctx, req)
}

func (c *Client) doRaw(ctx context.Context, req Request) (*http.Response, string, error) {
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	if method == "" {
		return nil, "", ErrInvalidMethod
	}
	total := len(c.endpoints)
	var lastErr error
	used := make(map[int]struct{}, total)
	for len(used) < total {
		idx := c.pickEndpoint(used)
		if idx < 0 {
			break
		}
		used[idx] = struct{}{}
		ep := c.endpoints[idx]
		httpReq, err := c.buildRequest(ctx, method, req, ep)
		if err != nil {
			return nil, "", err
		}
		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("upstream status code %d", resp.StatusCode)
			_ = resp.Body.Close()
			continue
		}
		return resp, ep.name, nil
	}
	if lastErr == nil {
		lastErr = ErrNoAvailableEndpoint
	}
	return nil, "", lastErr
}

// DoJSON 是 Do 的 JSON 辅助方法，自动读取响应体并反序列化到 out。
func (c *Client) DoJSON(ctx context.Context, req Request, out any) (int, string, []byte, error) {
	resp, endpointName, err := c.doRaw(ctx, req)
	if err != nil {
		return 0, endpointName, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, endpointName, nil, err
	}
	if out != nil {
		if err = json.Unmarshal(data, out); err != nil {
			return resp.StatusCode, endpointName, data, err
		}
	}
	return resp.StatusCode, endpointName, data, nil
}

func extractMessage(statusCode int, data []byte) string {
	if len(data) == 0 {
		return http.StatusText(statusCode)
	}
	var payload struct {
		Message string `json:"message"`
		Msg     string `json:"msg"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal(data, &payload); err == nil {
		if payload.Message != "" {
			return payload.Message
		}
		if payload.Msg != "" {
			return payload.Msg
		}
		if payload.Error != "" {
			return payload.Error
		}
	}
	return string(data)
}

// buildRequest 组装最终的 HTTP 请求并注入 endpoint 级与请求级 header。
func (c *Client) buildRequest(ctx context.Context, method string, req Request, ep endpointState) (*http.Request, error) {
	fullURL := ep.baseURL + "/" + strings.TrimLeft(req.Path, "/")
	if strings.TrimSpace(req.Path) == "" {
		fullURL = ep.baseURL
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, fullURL, bytes.NewReader(req.Body))
	if err != nil {
		return nil, err
	}
	if len(req.Query) > 0 {
		query := httpReq.URL.Query()
		for k, v := range req.Query {
			for _, item := range v {
				query.Add(k, item)
			}
		}
		httpReq.URL.RawQuery = query.Encode()
	}
	for k, v := range ep.headers {
		httpReq.Header.Set(k, v)
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	if ep.apiKey != "" {
		httpReq.Header.Set(c.apiKeyHeader, c.apiKeyPrefix+ep.apiKey)
	}
	if req.AgentGateway != "" {
		httpReq.Header.Set(c.agentGatewayHeader, req.AgentGateway)
	}
	return httpReq, nil
}

// pickEndpoint 按算法选择当前请求要访问的 endpoint。
func (c *Client) pickEndpoint(used map[int]struct{}) int {
	switch c.algorithm {
	case AlgorithmRandom:
		return c.pickRandom(used)
	case AlgorithmWeightedRound:
		return c.pickWeightedRound(used)
	default:
		return c.pickRoundRobin(used)
	}
}

// pickRoundRobin 轮询选择可用 endpoint。
func (c *Client) pickRoundRobin(used map[int]struct{}) int {
	total := len(c.endpoints)
	if total == 0 {
		return -1
	}
	start := int(c.roundRobinCounter.Add(1)-1) % total
	for i := 0; i < total; i++ {
		idx := (start + i) % total
		if _, ok := used[idx]; !ok {
			return idx
		}
	}
	return -1
}

// pickRandom 随机选择可用 endpoint。
func (c *Client) pickRandom(used map[int]struct{}) int {
	available := make([]int, 0, len(c.endpoints)-len(used))
	for i := range c.endpoints {
		if _, ok := used[i]; ok {
			continue
		}
		available = append(available, i)
	}
	if len(available) == 0 {
		return -1
	}
	c.randomMu.Lock()
	idx := available[c.random.Intn(len(available))]
	c.randomMu.Unlock()
	return idx
}

// pickWeightedRound 按权重轮询选择可用 endpoint。
func (c *Client) pickWeightedRound(used map[int]struct{}) int {
	total := len(c.weightedRoundBucket)
	if total == 0 {
		return c.pickRoundRobin(used)
	}
	start := int(c.weightedCursor.Add(1)-1) % total
	for i := 0; i < total; i++ {
		idx := c.weightedRoundBucket[(start+i)%total]
		if _, ok := used[idx]; !ok {
			return idx
		}
	}
	return -1
}
