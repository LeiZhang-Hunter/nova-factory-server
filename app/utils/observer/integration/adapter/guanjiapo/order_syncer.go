package guanjiapo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/errgo.v2/errors"
	"io"
	"net/http"
	"net/url"
	"nova-factory-server/app/utils/observer/integration/api"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/result"
	"strings"
	"time"
)

type orderSyncer struct {
	tokenURL string
}

func newOrderSyncer(tokenURL string) api.OrderSyncer {
	return &orderSyncer{
		tokenURL: tokenURL,
	}
}

// makeSign 签名
func (s *orderSyncer) makeSign(timestamp string, token string, cfg *ConfigSnapshot, method string, body string) (string, error) {
	var param map[string]string = make(map[string]string)
	param["app_key"] = cfg.Credentials.AppKey
	param["v"] = "1.0"
	param["format"] = "json"
	param["sign_method"] = "md5"
	param["method"] = method
	param["timestamp"] = timestamp
	param["token"] = token

	return generateMD5Sign(param, body, cfg.Credentials.AppSecret)
}

func (s *orderSyncer) SyncOrders(ctx context.Context, req event.OrderEvent) (result.OrderSyncResponse, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if req == nil || len(req.Orders()) == 0 {
		return nil, errors.New("orders不能为空")
	}
	snapshot, err := parseSnapshot(req.Config())
	if err != nil {
		return nil, err
	}
	token, err := resolveAccessToken(ctx, snapshot, req.Cache())
	if err != nil {
		return nil, err
	}
	openapiURL := s.tokenURL
	if strings.TrimSpace(snapshot.BaseURL) != "" {
		openapiURL = strings.TrimRight(strings.TrimSpace(snapshot.BaseURL), "/") + "/openapi"
	}
	parse, err := url.Parse(openapiURL)
	if err != nil {
		return nil, err
	}
	params := parse.Query()
	params.Set("app_key", snapshot.Credentials.AppKey)
	params.Set("v", "1.0")
	params.Set("format", "json")
	params.Set("sign_method", "md5")
	params.Set("method", "emall.order.synchronize")
	params.Set("timestamp", timestamp)
	params.Set("token", token)
	body := map[string]any{
		"orders": req.Orders(),
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	sign, err := s.makeSign(timestamp, token, snapshot, "emall.order.synchronize", string(payload))
	if err != nil {
		return nil, err
	}
	params.Set("sign", sign)
	parse.RawQuery = params.Encode()
	fmt.Println(parse.String())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, parse.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ret := &OrderSyncResponse{}
	if err = json.Unmarshal(respBytes, ret); err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(ret.Message)
		if msg == "" {
			msg = string(respBytes)
		}
		return nil, errors.New("订单同步失败: " + msg)
	}
	if ret.Code != 0 {
		return nil, errors.New(ret.Message)
	}
	return ret, nil
}
