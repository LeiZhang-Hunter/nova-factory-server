package client

import (
	"encoding/json"
	"errors"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/kind"
	"strings"
)

// Kind 集成系统类型标识
const (
	KindGuanJiaPo kind.Kind = "gjp_v1"
)

// Credentials 管家婆应用授权信息
type Credentials struct {
	AppKey          string `json:"appKey"`
	AppSecret       string `json:"appSecret"`
	Selfmallaccount string `json:"selfmallaccount"`
}

// ConfigSnapshot 管家婆集成配置快照
type ConfigSnapshot struct {
	SystemName      string            `json:"systemName"`
	Credentials     Credentials       `json:"credentials"`
	CheckURL        string            `json:"checkUrl"`
	BaseURL         string            `json:"baseUrl"`
	RedirectURL     string            `json:"redirect_url"`
	State           string            `json:"state"`
	Token           string            `json:"token"`
	AccessToken     string            `json:"accessToken"`
	Cookie          string            `json:"cookie"`
	Headers         map[string]string `json:"headers"`
	CodeTTL         string            `json:"codeTTL"`
	TokenTTL        string            `json:"tokenTTL"`
	RefreshTokenTTL string            `json:"refreshTokenTTL"`
}

// ParseSnapshot 解析集成配置JSON为配置快照
func ParseSnapshot(cfg config.Config) (*ConfigSnapshot, error) {
	if cfg == nil {
		return nil, errors.New("integration config不能为空")
	}
	s := &ConfigSnapshot{}
	if strings.TrimSpace(cfg.GetData()) == "" {
		return s, nil
	}
	if err := json.Unmarshal([]byte(cfg.GetData()), s); err != nil {
		return nil, err
	}
	return s, nil
}

// ApplyDefaults 为未配置的 TTL 字段填充默认值
func (c *ConfigSnapshot) ApplyDefaults() {
	if c.CodeTTL == "" {
		c.CodeTTL = "10m"
	}
	if c.TokenTTL == "" {
		c.TokenTTL = "24h"
	}
	if c.RefreshTokenTTL == "" {
		c.RefreshTokenTTL = "720h"
	}
}
