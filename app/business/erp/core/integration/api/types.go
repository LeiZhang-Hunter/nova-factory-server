package api

import (
	"context"
	"nova-factory-server/app/business/erp/setting/settingmodels"
	"strings"
)

type Kind string

const (
	KindGuanJiaPo Kind = "gjp_v1"
	KindKingdee   Kind = "金蝶"
)

type LoginState struct {
	Online   bool   `json:"online"`
	Message  string `json:"message"`
	Type     string `json:"type"`
	CheckURL string `json:"checkUrl"`
	Raw      string `json:"raw,omitempty"`
}

type Client interface {
	Kind() Kind
	CheckLoginState(ctx context.Context, cfg *settingmodels.IntegrationConfig, overrideURL string, overrideRedirectURL string) (*LoginState, error)
}

type ConfigSnapshot struct {
	CheckURL       string            `json:"checkUrl"`
	LoginStatusURL string            `json:"loginStatusUrl"`
	BaseURL        string            `json:"baseUrl"`
	AppKey         string            `json:"appkey"`
	AppSecret      string            `json:"appsecret"`
	RedirectURL    string            `json:"redirect_url"`
	State          string            `json:"state"`
	Token          string            `json:"token"`
	AccessToken    string            `json:"accessToken"`
	Cookie         string            `json:"cookie"`
	Headers        map[string]string `json:"headers"`
}

func ResolveCheckURL(overrideURL string, snapshot *ConfigSnapshot, fallbackPath string) string {
	if strings.TrimSpace(overrideURL) != "" {
		return strings.TrimSpace(overrideURL)
	}
	if snapshot == nil {
		return ""
	}
	if strings.TrimSpace(snapshot.CheckURL) != "" {
		return strings.TrimSpace(snapshot.CheckURL)
	}
	if strings.TrimSpace(snapshot.LoginStatusURL) != "" {
		return strings.TrimSpace(snapshot.LoginStatusURL)
	}
	if strings.TrimSpace(snapshot.BaseURL) != "" {
		return strings.TrimRight(strings.TrimSpace(snapshot.BaseURL), "/") + fallbackPath
	}
	return ""
}
