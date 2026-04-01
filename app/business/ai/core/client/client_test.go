package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRoundRobinWithGatewayHeader 验证轮询调度、API Key 注入与 Agent 网关头透传。
func TestRoundRobinWithGatewayHeader(t *testing.T) {
	var gotHost []string
	var gotAPIKey []string
	var gotGateway []string
	makeServer := func(host string) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotHost = append(gotHost, host)
			gotAPIKey = append(gotAPIKey, r.Header.Get("Authorization"))
			gotGateway = append(gotGateway, r.Header.Get("X-Agent-Gateway"))
			_, _ = w.Write([]byte(`{"ok":true}`))
		}))
	}
	s1 := makeServer("s1")
	defer s1.Close()
	s2 := makeServer("s2")
	defer s2.Close()

	c, err := NewClient(Config{
		Algorithm: AlgorithmRoundRobin,
		Endpoints: []Endpoint{
			{Name: "a", BaseURL: s1.URL, APIKey: "k1", Enabled: true},
			{Name: "b", BaseURL: s2.URL, APIKey: "k2", Enabled: true},
		},
		APIKeyHeader:       "Authorization",
		APIKeyPrefix:       "Bearer ",
		AgentGatewayHeader: "X-Agent-Gateway",
	})
	if err != nil {
		t.Fatalf("new client err: %v", err)
	}

	for i := 0; i < 4; i++ {
		statusCode, _, reqErr := c.Do(context.Background(), Request{
			Method:       http.MethodGet,
			Path:         "/v1/ping",
			AgentGateway: "gw-1",
		})
		if reqErr != nil {
			t.Fatalf("do req err: %v", reqErr)
		}
		if statusCode != http.StatusOK {
			t.Fatalf("unexpected status: %d", statusCode)
		}
	}
	if len(gotHost) != 4 {
		t.Fatalf("got request count %d", len(gotHost))
	}
	if gotHost[0] != "s1" || gotHost[1] != "s2" || gotHost[2] != "s1" || gotHost[3] != "s2" {
		t.Fatalf("round robin failed: %v", gotHost)
	}
	if gotAPIKey[0] != "Bearer k1" || gotAPIKey[1] != "Bearer k2" {
		t.Fatalf("api key header failed: %v", gotAPIKey)
	}
	for _, v := range gotGateway {
		if v != "gw-1" {
			t.Fatalf("gateway header failed: %v", gotGateway)
		}
	}
}

// TestFailoverOnServerError 验证上游返回 5xx 时会自动切换到下一个地址。
func TestFailoverOnServerError(t *testing.T) {
	first := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"ok":false}`))
	}))
	defer first.Close()
	second := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer second.Close()

	c, err := NewClient(Config{
		Algorithm: AlgorithmRoundRobin,
		Endpoints: []Endpoint{
			{Name: "bad", BaseURL: first.URL, APIKey: "k1", Enabled: true},
			{Name: "good", BaseURL: second.URL, APIKey: "k2", Enabled: true},
		},
	})
	if err != nil {
		t.Fatalf("new client err: %v", err)
	}
	statusCode, message, reqErr := c.Do(context.Background(), Request{
		Method: http.MethodPost,
		Path:   "/v1/chat",
		Body:   []byte(`{"message":"hello"}`),
	})
	if reqErr != nil {
		t.Fatalf("do req err: %v", reqErr)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", statusCode)
	}
	if message != `{"ok":true}` {
		t.Fatalf("unexpected message: %s", message)
	}
}
