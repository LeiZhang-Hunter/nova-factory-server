package integration

import (
	"nova-factory-server/app/business/erp/core/integration/api"
	"testing"
)

func TestRegistryCreateBuiltin(t *testing.T) {
	c, err := Create(api.KindGuanJiaPo)
	if err != nil {
		t.Fatalf("create guanjiapo failed: %v", err)
	}
	if c == nil {
		t.Fatal("client is nil")
	}
	if c.Kind() != api.KindGuanJiaPo {
		t.Fatalf("unexpected kind: %s", c.Kind())
	}
}

func TestRegistryCreateUnsupported(t *testing.T) {
	_, err := Create(api.Kind("未知"))
	if err == nil {
		t.Fatal("expected error for unsupported kind")
	}
}

func TestRegistryCreateByType(t *testing.T) {
	c, err := CreateByType("管家婆")
	if err != nil {
		t.Fatalf("create by type failed: %v", err)
	}
	if c == nil {
		t.Fatal("client is nil")
	}
	if c.Kind() != api.KindGuanJiaPo {
		t.Fatalf("unexpected kind: %s", c.Kind())
	}
}

func TestResolveCheckURL(t *testing.T) {
	s := &api.ConfigSnapshot{
		CheckURL:       "https://a.example.com/check",
		LoginStatusURL: "https://a.example.com/login/status",
		BaseURL:        "https://a.example.com",
	}
	if got := api.ResolveCheckURL("https://override/check", s, "/fallback"); got != "https://override/check" {
		t.Fatalf("override not applied: %s", got)
	}
	if got := api.ResolveCheckURL("", s, "/fallback"); got != "https://a.example.com/check" {
		t.Fatalf("checkUrl not applied: %s", got)
	}
	s.CheckURL = ""
	if got := api.ResolveCheckURL("", s, "/fallback"); got != "https://a.example.com/login/status" {
		t.Fatalf("loginStatusUrl not applied: %s", got)
	}
	s.LoginStatusURL = ""
	if got := api.ResolveCheckURL("", s, "/fallback"); got != "https://a.example.com/fallback" {
		t.Fatalf("baseUrl fallback not applied: %s", got)
	}
}
