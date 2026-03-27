package grasp

import (
	"encoding/json"
	"errors"
	"nova-factory-server/app/business/erp/setting/settingModels"
	"strings"
)

// Credentials 管家婆应用授权信息
type Credentials struct {
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
}

// ConfigSnapshot 集成配置快照
type ConfigSnapshot struct {
	SystemName  string            `json:"systemName"`
	Credentials Credentials       `json:"credentials"`
	CheckURL    string            `json:"checkUrl"`
	BaseURL     string            `json:"baseUrl"`
	RedirectURL string            `json:"redirect_url"`
	State       string            `json:"state"`
	Token       string            `json:"token"`
	AccessToken string            `json:"accessToken"`
	Cookie      string            `json:"cookie"`
	Headers     map[string]string `json:"headers"`
}

// ParseSnapshot 解析集成配置JSON为配置快照
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
