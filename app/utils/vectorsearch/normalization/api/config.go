package api

import "nova-factory-server/app/utils/gateway/v1/config/cfg"

type Config struct {
	Interceptors []*InterceptorConfig `json:"interceptors"`
}

// InterceptorConfig 表示 pipeline 构建和 step 初始化时使用的配置。
type InterceptorConfig struct {
	Enabled    *bool         `yaml:"enabled,omitempty"`
	Name       string        `yaml:"name,omitempty"`
	Type       string        `yaml:"type,omitempty" validate:"required"`
	Properties cfg.CommonCfg `yaml:",inline"`
}
