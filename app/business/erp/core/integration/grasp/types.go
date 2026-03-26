package grasp

import (
	"encoding/json"
	"errors"
	"nova-factory-server/app/business/erp/setting/settingModels"
	"strings"
)

type Credentials struct {
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
}
type ConfigSnapshot struct {
	SystemName     string            `json:"systemName"`
	Credentials    Credentials       `json:"credentials"`
	CheckURL       string            `json:"checkUrl"`
	LoginStatusURL string            `json:"loginStatusUrl"`
	BaseURL        string            `json:"baseUrl"`
	RedirectURL    string            `json:"redirect_url"`
	State          string            `json:"state"`
	Token          string            `json:"token"`
	AccessToken    string            `json:"accessToken"`
	Cookie         string            `json:"cookie"`
	Headers        map[string]string `json:"headers"`
}

func ParseSnapshot(cfg *settingModels.IntegrationConfig) (*ConfigSnapshot, error) {
	if cfg == nil {
		return nil, errors.New("integration config不能为空")
	}
	s := &ConfigSnapshot{}
	if strings.TrimSpace(cfg.Data) == "" {
		return s, nil
	}
	if err := json.Unmarshal([]byte(cfg.Data), s); err != nil {
		return nil, err
	}
	return s, nil
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
