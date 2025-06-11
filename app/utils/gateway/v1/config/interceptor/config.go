package interceptor

import (
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
)

type Config struct {
	Enabled    *bool         `yaml:"enabled,omitempty"`
	Name       string        `yaml:"name,omitempty"`
	Type       string        `yaml:"type,omitempty" validate:"required"`
	Properties cfg.CommonCfg `yaml:",inline"`
}
