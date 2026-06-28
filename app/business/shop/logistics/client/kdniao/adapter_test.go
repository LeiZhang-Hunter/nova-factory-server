package kdniao

import "testing"

func TestNewAdapter_NilConfig(t *testing.T) {
	adapter, err := NewAdapter(nil)
	if adapter == nil {
		t.Fatal("expected non-nil adapter even on nil config")
	}
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestNewAdapter_EmptyConfig(t *testing.T) {
	adapter, err := NewAdapter(&Config{})
	if err == nil {
		t.Fatal("expected error for empty config")
	}
	if adapter == nil {
		t.Fatal("expected non-nil adapter")
	}
}

func TestAdapter_Query_NotConfigured(t *testing.T) {
	a := &Adapter{}
	result, err := a.Query("SF", "SF123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success() {
		t.Fatal("expected Success=false for unconfigured adapter")
	}
	if result.Reason() == "" {
		t.Fatal("expected non-empty Reason")
	}
}

func TestAdapter_Query_RealCall(t *testing.T) {
	a, err := NewAdapter(&Config{
		EBusinessID: "1926491",
		AppKey:      "42e9951f-f2af-4beb-be94-fa3fa8b4cfb0",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := a.Query("STO", "773367326370601")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success() {
		t.Fatalf("expected Success=true, got Reason: %s", result.Reason())
	}
	// 测试单号暂无轨迹，State 应为 "0"
	if result.State() != "0" {
		t.Fatalf("expected State=0 for no-track single, got %s", result.State())
	}
}

func TestNewAdapter_ValidConfig_NoSDKCall(t *testing.T) {
	adapter, err := NewAdapter(&Config{EBusinessID: "fake", AppKey: "fake"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if adapter == nil {
		t.Fatal("expected non-nil adapter")
	}
}

func TestNewAdapter_ReturnsError(t *testing.T) {
	_, err := NewAdapter(nil)
	if err == nil || err.Error() != "cfg is nil" {
		t.Fatalf("expected 'cfg is nil', got: %v", err)
	}
}
