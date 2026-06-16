package client

import (
	"testing"
)

func TestGenerateMD5Sign(t *testing.T) {
	params := map[string]string{
		"app_key":     "test_key",
		"v":           "1.0",
		"format":      "json",
		"sign_method": "md5",
		"method":      "emall.goods.get",
		"timestamp":   "2025-01-01 12:00:00",
		"token":       "test_token",
	}
	body := `{"page":1,"pagesize":10}`
	secret := "test_secret"

	sign, err := GenerateMD5Sign(params, body, secret)
	if err != nil {
		t.Fatal(err)
	}
	if sign == "" {
		t.Error("expected non-empty sign")
	}

	// Verify deterministic: same inputs produce same sign
	sign2, _ := GenerateMD5Sign(params, body, secret)
	if sign != sign2 {
		t.Error("expected deterministic sign output")
	}
}

func TestGenerateMD5SignEmptySecret(t *testing.T) {
	_, err := GenerateMD5Sign(map[string]string{"k": "v"}, "body", "")
	if err == nil {
		t.Error("expected error for empty secret")
	}
}

func TestGenerateMD5SignSkipsSignParam(t *testing.T) {
	params := map[string]string{
		"app_key": "key",
		"sign":    "existing_sign",
	}
	body := "body"
	secret := "secret"

	sign1, _ := GenerateMD5Sign(params, body, secret)

	delete(params, "sign")
	sign2, _ := GenerateMD5Sign(params, body, secret)

	if sign1 != sign2 {
		t.Error("expected sign param to be skipped, producing same result")
	}
}

func TestResolveOpenAPIURL(t *testing.T) {
	t.Run("uses base URL from snapshot", func(t *testing.T) {
		snapshot := &ConfigSnapshot{BaseURL: "https://api.example.com/erp"}
		url := ResolveOpenAPIURL("http://default.url/token", snapshot)
		if url != "https://api.example.com/erp/openapi" {
			t.Errorf("expected base+openapi, got %s", url)
		}
	})

	t.Run("falls back to default URL", func(t *testing.T) {
		snapshot := &ConfigSnapshot{}
		url := ResolveOpenAPIURL("http://default.url/token", snapshot)
		if url != "http://default.url/token" {
			t.Errorf("expected default URL, got %s", url)
		}
	})

	t.Run("nil snapshot falls back", func(t *testing.T) {
		url := ResolveOpenAPIURL("http://default.url/token", nil)
		if url != "http://default.url/token" {
			t.Errorf("expected default URL for nil snapshot, got %s", url)
		}
	})
}

func TestParseSnapshot(t *testing.T) {
	t.Run("parses valid JSON config", func(t *testing.T) {
		cfg := &testConfig{data: `{"systemName":"test","credentials":{"appKey":"k","appSecret":"s"}}`}
		snapshot, err := ParseSnapshot(cfg)
		if err != nil {
			t.Fatal(err)
		}
		if snapshot.SystemName != "test" {
			t.Errorf("expected SystemName=test, got %s", snapshot.SystemName)
		}
		if snapshot.Credentials.AppKey != "k" {
			t.Errorf("expected AppKey=k, got %s", snapshot.Credentials.AppKey)
		}
	})

	t.Run("nil config returns error", func(t *testing.T) {
		_, err := ParseSnapshot(nil)
		if err == nil {
			t.Error("expected error for nil config")
		}
	})

	t.Run("empty data returns empty snapshot", func(t *testing.T) {
		cfg := &testConfig{data: ""}
		snapshot, err := ParseSnapshot(cfg)
		if err != nil {
			t.Fatal(err)
		}
		if snapshot.SystemName != "" {
			t.Errorf("expected empty SystemName, got %s", snapshot.SystemName)
		}
	})
}

type testConfig struct {
	data string
}

func (c *testConfig) GetOverrideURL() string      { return "" }
func (c *testConfig) GetMetadata() map[string]any { return nil }
func (c *testConfig) GetData() string             { return c.data }
func (c *testConfig) GetType() string             { return "gjp_v1" }
func (c *testConfig) GetStatus() *bool            { s := true; return &s }
